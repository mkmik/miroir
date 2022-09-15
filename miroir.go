package miroir

import (
	"bytes"
	"io"
	"sync"
)

const (
	LEFT = 1 << iota
	RIGHT
)

// New returns a pair of readers that will both independently read from the input reader.
//
// If r implements ReadCloser upstream reader will be closed when both returned readers are closed.
// If the upstream reader returns an error on close, it will be visible only when the last reader
// is closed (the first Close call issued against any of the two readers will always return `nil`).
func New(r io.Reader) (io.ReadCloser, io.ReadCloser) {
	buf := bytes.NewBuffer(nil)
	m := &miroir{
		raw: r,
		r:   io.TeeReader(r, buf),
		buf: buf,
	}
	return &reader{m: m, id: LEFT}, &reader{m: m, id: RIGHT}
}

type miroir struct {
	raw io.Reader // so we can close the upstream reader
	r   io.Reader
	buf *bytes.Buffer
	sync.Mutex
	closed int
}

func (m *miroir) read(start int, p []byte) (int, error) {
	m.Lock()
	defer m.Unlock()

	if start < m.buf.Len() {
		n, err := bytes.NewReader(m.buf.Bytes()).ReadAt(p, int64(start))
		if err == io.EOF {
			// consume the rest from the upstream reader
			if n < len(p) {
				n2, err := m.r.Read(p[n:])
				return n + n2, err
			}
			err = nil
		}
		return n, err
	} else {
		return m.r.Read(p)
	}
}

func (m *miroir) close(id int) error {
	m.Lock()
	defer m.Unlock()

	m.closed |= id
	if m.closed == LEFT|RIGHT {
		if c, ok := m.raw.(io.ReadCloser); ok {
			return c.Close()
		}
	}
	return nil
}

type reader struct {
	m   *miroir
	pos int
	id  int
}

func (r *reader) Read(p []byte) (int, error) {
	n, err := r.m.read(r.pos, p)
	r.pos += n
	return n, err
}

func (r *reader) Close() error {
	return r.m.close(r.id)
}
