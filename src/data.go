package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var METADATA_FILE string = "metadata.mt"
var INDEX_DATA_FILE string = "index.idx"
var DATA_FILE_EXTENSION string = ".dat"
var TAILER_EXTENSION string = ".trl"
var DATA_FILE_PREFIX string = "data-"
var MAX_SEGMENT_SIZE int64 = 512 * 1024 * 1024

var currentDirectory string

type IndexData struct {
	MessageId    int32
	SegmentId    int32
	Size         int32
	OffsetInData int64
	Time         int64
}

type Metadata struct {
	CurrentSegment     int32
	CurrentOffsetWrite int64
	CountSegment       int32
	CountMessage       int32
}

type Tailer struct {
	Messageid    int64
	LastReadTime int64
}

type DataInfo struct {
	Size int32
	Time int64
}

// func main() {
// 	currentDirectory = getCurrentDirectory()
// }

func writeIndexFile(file *os.File, index IndexData) {
	err := binary.Write(file, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}
}

func writeDataFile(file *os.File, data DataInfo, reader *bufio.Reader) {
	binary.Write(file, binary.LittleEndian, data.Size)
	binary.Write(file, binary.LittleEndian, data.Time)
	tee := io.TeeReader(reader, os.Stdout)
	_, err := io.Copy(file, tee)

	if err != nil {
		panic(err)
	}

}

func readIndexFile(file *os.File, index *IndexData, offset int64) {
	file.Seek(offset, 0)
	err := binary.Read(file, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}
}

func writeTailerFile(file *os.File, tailer Tailer) {
	err := binary.Write(file, binary.LittleEndian, tailer)
	if err != nil {
		panic(err)
	}
}

func readTailerFile(file *os.File, tailer *Tailer) {
	err := binary.Read(file, binary.LittleEndian, tailer)
	if err != nil {
		panic(err)
	}
}

func writeMetadataFile(file *os.File, metadata Metadata) {

	err := binary.Write(file, binary.LittleEndian, metadata)

	if err != nil {
		panic(err)
	}
}

func readMetadataFile(file *os.File, metadata *Metadata) {
	file.Seek(0, io.SeekStart)
	err := binary.Read(file, binary.LittleEndian, metadata)
	if err != nil {
		panic(err)
	}
}

func writeMessage(header Header, reader *bufio.Reader) {

	queuePath := getQueuePath(header.QueueName)

	fmt.Println(queuePath)
	if isDirExists(queuePath) {
		writeExistQueue(header, reader)
	} else {
		os.MkdirAll(queuePath, 0777)
		writeNewQueue(header, reader)
	}
}

