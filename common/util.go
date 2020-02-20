package common

import (
	"fmt"
	"io"

	"github.com/vmihailenco/msgpack/v4"
)

var UUIDMAP map[string]int = make(map[string]int)

func UUID(key string) string {
	if _, ok := UUIDMAP[key]; !ok {
		UUIDMAP[key] = 0
	}

	UUIDMAP[key]++
	return fmt.Sprint(UUIDMAP[key])
}

func ReadObject(r io.Reader, o interface{}) (err error) {
	lenBuf = []byte{0}
	data := []byte{}
	var n int 
	for {
		if _, err = io.ReadAtLeast(r, lenBuf, 1); err != nil {
			return err
		}

		l := int(lenBuf[0])
		if l == 0 {
			err = msgpack.Unmarshal(data, o)
			break
		}

		buf := make([]byte, l)
		if n, err = io.ReadAtLeast(r, buf, l); err != nil {
			return err
		}

		data = append(data, buf...)
	}
} 

func WriteObject(w io.Writer, o interface{}) (err error){ 
	var data []byte
	if data, err = msgpack.Marshal(o); err != nil {
		return err
	}

	l := len(data)
	for l > 0 {
		cl := 255
		if cl > l {
			cl = l
		}

		if _, err = w.Write([]byte{byte(cl)}); err != nil {
			return err
		}

		if l == 0 {
			break
		}

		buf := data[:cl]
		if _, err = w.Write(buf); err != nil {
			return err
		}

		l -= cl
	}

}