package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestWriteAndReadMetadata(t *testing.T) {

	dir := t.TempDir()

	metadataPath := dir + "\\" + METADATA_FILE

	file, err := os.Create(metadataPath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	metadata := Metadata{1, 2, 3, 4, 5}

	writeMetadataFile(file, metadata)

	metadataTest := Metadata{}
	readMetadataFile(file, &metadataTest)

	fmt.Println(metadata.CountMessage)
	fmt.Println(metadataTest.CountMessage)

	if (metadata.CountMessage != metadataTest.CountMessage) ||
		(metadata.CountSegment != metadataTest.CountSegment) ||
		(metadata.CurrentSegment != metadataTest.CurrentSegment) ||
		(metadata.CurrentOffsetWrite != metadataTest.CurrentOffsetWrite) {
		t.Error("metadata не совпадают")
	}
}

func TestWriteAndReadIndex(t *testing.T) {

	dir := t.TempDir()

	indexPath := dir + "\\" + INDEX_DATA_FILE

	file, errInfile := os.Create(indexPath)

	if errInfile != nil {
		panic(errInfile)
	}

	defer file.Close()

	index := IndexData{1, 2, 3, 56012, 5}

	writeIndexFile(file, index)

	indexRead := IndexData{}
	readIndexFile(file, &indexRead, 0)

	if (index.MessageId != indexRead.MessageId) ||
		(index.SegmentId != indexRead.SegmentId) ||
		(index.Size != indexRead.Size) ||
		(index.Time != indexRead.Time) ||
		(index.OffsetInData != indexRead.OffsetInData) {
		t.Error("index не совпадают")
	}

}

func TestWriteAndReadData(t *testing.T) {

	dir := t.TempDir()

	segment := getSegmentName(50)

	dataPath := dir + "\\" + segment

	file, errDatafile := os.Create(dataPath)

	if errDatafile != nil {
		panic(errDatafile)
	}

	defer file.Close()

	dataInfo := DataInfo{1000, 20000}

	dataText := "hello world111"

	reader := bufio.NewReader(strings.NewReader(dataText))

	writeDataFile(file, dataInfo, reader)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	readDataFile(file, 0, int64(len(dataText)), writer)

	dataRead := buf.String()

	fmt.Println(dataRead)

	if dataRead != dataText {
		t.Error("data не совпадают")
	}
}