func writeNewQueue(header Header, reader *bufio.Reader) {
	currentDirectory = getCurrentDirectory()
	metadata := Metadata{}
	metadata = Metadata{1, 12 + int64(header.TailerOrPayloadSize), 1, 1}

	timeNow := getTimeNow()

	metadataPath := getMetadataPath(header.QueueName)
	metadataFile, err := os.OpenFile(metadataPath, os.O_WRONLY|os.O_CREATE, 0666)
	defer metadataFile.Close()
	if err != nil {
		panic(err)
	}

	writeMetadataFile(metadataFile, metadata)

	indexFile, err := os.OpenFile(getIndexPath(header.QueueName), os.O_WRONLY|os.O_CREATE, 0666)
	defer indexFile.Close()

	index := IndexData{1, 1, header.TailerOrPayloadSize, 0, timeNow}
	writeIndexFile(indexFile, index)

	dataFile, err := os.OpenFile(getSegmentPath(header.QueueName, 1), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer dataFile.Close()

	dataInfo := DataInfo{header.TailerOrPayloadSize, timeNow}

	writeDataFile(dataFile, dataInfo, reader)

	if err != nil {
		panic(err)
	}
}

func writeExistQueue(header Header, reader *bufio.Reader) {
	metadata := Metadata{}
	metadataFile, err := os.Open(getMetadataPath(header.QueueName))
	defer metadataFile.Close()

	if err != nil {
		panic(err)
	}
	timeNow := getTimeNow()

	err = binary.Read(metadataFile, binary.LittleEndian, &metadata)
	metadata.CurrentOffsetWrite = metadata.CurrentOffsetWrite + 12 + int64(header.TailerOrPayloadSize)
	isNewSegment := MAX_SEGMENT_SIZE-(metadata.CurrentOffsetWrite+12) < int64(header.TailerOrPayloadSize)

	if isNewSegment {
		metadata.CountSegment = metadata.CountSegment + 1
		metadata.CurrentSegment = metadata.CurrentSegment + 1
	}
	metadata.CountMessage = metadata.CountMessage + 1

	indexFile, err := os.OpenFile(getIndexPath(header.QueueName), os.O_WRONLY|os.O_APPEND, 0666)
	defer indexFile.Close()

	index := IndexData{
		metadata.CountMessage,
		metadata.CurrentSegment,
		header.TailerOrPayloadSize,
		metadata.CurrentOffsetWrite,
		timeNow}

	writeIndexFile(indexFile, index)

	dataFile, err := os.OpenFile(getSegmentPath(header.QueueName, metadata.CurrentSegment), os.O_WRONLY|os.O_APPEND, 0666)
	defer dataFile.Close()

	dataInfo := DataInfo{header.TailerOrPayloadSize, timeNow}

	writeDataFile(dataFile, dataInfo, reader)
}

func getMessage(header Header, writer *bufio.Writer) {

	metadata := Metadata{}
	metadataFile, err := os.Open(getMetadataPath(header.QueueName))
	defer metadataFile.Close()

	queuePath := getQueuePath(header.QueueName)
	timeNow := getTimeNow()

	if !isDirExists(queuePath) {
		writer.Write([]byte("Queue not exist"))
		return
	}

	tailerPath := getTailerPath(header.QueueName, header.TailerName)
	os.MkdirAll(tailerPath, 0777)

	tailer := Tailer{}

	tailerFile, err := os.OpenFile(tailerPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer tailerFile.Close()

	if err != nil {
		panic(err)
	}

	tailerFileInfo, errStat := os.Stat(tailerPath)

	if errStat != nil {
		panic(err)
	}

	if tailerFileInfo.Size() == 0 {
		tailer.Messageid = 1
		tailer.LastReadTime = timeNow
	} else {
		readTailerFile(tailerFile, &tailer)
	}

	if int64(metadata.CountMessage) == tailer.Messageid {
		writer.Write(collectErrorBuf("New messages not exist"))
		return
	}

	tailer.Messageid = tailer.Messageid + 1
	tailer.LastReadTime = timeNow

	index := IndexData{}

	indexFile, err := os.OpenFile(getIndexPath(header.QueueName), os.O_RDWR, 0666)

	readIndexFile(indexFile, &index, 28*tailer.Messageid)

}

func getCurrentDirectory() string {
	return filepath.FromSlash("C:/dev/projects/my/Go-language/un-queue-go/resources")
}

func getSegmentName(segmentId int32) string {
	segmentString := fmt.Sprintf("%08d", segmentId)
	return DATA_FILE_PREFIX + segmentString + DATA_FILE_EXTENSION
}

func getSegmentPath(queueName string, segmentId int32) string {
	return filepath.Join(currentDirectory, queueName, getSegmentName(segmentId))
}

func getMetadataPath(queueName string) string {
	return filepath.FromSlash(filepath.Join(currentDirectory, queueName, METADATA_FILE))
}

func getIndexPath(queueName string) string {
	return filepath.Join(currentDirectory, queueName, INDEX_DATA_FILE)
}

func getQueuePath(queueName string) string {
	return filepath.FromSlash(filepath.Join(currentDirectory, queueName))
}

func getTailerPath(queueName string, tailer string) string {
	return filepath.Join(currentDirectory, queueName, "tailer", tailer, TAILER_EXTENSION)
}
