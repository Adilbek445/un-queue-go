package main

import (
	"bufio"
	"fmt"
	"os"
)

var METADATA_FILE string = "metadata.mt"
var INDEX_DATA_FILE string = "index.idx"
var DATA_FILE_EXTENSION string = ".dat"
var TAILER_EXTENSION string = ".trl"
var DATA_FILE_PREFIX string = "data-"
var MAX_SEGMENT_SIZE int64 = 512 * 1024 * 1024
var INDEX_RECORD_SIZE = 28
var DATA_HEADER_SIZE = 12

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
	QueueDataSize      int64
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

type QueueStat struct {
	CountSegment       int32
	CountMessage       int32
	FirstWriteTime     int64
	LastWriteTime      int64
	TailerCountMessage int32
	TailerLastTime     int64
	QueueDataSize      int64
}

func writeMessage(header Header, reader *bufio.Reader) {

	queuePath := getQueuePath(header.QueueName)
	metadataPath := getMetadataPath(header.QueueName)
	fmt.Println(queuePath)
	if isDirExists(queuePath) && IsFileExists(metadataPath) {
		writeExistQueue(header, reader)
	} else {
		os.MkdirAll(queuePath, 0777)
		writeNewQueue(header, reader)
	}
}

func writeNewQueue(header Header, reader *bufio.Reader) {
	currentDirectory = getCurrentDirectory()
	metadata := Metadata{}
	metadata = Metadata{1, int64(DATA_HEADER_SIZE) + int64(header.TailerOrPayloadSize), int64(header.TailerOrPayloadSize), 1, 1}

	timeNow := getTimeNow()

	metadataPath := getMetadataPath(header.QueueName)
	metadataFile, metaOpenErr := os.OpenFile(metadataPath, os.O_WRONLY|os.O_CREATE, 0777)

	if metaOpenErr != nil {
		panic(metaOpenErr)
	}

	defer metadataFile.Close()

	writeMetadataFile(metadataFile, metadata)

	indexFile, indexOpenErr := os.OpenFile(getIndexPath(header.QueueName), os.O_WRONLY|os.O_CREATE, 0777)
	if indexOpenErr != nil {
		panic(indexOpenErr)
	}

	defer indexFile.Close()

	index := IndexData{1, 1, header.TailerOrPayloadSize, 0, timeNow}
	fmt.Printf("%#v\n", index)
	writeIndexFile(indexFile, index)

	dataFile, dataOpenErr := os.OpenFile(getSegmentPath(header.QueueName, 1), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if dataOpenErr != nil {
		panic(dataOpenErr)
	}

	defer dataFile.Close()

	dataInfo := DataInfo{header.TailerOrPayloadSize, timeNow}

	writeDataFile(dataFile, dataInfo, reader)

}

func writeExistQueue(header Header, reader *bufio.Reader) {
	metadata := Metadata{}
	metadataFile, metaOpenErr := os.OpenFile(getMetadataPath(header.QueueName), os.O_RDWR, 0777)

	if metaOpenErr != nil {
		panic(metaOpenErr)
	}

	timeNow := getTimeNow()

	readMetadataFile(metadataFile, &metadata)

	fmt.Printf("%#v\n", metadata)

	currentIndexOffsetWrite := metadata.CurrentOffsetWrite

	metadata.CurrentOffsetWrite = metadata.CurrentOffsetWrite + int64(DATA_HEADER_SIZE) + int64(header.TailerOrPayloadSize)
	isNewSegment := MAX_SEGMENT_SIZE-(metadata.CurrentOffsetWrite+int64(DATA_HEADER_SIZE)) < int64(header.TailerOrPayloadSize)

	if isNewSegment {
		metadata.CountSegment = metadata.CountSegment + 1
		metadata.CurrentSegment = metadata.CurrentSegment + 1
	}
	metadata.CountMessage = metadata.CountMessage + 1
	metadata.QueueDataSize = metadata.QueueDataSize + int64(header.TailerOrPayloadSize)
	indexFile, indexOpenErr := os.OpenFile(getIndexPath(header.QueueName), os.O_WRONLY|os.O_APPEND, 0777)

	if indexOpenErr != nil {
		panic(indexOpenErr)
	}

	defer indexFile.Close()

	fmt.Printf("%#v\n", metadata)

	index := IndexData{
		metadata.CountMessage,
		metadata.CurrentSegment,
		header.TailerOrPayloadSize,
		currentIndexOffsetWrite,
		timeNow}
	fmt.Printf("%#v\n", index)
	writeIndexFile(indexFile, index)

	dataFile, dataOpenErr := os.OpenFile(getSegmentPath(header.QueueName, metadata.CurrentSegment), os.O_WRONLY|os.O_APPEND, 0777)

	if dataOpenErr != nil {
		panic(dataOpenErr)
	}

	defer dataFile.Close()

	dataInfo := DataInfo{header.TailerOrPayloadSize, timeNow}

	writeDataFile(dataFile, dataInfo, reader)
	writeMetadataFile(metadataFile, metadata)
	defer metadataFile.Close()
}

func getMessage(header Header, writer *bufio.Writer) {

	metadata := Metadata{}
	metadataFile, errMfile := os.Open(getMetadataPath(header.QueueName))
	readMetadataFile(metadataFile, &metadata)

	if errMfile != nil {
		panic(errMfile)
	}
	defer metadataFile.Close()

	queuePath := getQueuePath(header.QueueName)
	timeNow := getTimeNow()

	if !isDirExists(queuePath) {
		writer.Write(collectErrorBuf("Queue not exist"))
		writer.Flush()
		return
	}

	tailerDir := getTailerDirectory(header.QueueName)
	os.MkdirAll(tailerDir, 0777)

	tailer := Tailer{}
	tailerPath := getTailerPath(header.QueueName, header.TailerName)
	tailerFile, errTfile := os.OpenFile(tailerPath, os.O_RDWR|os.O_CREATE, 0777)

	if errTfile != nil {
		panic(errTfile)
	}
	defer tailerFile.Close()

	tailerFileInfo, errStat := os.Stat(tailerPath)

	if errStat != nil {
		panic(errStat)
	}

	if tailerFileInfo.Size() == 0 {
		tailer.Messageid = 0
		tailer.LastReadTime = timeNow
		fmt.Println("tailer пустой")
	} else {
		readTailerFile(tailerFile, &tailer)
	}

	fmt.Printf("%#v\n", metadata)

	if int64(metadata.CountMessage) == tailer.Messageid {
		fmt.Println("New messages not exist")
		writer.Write(collectErrorBuf("New messages not exist"))
		writer.Flush()
		return
	}

	index := IndexData{}

	indexFile, errInfile := os.OpenFile(getIndexPath(header.QueueName), os.O_RDWR, 0777)

	if errInfile != nil {
		panic(errInfile)
	}

	defer indexFile.Close()

	offset := int64(INDEX_RECORD_SIZE) * (tailer.Messageid)
	fmt.Printf("Смещение : %d \n", offset)
	readIndexFile(indexFile, &index, offset)

	dataFile, errDfile := os.OpenFile(getSegmentPath(header.QueueName, index.SegmentId), os.O_RDONLY, 0777)

	if errDfile != nil {
		panic(errDfile)
	}

	defer dataFile.Close()

	fmt.Println("До чтения")

	readDataFile(dataFile, index.OffsetInData, int64(index.Size), writer)

	tailer.Messageid = tailer.Messageid + 1
	tailer.LastReadTime = timeNow
	writeTailerFile(tailerFile, tailer)

	fmt.Printf("%#v\n", tailer)
	fmt.Printf("%#v\n", index)
	fmt.Println("Данные прочитаны")
}

func getStat(header Header, writer *bufio.Writer) {
	metadata := Metadata{}
	metadataFile, errMfile := os.Open(getMetadataPath(header.QueueName))
	readMetadataFile(metadataFile, &metadata)

	if errMfile != nil {
		panic(errMfile)
	}
	defer metadataFile.Close()

	queuePath := getQueuePath(header.QueueName)

	if !isDirExists(queuePath) {
		writer.Write(collectErrorBuf("Queue not exist"))
		writer.Flush()
		return
	}

	tailerDir := getTailerDirectory(header.QueueName)
	os.MkdirAll(tailerDir, 0777)

	tailer := Tailer{0, 0}
	tailerPath := getTailerPath(header.QueueName, header.TailerName)
	tailerFile, errTfile := os.OpenFile(tailerPath, os.O_RDWR|os.O_CREATE, 0777)

	tailerFileInfo, errStat := os.Stat(tailerPath)

	if errStat != nil {
		panic(errStat)
	}

	if tailerFileInfo.Size() == 0 {
		fmt.Println("tailer пустой")
	} else {
		readTailerFile(tailerFile, &tailer)
	}

	if errTfile != nil {
		panic(errTfile)
	}
	defer tailerFile.Close()

	indexFirstMessage := IndexData{}
	indexLastMessage := IndexData{}

	indexFilePath := getIndexPath(header.QueueName)

	indexFile, errInfile := os.OpenFile(indexFilePath, os.O_RDWR, 0777)

	if errInfile != nil {
		panic(errInfile)
	}
	defer indexFile.Close()

	indexFileInfo, _ := os.Stat(indexFilePath)

	readIndexFile(indexFile, &indexFirstMessage, 0)
	readIndexFile(indexFile, &indexLastMessage, indexFileInfo.Size()-int64(INDEX_RECORD_SIZE))

	stat := QueueStat{}

	stat.CountMessage = metadata.CountMessage
	stat.CountSegment = metadata.CountSegment
	stat.FirstWriteTime = indexFirstMessage.Time
	stat.LastWriteTime = indexLastMessage.Time
	stat.TailerCountMessage = int32(tailer.Messageid)
	stat.TailerLastTime = tailer.LastReadTime
	stat.QueueDataSize = metadata.QueueDataSize

	writer.Write(queueStatToBuf(stat))
	writer.Flush()
}

func checkNewMessage(header Header, writer *bufio.Writer) {
	metadata := Metadata{}
	metadataFile, errMfile := os.Open(getMetadataPath(header.QueueName))
	readMetadataFile(metadataFile, &metadata)
	if errMfile != nil {
		panic(errMfile)
	}
	defer metadataFile.Close()

	queuePath := getQueuePath(header.QueueName)

	if !isDirExists(queuePath) {
		writer.Write(collectErrorBuf("Queue not exist"))
		writer.Flush()
		return
	}

	tailerPath := getTailerPath(header.QueueName, header.TailerName)
	if !IsFileExists(tailerPath) {
		writer.Write(collectErrorBuf("Tailer not exist"))
		writer.Flush()
		return
	}

	tailer := Tailer{0, 0}
	tailerFile, errTfile := os.OpenFile(tailerPath, os.O_RDONLY, 0777)
	if errTfile != nil {
		panic(errTfile)
	}
	defer tailerFile.Close()

	readTailerFile(tailerFile, &tailer)

	resp := byte(0x01)

	if int64(metadata.CountMessage) == tailer.Messageid {
		fmt.Println("New messages not exist")
		resp = byte(0x00)
	}

	writer.Write(newMessageCheckBuf(resp))
	writer.Flush()

}
