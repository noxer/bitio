package bitio

import "io"

// Reader offers support for reads of arbitrary bit length.
type Reader struct {
	r    io.Reader
	len  uint
	bits byte
}

// NewReader creates a new bit reader for r. If r is already a bit reader it just returns r.
func NewReader(r io.Reader) *Reader {
	// just return when this is already the right reader
	if br, ok := r.(*Reader); ok {
		return br
	}
	return &Reader{
		r: r,
	}
}

// fillBuffer makes sure that at least one bit is in the bits buffer, otherwise it will try to read a new byte into the buffer.
func (r *Reader) fillBuffer() error {
	if r.len != 0 {
		return nil
	}

	var buf [1]byte
	_, err := r.r.Read(buf[:])
	if err != nil {
		return err
	}

	r.bits = buf[0]
	r.len = 8

	return nil
}

// ReadUint8 reads n bits from the source into the returned byte. Will not read more than 8 bits.
func (r *Reader) ReadUint8(n uint) (uint8, error) {
	if n > 8 {
		n = 8
	}

	// make sure the buffer is filled
	if err := r.fillBuffer(); err != nil {
		return 0, err
	}

	// fast path, we've got exactly the amount of bits needed to fulfill the request
	if r.len == n {
		r.len = 0
		return r.bits, nil
	}

	// make sure we're not trying to read too many bits
	var rem uint
	if n > r.len {
		rem = n - r.len
		n = r.len
	}

	r.len -= n
	res := r.bits >> r.len
	r.bits &= 1<<r.len - 1

	// we may be done here
	if rem == 0 {
		return res, nil
	}

	// shift the partial result and read the other part
	res <<= rem
	o, err := r.ReadUint8(rem)
	return res | o, err
}

func (r *Reader) Read(p []byte) (int, error) {
	var (
		i   = -1
		err error
	)
	for i = range p {
		p[i], err = r.ReadUint8(8)
		if err != nil {
			break
		}
	}

	return i + 1, err
}

// ReadBits reads n bits into p. The function will not read more bits that can be contained in p.
func (r *Reader) ReadBits(p []byte, n uint) (int, error) {
	// make sure to only read the appropriate number of bits
	if maxLen := uint(len(p)) * 8; maxLen < n {
		n = maxLen
	}

	// fill the buffer with the full bytes
	c, err := r.Read(p[:n/8])
	if err != nil {
		return c, err
	}

	n %= 8
	if n > 0 {
		p[c], err = r.ReadUint8(n)
		c++
	}

	return c, err
}

// Bool reads a single bit from the reader and returns it as a boolean.
func (r *Reader) Bool() (bool, error) {
	b, err := r.Bit()
	return b == 1, err
}

// Bit reads a single bit from the reader and returns it as a bool (1 == true, 0 == false).
func (r *Reader) Bit() (b byte, err error) {
	if r.len == 0 {
		err = r.fillBuffer()
		if err != nil {
			return
		}
	}

	r.len--
	b = r.bits >> r.len
	r.bits -= b << r.len

	return
}

// Bools reads n bits from the reader and returns them as a []bool.
func (r *Reader) Bools(n int) ([]bool, error) {
	bits := make([]bool, n)
	for i := 0; i < n; i++ {
		b, err := r.Bool()
		if err != nil {
			return bits, err
		}
		bits[i] = b
	}
	return bits, nil
}

// BitsBuffered returns the number of bits that are currently buffered in the reader. This value is valid until the next read.
func (r *Reader) BitsBuffered() uint {
	return r.len
}
