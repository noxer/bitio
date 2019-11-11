package bitio

import (
	"bytes"
	"testing"
)

func TestBitio(t *testing.T) {
	data := []byte{42, 66, 130, 255, 0, 0, 34, 89}
	lengths := []uint{3, 8, 1, 0, 7, 7, 1, 8, 1, 8, 4, 4, 4, 6, 2}
	res := make([]byte, len(lengths))
	var err error

	r := NewReader(bytes.NewReader(data))
	for i, l := range lengths {
		res[i], err = r.ReadUint8(l)
		if err != nil {
			t.Error(err)
		}
	}

	buf := &bytes.Buffer{}
	w := NewWriter(buf)
	for i, l := range lengths {
		err = w.WriteBits8(res[i], l)
		if err != nil {
			t.Error(err)
		}
	}

	if len(buf.Bytes()) != len(data) {
		t.Errorf("Slices don't match: %v vs. expected %v", buf.Bytes(), data)
	}
	for i, b := range buf.Bytes() {
		if data[i] != b {
			t.Errorf("Expected %d, got %d in %d", data[i], b, i)
		}
	}
}
