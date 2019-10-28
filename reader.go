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
	if br, ok := r.(*Reader); ok {
		return br
	}
	return &Reader{
		r: r,
	}
}

// fillBits makes sure that at least one bit is in the bits buffer, otherwise it will try to read a new byte into the buffer.
func (r *Reader) fillBits() error {
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

// ReadBits8 reads n bits from the source into the returned byte. Will not read more than 8 bits.
func (r *Reader) ReadBits8(n uint) (uint8, error) {
	if n > 8 {
		n = 8
	}

	// make sure the buffer is filled
	if err := r.fillBits(); err != nil {
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
	o, err := r.ReadBits8(rem)
	return res | o, err
}

func (r *Reader) Read(p []byte) (int, error) {
	var (
		i   = -1
		err error
	)
	for i = range p {
		p[i], err = r.ReadBits8(8)
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
		p[c], err = r.ReadBits8(n)
		c++
	}

	return c, err
}

// Bit reads a single bit from the reader and returns it as a bool (1 == true, 0 == false).
func (r *Reader) Bit() (bool, error) {
	b, err := r.ReadBits8(1)
	return b == 1, err
}

// BitsBuffered returns the number of bits that are currently buffered in the reader. This value is valid until the next read.
func (r *Reader) BitsBuffered() uint {
	return r.len
}
