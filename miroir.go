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
		n, err := bytes.NewReader(m.buf.Bytes()).ReadAt(p, int64(start))
		//		log.Printf("GOT FROM BUFFER n=%d, err=%v", n, err)
		if err == io.EOF {
			//			log.Printf("GOT EOF IN BUFFER n=%d, len(p)=%d", n, len(p))
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

type reader struct {
	m   *miroir
	pos int
}

func (r *reader) Read(p []byte) (int, error) {
	n, err := r.m.read(r.pos, p)
	r.pos += n
	return n, err
}
