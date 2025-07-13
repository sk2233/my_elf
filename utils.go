package main

import (
	"encoding/binary"
	"io"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsNum(val string) bool {
	if len(val) == 0 {
		return false
	}
	for i := 0; i < len(val); i++ {
		if val[i] < '0' || val[i] > '9' {
			return false
		}
	}
	return true
}

func WriteAny(w io.Writer, data interface{}) {
	err := binary.Write(w, binary.LittleEndian, data)
	HandleErr(err)
}

func FromU32(val uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, val)
	return res
}

func FromU64(val uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, val)
	return res
}

func WriteByte(w io.Writer, data []byte) {
	_, err := w.Write(data)
	HandleErr(err)
}
