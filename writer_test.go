package bitio

import (
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	data := []byte{0b110011, 0b001111, 0b1, 0b11110000}
	lengths := []uint{6, 6, 1, 8}
	results := []byte{0b11001100, 0b11111111, 0b10000000}

	buf := bytes.Buffer{}
	w := NewWriter(&buf)
	for i, d := range data {
		err := w.WriteBits8(d, lengths[i])
		if err != nil {
			t.Error(err)
		}
	}
	err := w.Flush()
	if err != nil {
		t.Error(err)
	}

	if len(buf.Bytes()) != len(results) {
		t.Errorf("Results don't match: %v vs. expected %v", buf.Bytes(), results)
	}
	for i, b := range buf.Bytes() {
		if b != results[i] {
			t.Errorf("Unexpected result %d != expected %d in %d", b, results[i], i)
		}
	}
}
