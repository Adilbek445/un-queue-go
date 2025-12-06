package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func getTimeNow() int64 {
	now := time.Now()
	return now.UnixMilli()
}

func serializeHeader(buf [7]byte) Header {
	header := Header{}

	header.HeaderFlag = int8(buf[0])
	header.QueueNameSize = int16(binary.BigEndian.Uint16(buf[1:3]))
	fmt.Println("QueueNameSize %u", header.QueueNameSize)
	header.TailerOrPayloadSize = int32(binary.BigEndian.Uint32(buf[3:7]))
	return header
}

func collectErrorBuf(errorMessage string) []byte {
	header := byte(0x02)
	errorMessageSize := int16(len(errorMessage))

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(errorMessageSize))

	errorMessageByte := []byte(errorMessage)
	payload := []byte{header}
	payload = append(payload, b...)
	payload = append(payload, errorMessageByte...)
	return payload
}

func queueStatToBuf(stat QueueStat) []byte {
	b := make([]byte, 44)
	binary.BigEndian.PutUint32(b[0:4], uint32(stat.CountSegment))
	binary.BigEndian.PutUint32(b[4:8], uint32(stat.CountMessage))
	binary.BigEndian.PutUint64(b[8:16], uint64(stat.FirstWriteTime))
	binary.BigEndian.PutUint64(b[16:24], uint64(stat.LastWriteTime))
	binary.BigEndian.PutUint32(b[24:28], uint32(stat.TailerCountMessage))
	binary.BigEndian.PutUint64(b[28:36], uint64(stat.TailerLastTime))
	binary.BigEndian.PutUint64(b[36:44], uint64(stat.QueueDataSize))

	header := byte(0x01)

	payload := []byte{header}
	payload = append(payload, b...)
	return payload
}

func newMessageCheckBuf(newMessageExist byte) []byte {
	header := byte(0x01)
	payload := []byte{header}
	payload = append(payload, newMessageExist)
	return payload
}

func IsAlphanumericASCII(s string) bool {
	for _, r := range s {
		isLetterOrDigit := (r >= '0' && r <= '9') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z')

		if !isLetterOrDigit {
			return false
		}
	}

	return len(s) > 0
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
	return filepath.Join(getTailerDirectory(queueName), tailer+TAILER_EXTENSION)
}

func getTailerDirectory(queueName string) string {
	return filepath.Join(currentDirectory, queueName, "tailer")
}
