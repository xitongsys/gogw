package common

import (
	"io"
	"net/http"
)

type Conn interface {
	io.Reader
	io.Writer
	io.Closer
}

type HttpConn struct {
	Reader *http.Request
	Writer http.ResponseWriter
}

func NewHttpConn(r *http.Request, w http.ResponseWriter) *HttpConn {
	return &HttpConn{
		Reader: r,
		Writer: w,
	}
}

func (hc *HttpConn) Read(bs []byte) (int ,error){
	return hc.Reader.Body.Read(bs)
}

func (hc *HttpConn) Write(bs []byte) (int, error){
	return hc.Writer.Write(bs)
}

func (hc *HttpConn) Close() error {
	return hc.Reader.Body.Close()
}