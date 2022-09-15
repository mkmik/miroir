package miroir

import (
	"bytes"
	"io"
	"sync"
)

// NewMiroir returns a pair of readers that will both independently read from the input reader.
func NewMiroir(r io.Reader) (io.Reader, io.Reader) {
	buf := bytes.NewBuffer(nil)
	m := &miroir{
		r:   io.TeeReader(r, buf),
		buf: buf,
	}
	return &reader{m: m}, &reader{m: m}
}

type miroir struct {
	r   io.Reader
	buf *bytes.Buffer
	sync.Mutex
}

func (m *miroir) read(start int, p []byte) (int, error) {
	m.Lock()
	defer m.Unlock()

	if start < m.buf.Len() {
		return m.buf.Read(p)
	} else {
		return m.r.Read(p)
	}
}

type reader struct {
	m   *miroir
	pos int
}

func (r *reader) Read(p []byte) (int, error) {
	n, err := r.m.read(r.pos, p)
	r.pos += n
	return n, err
}
