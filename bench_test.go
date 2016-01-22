// Copyright (C) 2014 Jakob Borg. All rights reserved. Use of this source code
// is governed by an MIT-style license that can be found in the LICENSE file.

package xdr_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/calmh/xdr"
)

type XDRBenchStruct struct {
	I1  uint64
	I2  uint32
	I3  uint16
	I4  uint8
	Bs0 []byte // max:128
	Bs1 []byte
	S0  string // max:128
	S1  string
}

var res []byte // no to be optimized away
var s = XDRBenchStruct{
	I1:  42,
	I2:  43,
	I3:  44,
	I4:  45,
	Bs0: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
	Bs1: []byte{11, 12, 13, 14, 15, 16, 17, 18, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	S0:  "Hello World! String one.",
	S1:  "Hello World! String two.",
}
var e []byte

func init() {
	e, _ = s.MarshalXDR()
}

func BenchmarkThisMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, _ = s.MarshalXDR()
	}

	b.ReportAllocs()
}

func BenchmarkThisUnmarshal(b *testing.B) {
	var t XDRBenchStruct
	for i := 0; i < b.N; i++ {
		err := t.UnmarshalXDR(e)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkThisEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := s.EncodeXDR(ioutil.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkThisEncoder(b *testing.B) {
	w := xdr.NewWriter(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		_, err := s.EncodeXDRInto(w)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

type repeatReader struct {
	data []byte
}

func (r *repeatReader) Read(bs []byte) (n int, err error) {
	if len(bs) > len(r.data) {
		err = io.EOF
	}
	n = copy(bs, r.data)
	r.data = r.data[n:]
	return n, err
}

func (r *repeatReader) Reset(bs []byte) {
	r.data = bs
}

func BenchmarkThisDecode(b *testing.B) {
	rr := &repeatReader{e}
	var t XDRBenchStruct
	for i := 0; i < b.N; i++ {
		err := t.DecodeXDR(rr)
		if err != nil {
			b.Fatal(err)
		}
		rr.Reset(e)
	}

	b.ReportAllocs()
}

func BenchmarkThisDecoder(b *testing.B) {
	rr := &repeatReader{e}
	r := xdr.NewReader(rr)
	var t XDRBenchStruct
	for i := 0; i < b.N; i++ {
		err := t.DecodeXDRFrom(r)
		if err != nil {
			b.Fatal(err)
		}
		rr.Reset(e)
	}

	b.ReportAllocs()
}

func BenchmarkReadString(b *testing.B) {
	buf := new(bytes.Buffer)
	xw := xdr.NewWriter(buf)
	orig := "This is short string for benchmarking purposes"
	xw.WriteString(orig)
	bs := buf.Bytes()

	r := &repeatReader{bs}
	xr := xdr.NewReader(r)

	var s string
	for i := 0; i < b.N; i++ {
		s = xr.ReadString()
		r.Reset(bs)
	}

	if s != orig {
		b.Fatalf("Wrong result, got %q instead of %q", s, orig)
	}
	b.ReportAllocs()
}

func BenchmarkWriteString(b *testing.B) {
	xw := xdr.NewWriter(ioutil.Discard)
	orig := "This is short string for benchmarking purposes"

	for i := 0; i < b.N; i++ {
		xw.WriteString(orig)
	}

	b.ReportAllocs()
}

func BenchmarkWriteBytes(b *testing.B) {
	xw := xdr.NewWriter(ioutil.Discard)
	orig := []byte("This is short string for benchmarking purposes")

	for i := 0; i < b.N; i++ {
		xw.WriteBytes(orig)
	}

	b.ReportAllocs()
}

func BenchmarkWriteUint32(b *testing.B) {
	xw := xdr.NewWriter(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		xw.WriteUint32(0x01020304)
	}

	b.ReportAllocs()
}

func BenchmarkWriteUint64(b *testing.B) {
	xw := xdr.NewWriter(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		xw.WriteUint64(0x0102030405060708)
	}

	b.ReportAllocs()
}

func BenchmarkReadBytesMax(b *testing.B) {
	buf := new(bytes.Buffer)
	xw := xdr.NewWriter(buf)
	orig := "This is short string for benchmarking purposes"
	xw.WriteString(orig)
	bs := buf.Bytes()

	r := &repeatReader{bs}
	xr := xdr.NewReader(r)

	var s []byte
	for i := 0; i < b.N; i++ {
		s = xr.ReadBytesMax(64)
		r.Reset(bs)
	}

	if string(s) != orig {
		b.Fatalf("Wrong result, got %q instead of %q", s, orig)
	}
	b.ReportAllocs()
}

func BenchmarkReadUint32(b *testing.B) {
	bs := []byte{1, 2, 3, 4}
	want := uint32(0x01020304)

	r := &repeatReader{bs}
	xr := xdr.NewReader(r)

	var s uint32
	for i := 0; i < b.N; i++ {
		s = xr.ReadUint32()
		r.Reset(bs)
	}

	if s != want {
		b.Fatalf("Wrong result, got %d instead of %d", s, want)
	}
	b.ReportAllocs()
}

func BenchmarkReadUint64(b *testing.B) {
	bs := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	want := uint64(0x0102030405060708)

	r := &repeatReader{bs}
	xr := xdr.NewReader(r)

	var s uint64
	for i := 0; i < b.N; i++ {
		s = xr.ReadUint64()
		r.Reset(bs)
	}

	if s != want {
		b.Fatalf("Wrong result, got %d instead of %d", s, want)
	}
	b.ReportAllocs()
}
