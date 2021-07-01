package zstd

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"fmt"
	// "github.com/pierrec/lz4"
	"github.com/klauspost/compress/zstd"
)

var (
	s128    = mustLoadFile("testdata/128.txt")
	s256    = mustLoadFile("testdata/256.txt")
	s512     = mustLoadFile("testdata/512.txt")
	s1024    = mustLoadFile("testdata/1024.txt")
	s2048 = mustLoadFile("testdata/2048.txt")
	s4096 = mustLoadFile("testdata/4096.txt")
	s8192  = mustLoadFile("testdata/8192.txt")
	s16384 = mustLoadFile("testdata/16384.txt")
	s32768 = mustLoadFile("testdata/32768.txt")
)

// 测试内存：go test -v -bench=Compress -benchmem lz4_test.go
func BenchmarkZstdCompress128(b *testing.B) { benchmarkZstdCompress(b, s128) }
func BenchmarkZstdCompress256(b *testing.B) { benchmarkZstdCompress(b, s256) }
func BenchmarkZstdCompress512(b *testing.B)  { benchmarkZstdCompress(b, s512) }
func BenchmarkZstdCompress1024(b *testing.B)   { benchmarkZstdCompress(b, s1024) }
func BenchmarkZstdCompress2048(b *testing.B)   { benchmarkZstdCompress(b, s2048) }
func BenchmarkZstdCompress4096(b *testing.B)   { benchmarkZstdCompress(b, s4096) }
func BenchmarkZstdCompress8192(b *testing.B)   { benchmarkZstdCompress(b, s8192) }
func BenchmarkZstdCompress16384(b *testing.B)   { benchmarkZstdCompress(b, s16384) }
func BenchmarkZstdCompress32768(b *testing.B)   { benchmarkZstdCompress(b, s32768) }

func BenchmarkZstdUnCompress128(b *testing.B) { benchmarkZstdUncompress(b, s128) }
func BenchmarkZstdUnCompress256(b *testing.B) { benchmarkZstdUncompress(b, s256) }
func BenchmarkZstdUnCompress512(b *testing.B) { benchmarkZstdUncompress(b, s512) }
func BenchmarkZstdUnCompress1024(b *testing.B) { benchmarkZstdUncompress(b, s1024) }
func BenchmarkZstdUnCompress2048(b *testing.B) { benchmarkZstdUncompress(b, s2048) }
func BenchmarkZstdUnCompress4096(b *testing.B) { benchmarkZstdUncompress(b, s4096) }
func BenchmarkZstdUnCompress8192(b *testing.B) { benchmarkZstdUncompress(b, s8192) }
func BenchmarkZstdUnCompress16384(b *testing.B) { benchmarkZstdUncompress(b, s16384) }
func BenchmarkZstdUnCompress32768(b *testing.B) { benchmarkZstdUncompress(b, s32768) }

func benchmarkZstdCompress(b *testing.B, uncompressed []byte) {
	buf := bytes.NewBuffer(nil)
	zw, _ := zstd.NewWriter(buf)
	// zw, _ := zstd.NewWriter(buf, zstd.WithEncoderLevel(zstd.SpeedFastest))
	r := bytes.NewReader(uncompressed)

    compressedSize, _ := io.Copy(zw, r)
    zw.Close()
    // fmt.Println(len(uncompressed), len(buf.Bytes()))

    b.SetBytes(compressedSize)
	b.ReportAllocs()
	b.ResetTimer()
    
 //    for i := 0; i < b.N; i++ {
 //    	buf.Reset()
	// 	r.Reset(uncompressed)
	// 	zw.Reset(buf)
	// 	_, _ = io.Copy(zw, r)
	// 	zw.Close()
	// }

	// fmt.Println(len(uncompressed), len(buf.Bytes()))

	ret := []byte{}
    for i := 0; i < b.N; i++ {
    	ret = zw.EncodeAll(uncompressed, make([]byte, 0, len(uncompressed)))
	}

	fmt.Println(len(uncompressed), len(ret))
}

func benchmarkZstdUncompress(b *testing.B, s []byte) {
	buf := bytes.NewBuffer(nil)
	zw, _ := zstd.NewWriter(buf)
	// zw, _ := zstd.NewWriter(buf, zstd.WithEncoderLevel(zstd.SpeedFastest))
	rr := bytes.NewReader(s)

    _, _ = io.Copy(zw, rr)
    zw.Close()
    // client := NewClient(string(s))
    compressed := buf.Bytes()
    // fmt.Println(compressed)
    // br := []byte{}

	r := bytes.NewReader(compressed)
	zr, _ := zstd.NewReader(r)
	
	buf.Reset()
	// Determine the uncompressed size of testfile.
	uncompressedSize, err := io.Copy(buf, zr)
	// fmt.Println(uncompressedSize, len(buf.Bytes()))
	if err != nil {
		b.Fatal(err)
	}
	// zr.Close() //这个库重用reader的时候不能close

	b.SetBytes(uncompressedSize)
	b.ReportAllocs()
	b.ResetTimer()

	// for i := 0; i < b.N; i++ {
	// 	r.Reset(compressed)
	// 	zr.Reset(r)
	// 	buf.Reset()
	// 	_, _ = io.Copy(buf, zr)
	// 	// zr.Close()
	// }

	// fmt.Println(len(s), len(compressed), len(buf.Bytes()))

	ret := []byte{}
    for i := 0; i < b.N; i++ {
    	ret, _ = zr.DecodeAll(compressed, nil)
	}

	fmt.Println(len(s), len(compressed), len(ret))
}

func mustLoadFile(f string) []byte {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return b
}
