package main

import (
	"bufio"
	"bytes"
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

// func main() {
// 	currentDirectory = getCurrentDirectory()
// }

func writeIndexFile(file *os.File, index IndexData) {
	err := binary.Write(file, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}
}

func writeDataFile(file *os.File, size int32, time int64) {

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, size)
	binary.Write(buf, binary.LittleEndian, time)

	err := binary.Write(file, binary.LittleEndian, buf)
	if err != nil {
		panic(err)
	}
}

func readIndexFile(file *os.File, index IndexData) {
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

func readTailerFile(file *os.File, tailer Tailer) {
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

func readMetadataFile(file *os.File, metadata Metadata) {
	err := binary.Read(file, binary.LittleEndian, metadata)
	if err != nil {
		panic(err)
	}
}

func writeMessage(header Header, reader *bufio.Reader) {

	metadataPath := getMetadataPath(header.QueueName)
	if fileExists(metadataPath) {

	} else {

	}
}

func writeNewQueue(header Header, reader *bufio.Reader) {
	currentDirectory = getCurrentDirectory()
	metadata := Metadata{}
	metadata = Metadata{1, 12 + int64(header.TailerOrPayloadSize), 1, 1}

	metadataPath := filepath.FromSlash(getMetadataPath(header.QueueName))

	fmt.Println(currentDirectory)
	fmt.Println(len(header.QueueName))
	fmt.Println(metadataPath)

	metadataFile, err := os.OpenFile(metadataPath, os.O_WRONLY|os.O_CREATE, 0666)
	defer metadataFile.Close()
	if err != nil {
		panic(err)
	}

	writeMetadataFile(metadataFile, metadata)

	indexFile, err := os.OpenFile(filepath.FromSlash(getIndexPath(header.QueueName)), os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	timeNow := getTimeNow()

	index := IndexData{1, 1, header.TailerOrPayloadSize, 0, timeNow}
	writeIndexFile(indexFile, index)

	dataFile, err := os.OpenFile(filepath.FromSlash(getSegmentPath(header.QueueName, 1)), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	writeDataFile(dataFile, header.TailerOrPayloadSize, timeNow)

	_, err = io.Copy(dataFile, reader)

	if err != nil {
		panic(err)
	}
}

func writeExistQueue(header Header, reader *bufio.Reader) {
	metadata := Metadata{}
	file, err := os.Open(getMetadataPath(header.QueueName))
	if err != nil {
		panic(err)
	}
	err = binary.Read(file, binary.LittleEndian, &metadata)
	metadata.CurrentOffsetWrite = metadata.CurrentOffsetWrite + 12 + int64(header.TailerOrPayloadSize)
	isNewSegment := MAX_SEGMENT_SIZE-(metadata.CurrentOffsetWrite+12) < int64(header.TailerOrPayloadSize)

	if isNewSegment {
		metadata.CountSegment = metadata.CountSegment + 1
		metadata.CurrentSegment = metadata.CurrentSegment + 1
	}
	metadata.CountMessage = metadata.CountMessage + 1
}

func getMessage() {

}

func getCurrentDirectory() string {
	//	return "C:\\temp"
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
	return currentDirectory + "/" + queueName + "/" + METADATA_FILE
}

func getIndexPath(queueName string) string {
	return currentDirectory + "/" + queueName + "/" + INDEX_DATA_FILE
}
