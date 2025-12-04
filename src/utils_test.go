package main

import (
	"encoding/binary"
	"testing"
)

func TestCollectErrorBuf(t *testing.T) {
	test := "hello"
	buf := collectErrorBuf(test)

	headerByte := byte(0x02)

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(len(test)))

	if headerByte != buf[0] {
		t.Error("разные заголовки")
	}

	lenght := int16(binary.BigEndian.Uint16(buf[1:3]))

	if int16(len(test)) != lenght {
		t.Error("разный размер сообщения")
	}

	text := string(buf[3 : 3+len(test)])

	if test != text {
		t.Error("разный message")
	}

}

func TestIsAlphanumericASCII(t *testing.T) {

	if IsAlphanumericASCII("test") == false {
		t.Error("Текст не удовлетворяет условию")
	}

	if IsAlphanumericASCII("тест") != false {
		t.Error("Текст не удовлетворяет условию")
	}

}

func TestQueueStatToBuf(t *testing.T) {
	stat := QueueStat{1, 1, 1, 1, 1, 1, 1}

	buf := queueStatToBuf(stat)

	countSegment := binary.BigEndian.Uint32(buf[1:5])
	countMessage := binary.BigEndian.Uint32(buf[5:9])
	firstWriteTime := binary.BigEndian.Uint64(buf[9:17])
	lastWriteTime := binary.BigEndian.Uint64(buf[17:25])
	tailerCountMessage := binary.BigEndian.Uint32(buf[25:29])
	tailerLastTime := binary.BigEndian.Uint64(buf[29:37])
	queueDataSize := binary.BigEndian.Uint64(buf[37:45])

	if buf[0] != byte(0x01) ||
		countSegment != uint32(stat.CountMessage) ||
		countMessage != uint32(stat.CountMessage) ||
		firstWriteTime != uint64(stat.FirstWriteTime) ||
		lastWriteTime != uint64(stat.LastWriteTime) ||
		tailerCountMessage != uint32(stat.TailerCountMessage) ||
		tailerLastTime != uint64(stat.TailerLastTime) ||
		queueDataSize != uint64(stat.QueueDataSize) {
		t.Error("Неправильная структура QueueStat")
	}

}
