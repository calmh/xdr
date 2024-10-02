package xdr

import (
	"errors"
	"io"
	"testing"
)

func TestMissingPaddingUnmarshal(t *testing.T) {
	// A three-byte byte slice, missing the fourth padding byte
	v := []byte{0x00, 0x00, 0x00, 0x03, 0xCA, 0xFE, 0xFE}
	u := &Unmarshaller{Data: v}
	u.UnmarshalBytesMax(4)
	if !errors.Is(u.Error, io.ErrUnexpectedEOF) {
		t.Fatal(`Expected "Unexpected EOF", got`, u.Error)
	}
}
