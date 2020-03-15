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
	m *monitor.SpeedMonitor) error {

	newWriter, newReader := w, r
	var err error

	if compress {
		newWriter = zlib.NewWriter(w)
	}

	if decompress {
		if newReader, err = zlib.NewReader(r); err != nil {
			return err
		}
	}

	data := make([]byte, PACKSIZE)
	for {
		n, err := newReader.Read(data)
		if err != nil {
			return err
		}

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

	return nil
}

//Copy one packet only, for http1.0
func CopyOne(w io.Writer, r io.Reader, 
	compress, decompress bool, 
	m *monitor.SpeedMonitor) error {

	newWriter, newReader := w, r
	var err error

	if compress {
		newWriter = zlib.NewWriter(w)
	}

	if decompress {
		if newReader, err = zlib.NewReader(r); err != nil {
			return err
		}
		
	}

	data := make([]byte, PACKSIZE)

	n, err := newReader.Read(data)
	if err != nil {
		return err
	}

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

	return nil
}
