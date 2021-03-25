package logger

import (
	"io"
)

type Handler interface {
	io.Closer
	io.Writer
}

type StreamHandler struct {
	io.Writer
}

func NewStreamHandler(w io.Writer) (*StreamHandler, error) {
	return &StreamHandler{Writer:w}, nil
}

func (p *StreamHandler) Write(b []byte) (n int, err error) {
	return p.Writer.Write(b)
}

func (h *StreamHandler) Close() error {
	return nil
}

type NullHandler struct {
}

func NewNullHandler() (*NullHandler, error) {
	return &NullHandler{}, nil
}

func (p *NullHandler) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (p *NullHandler) Close() error {
	return nil
}