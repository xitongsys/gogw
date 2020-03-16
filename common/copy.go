package common

import (
	"io"
	"compress/zlib"
	"net/http"

	"gogw/monitor"
)

const (
	PACKSIZE = 1024 * 1024
)

func Copy(w io.Writer, r io.Reader, 
	compress, decompress bool, 
	m *monitor.SpeedMonitor) (int, error) {

	newWriter, newReader := w, r
	var err error

	if compress {
		newWriter = zlib.NewWriter(w)
	}

	if decompress {
		if newReader, err = zlib.NewReader(r); err != nil {
			return 0, err
		}
	}

	data := make([]byte, PACKSIZE)
	tot := 0
	for {
		n, err := newReader.Read(data)
		if n<=0 && err != nil {
			return tot, err
		}

		if n <= 0 {
			continue
		}

		tot += n

		n, err = newWriter.Write(data[:n])
		if err != nil {
			//return err
		}

		if ww, ok := newWriter.(*zlib.Writer); ok {
			ww.Flush()
		}

		if ww, ok := w.(http.Flusher); ok {
			ww.Flush()
		}

		//monitor
		if m != nil {
			m.Add(int64(n))
		}
	}

	return tot, nil
}

//Copy one packet only, for http1.0
func CopyOne(w io.Writer, r io.Reader, 
	compress, decompress bool, 
	m *monitor.SpeedMonitor) (int, error) {

	newWriter, newReader := w, r
	var err error
	var n int 

	if compress {
		newWriter = zlib.NewWriter(w)
	}

	if decompress {
		if newReader, err = zlib.NewReader(r); err != nil {
			return 0, err
		}
	}

	data := make([]byte, PACKSIZE)

	for n == 0 && err == nil {
		n, err = newReader.Read(data)
		if n <=0 && err != nil {
			return 0, err
		}

		if n <= 0 {
			continue
		}

		_, err = newWriter.Write(data[:n])
		if err != nil {
			return 0, err
		}

		if ww, ok := newWriter.(*zlib.Writer); ok {
			ww.Flush()
		}

		//monitor
		if m != nil {
			m.Add(int64(n))
		}
	}
		
	return n, err
}

//used in http1.0
func CopyAll(w io.Writer, r io.Reader, 
	compress, decompress bool, 
	m *monitor.SpeedMonitor) (int, error) {

	newWriter, newReader := w, r
	var err error

	if compress {
		newWriter = zlib.NewWriter(w)
	}

	if decompress {
		if newReader, err = zlib.NewReader(r); err != nil {
			return 0, err
		}
	}

	defer func(){
		if ww, ok := newWriter.(*zlib.Writer); ok {
			ww.Flush()
		}
	}()

	data := make([]byte, PACKSIZE)
	tot := 0
	for {
		n, err := newReader.Read(data)
		if n <= 0 && err != nil {
			return tot, err
		}

		if n <= 0 {
			continue
		}

		tot += n

		_, err = newWriter.Write(data[:n])
		if err != nil {
			//return err
		}

		//monitor
		if m != nil {
			m.Add(int64(n))
		}
	}

	return tot, nil
}