package common

import (
	"io"
	"net"
	"time"
)


type UDPConn struct {
	remoteAddr net.Addr
	conn *net.UDPConn	
	pipeReader *io.PipeReader
	PipeWriter *io.PipeWriter
}

func NewUDPConn(remoteAddr net.Addr, conn *net.UDPConn) (*UDPConn){
	r, w := io.Pipe()
	return & UDPConn {
		remoteAddr: remoteAddr,
		conn: conn,
		pipeReader: r,
		PipeWriter: w,
	}
}

func (uc *UDPConn) Read(b []byte) (n int , err error) {
	return uc.pipeReader.Read(b)
}

func (uc *UDPConn) Write(b []byte) (n int , err error) {
	return uc.conn.WriteTo(b, uc.remoteAddr)
}

func (uc *UDPConn) Close() error {
	return uc.conn.Close()
}

func (uc *UDPConn) LocalAddr() net.Addr {
	return uc.conn.LocalAddr()
}

func (uc *UDPConn) RemoteAddr() net.Addr {
	return uc.remoteAddr
}

func (uc *UDPConn) SetDeadline(t time.Time) error {
	return nil
}

func (uc *UDPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (uc *UDPConn) SetWriteDeadline(t time.Time) error {
	return nil
}