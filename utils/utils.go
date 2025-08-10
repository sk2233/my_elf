package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseNum(val string) (int, bool) {
	if len(val) == 0 {
		return 0, false
	}
	if val[0] == '0' && len(val) > 2 { // 非 10 进制
		base := 0
		switch val[len(val)-1] {
		case 'o':
			base = 8
		default:
			panic(fmt.Sprintf("unknown base %v", val[len(val)-1]))
		}
		res, err := strconv.ParseInt(val[1:len(val)-1], base, 64)
		HandleErr(err)
		return int(res), true
	}
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false
	}
	return int(res), true
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
