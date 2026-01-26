package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Stat struct {
	QueueName    string `json:"queueName"`
	DataSize     int64  `json:"dataSize"`
	CountSegment int32  `json:"countSegment"`
	CountMessage int32  `json:"countMessage"`
	DataInfo     []Data `json:"dataInfo"`
}

type Data struct {
	Size    int64  `json:"size"`
	Time    string `json:"time"`
	Payload string `json:"payload"`
}

func getInfoByQueueName(queueName string) {
	path := getSegmentPath(queueName, 1)
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat := Stat{}

	stat.QueueName = queueName

	buffer := make([]byte, 12)
	var datas []Data

	for {

		_, err = io.ReadFull(file, buffer)

		if err == io.EOF {
			break
		}

		size := binary.LittleEndian.Uint32(buffer[0:4])
		date := binary.LittleEndian.Uint64(buffer[4:12])

		t := time.UnixMilli(int64(date))
		dateStringRFC3339 := t.Format(time.RFC3339)

		data := Data{}
		data.Size = int64(int64(size))
		data.Time = dateStringRFC3339

		dataBuffer := make([]byte, size)

		_, err = io.ReadFull(file, dataBuffer)

		if err == io.EOF {
			break
		}

		data.Payload = string(dataBuffer)
		datas = append(datas, data)
	}

	stat.DataInfo = datas

	jsonData, _ := json.Marshal(stat)

	outPath := currentDirectory + "/" + "out.txt"

	outFile, err := os.Create(outPath)

	outFile.Write(jsonData)
	defer outFile.Close()
	// fmt.Println(string(jsonData))
}
