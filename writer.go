package bitio

import "io"

// Writer offers support for writes of arbitrary bit length.
type Writer struct {
	w    io.Writer
	free uint
	bits byte
}

// NewWriter creates a new bit writer. If w is already a bit writer it will return w.
func NewWriter(w io.Writer) *Writer {
	if bw, ok := w.(*Writer); ok {
		return bw
	}
	return &Writer{
		free: 8,
		w:    w,
	}
}

// Flush writes the remaining buffer into the writer. If the buffer is empty this is a no-op.
func (w *Writer) Flush() error {
	if w.free != 8 {
		return nil
	}
	_, err := w.w.Write([]byte{w.bits})
	w.free = 8
	return err
}

// WriteBits8 writes up to 8 bits into the writer. When a full byte is reached in the buffer it is written to the underlying writer.
func (w *Writer) WriteBits8(b byte, n uint) error {
	if n > 8 {
		n = 8
	}

	var rem uint
	if n > w.free {
		rem = n - w.free
		n = w.free
	}

	w.free -= n
	w.bits <<= n
	w.bits |= (b & (1<<n - 1))

	var err error
	if w.free == 0 {
		err = w.Flush()
	}
	if rem == 0 {
		return err
	}

	b >>= n
	return w.WriteBits8(b, rem)
}

func (w *Writer) Write(p []byte) (int, error) {
	for i, b := range p {
		err := w.WriteBits8(b, 8)
		if err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// WriteBits writes n bits from p. The function will not write more bits that can be contained in p.
func (w *Writer) WriteBits(p []byte, n uint) (int, error) {
	// make sure to only write the appropriate number of bits
	if maxLen := uint(len(p)) * 8; maxLen < n {
		n = maxLen
	}

	c, err := w.Write(p[:n/8])
	if err != nil {
		return c, err
	}

	m := n % 8
	if m > 0 {
		err = w.WriteBits8(p[n/8], m)
		c++
	}

	return c, err
}

// Bit writes a single bit into the writer.
func (w *Writer) Bit(b bool) error {
	var by byte
	if b {
		by = 1
	}
	return w.WriteBits8(by, 1)
}

// BitsBuffered returns the number of bits that are currently buffered in the writer. This value is valid until the next write.
func (w *Writer) BitsBuffered() uint {
	return 8 - w.free
}
