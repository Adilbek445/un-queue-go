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

	metadata := Metadata{1, 2, 3, 4}

	writeMetadataFile(file, &metadata)

	metadataTest := Metadata{}
	readMetadataFile(file, &metadataTest)

	fmt.Println(metadata.CountMessage)
	fmt.Println(metadataTest.CountMessage)

	if err != nil {
		panic(err)
	}
}
