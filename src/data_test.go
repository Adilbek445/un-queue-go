package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSum(t *testing.T) {

	dir := t.TempDir()

	metadataPath := dir + "\\" + METADATA_FILE

	file, err := os.Create(metadataPath)
	defer file.Close()

	if err != nil {
		panic(err)
	}
	metadata := Metadata{1, 2, 3, 4}

	writeMetadataFile(file, metadata)

	metadataTest := Metadata{}
	readMetadataFile(file, &metadataTest)

	fmt.Println(metadata.CountMessage)
	fmt.Println(metadataTest.CountMessage)

	if (metadata.CountMessage != metadataTest.CountMessage) &&
		(metadata.CountSegment != metadataTest.CountSegment) &&
		(metadata.CurrentSegment != metadataTest.CurrentSegment) &&
		(metadata.CurrentOffsetWrite != metadataTest.CurrentOffsetWrite) {
		t.Error("metadata не совпадают")
	}

}
