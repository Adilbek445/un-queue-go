package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func writeDataFile(file *os.File, data DataInfo, reader *bufio.Reader) {
	binary.Write(file, binary.LittleEndian, data.Size)
	binary.Write(file, binary.LittleEndian, data.Time)
	tee := io.TeeReader(reader, os.Stdout)
	_, err := io.Copy(file, tee)

	if err != nil {
		panic(err)
	}

}

func readDataFile(file *os.File, offset int64, size int64, writer *bufio.Writer) {
	offset = offset + int64(DATA_HEADER_SIZE)
	fmt.Println("READ SIZE:", size)
	_, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}

	logWriter := io.MultiWriter(writer, os.Stdout)

	_, err = io.CopyN(logWriter, file, size)
	if err != nil && err != io.EOF {
		panic(err)
	}

	writer.Flush()

}

func writeIndexFile(file *os.File, index IndexData) {
	err := binary.Write(file, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}
}

func readIndexFile(file *os.File, index *IndexData, offset int64) {
	file.Seek(offset, io.SeekStart)
	err := binary.Read(file, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}
}

func writeTailerFile(file *os.File, tailer Tailer) {
	file.Seek(0, io.SeekStart)
	err := binary.Write(file, binary.LittleEndian, tailer)
	if err != nil {
		panic(err)
	}
}

func readTailerFile(file *os.File, tailer *Tailer) {
	file.Seek(0, io.SeekStart)
	err := binary.Read(file, binary.LittleEndian, tailer)
	if err != nil {
		panic(err)
	}
}

func writeMetadataFile(file *os.File, metadata Metadata) {
	file.Seek(0, io.SeekStart)
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
