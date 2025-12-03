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
