package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

var METADATA_FILE string = "metadata.mt"
var INDEX_DATA_FILE string = "index.idx"
var DATA_FILE_EXTENSION string = ".dat"
var TAILER_EXTENSION string = ".trl"
var DATA_FILE_PREFIX string = "data-"

var currentDirectory string = ""

func main() {

	file, err := os.Create(getCurrentDirectory() + "/" + "my" + "/" + "metadata.mt")
	defer file.Close()

	if err != nil {
		panic(err)
	}

}

func write(file *os.File) {

	buf := new(bytes.Buffer)

	if file == nil {
		fmt.Println("SSSSSSSSSS")
	}

	// Пишем в буфер в LittleEndian
	binary.Write(buf, binary.LittleEndian, 123123)
	binary.Write(buf, binary.LittleEndian, 3434343)

	// Один вызов Write для файла
	file.Write(buf.Bytes())
}

func getSegmentName(segmentId int32) string {
	segmentString := fmt.Sprintf("%08d", segmentId)
	return currentDirectory + DATA_FILE_PREFIX + segmentString + DATA_FILE_EXTENSION
}

func getCurrentDirectory() string {
	return "C:\\dev\\projects\\my\\Go-language\\un-queue-go\\resources\\"
}
