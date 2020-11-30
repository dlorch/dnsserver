package main

import (
	"encoding/binary"
	"io"
)

func Write(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}
