package bitio

import (
	"bytes"
	"testing"
)

func TestNewReader(t *testing.T) {
	r1 := NewReader(nil)
	r2 := NewReader(r1)

	if r1 != r2 {
		t.Error("Expected r1 and r2 to be equal")
	}
}

func TestBool(t *testing.T) {
	results := []bool{true, false, true, false, true, false, true, false, true, true, true, true, true, true, true, true}
	r := NewReader(bytes.NewReader([]byte{0b10101010, 0b11111111}))

	for i, e := range results {
		b, err := r.Bool()
		if err != nil {
			t.Errorf("Error in bit #%d: %s", i, err)
		}
		if e != b {
			t.Errorf("Mismatch in bit #%d: expected %t != got %t", i, e, b)
		}
	}

	_, err := r.Bool()
	if err == nil {
		t.Error("Expected error at end of input data, got none!")
	}
}

func TestReadBits8(t *testing.T) {
	results := []byte{0b101, 0b010, 0b10111, 0b1, 0b1111}
	args := []uint{3, 3, 5, 1, 4}
	r := NewReader(bytes.NewReader([]byte{200, 0b10101010, 0b11111111}))

	b, err := r.ReadUint8(42)
	if err != nil {
		t.Error("Expected no error at start of input data")
	}
	if b != 200 {
		t.Errorf("Expected result 200, got %d", b)
	}

	for i, e := range results {
		b, err = r.ReadUint8(args[i])
		if err != nil {
			t.Errorf("Error in read #%d: %s", i, err)
		}
		if e != b {
			t.Errorf("Mismatch in read #%d: expected %d != got %d", i, e, b)
		}
	}

	_, err = r.ReadUint8(4)
	if err == nil {
		t.Error("Expected error at end of input data, got none!")
	}
}
