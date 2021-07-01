package lz4

import (
	"sync"
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"fmt"
	"github.com/pierrec/lz4"
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
func BenchmarkLz4Compress128(b *testing.B) { benchmarkLz4Compress(b, s128) }
func BenchmarkLz4Compress256(b *testing.B) { benchmarkLz4Compress(b, s256) }
func BenchmarkLz4Compress512(b *testing.B)  { benchmarkLz4Compress(b, s512) }
func BenchmarkLz4Compress1024(b *testing.B)   { benchmarkLz4Compress(b, s1024) }
func BenchmarkLz4Compress2048(b *testing.B)   { benchmarkLz4Compress(b, s2048) }
func BenchmarkLz4Compress4096(b *testing.B)   { benchmarkLz4Compress(b, s4096) }
func BenchmarkLz4Compress8192(b *testing.B)   { benchmarkLz4Compress(b, s8192) }
func BenchmarkLz4Compress16384(b *testing.B)   { benchmarkLz4Compress(b, s16384) }
func BenchmarkLz4Compress32768(b *testing.B)   { benchmarkLz4Compress(b, s32768) }

func BenchmarkLz4UnCompress128(b *testing.B) { benchmarkLz4Uncompress(b, s128) }
func BenchmarkLz4UnCompress256(b *testing.B) { benchmarkLz4Uncompress(b, s256) }
func BenchmarkLz4UnCompress512(b *testing.B) { benchmarkLz4Uncompress(b, s512) }
func BenchmarkLz4UnCompress1024(b *testing.B) { benchmarkLz4Uncompress(b, s1024) }
func BenchmarkLz4UnCompress2048(b *testing.B) { benchmarkLz4Uncompress(b, s2048) }
func BenchmarkLz4UnCompress4096(b *testing.B) { benchmarkLz4Uncompress(b, s4096) }
func BenchmarkLz4UnCompress8192(b *testing.B) { benchmarkLz4Uncompress(b, s8192) }
func BenchmarkLz4UnCompress16384(b *testing.B) { benchmarkLz4Uncompress(b, s16384) }
func BenchmarkLz4UnCompress32768(b *testing.B) { benchmarkLz4Uncompress(b, s32768) }

func benchmarkLz4Compress(b *testing.B, uncompressed []byte) {
	/* ------------------官方方法---------------------- */
	buf := bytes.NewBuffer(nil)
	zw := lz4.NewWriter(buf)
	r := bytes.NewReader(uncompressed)

	// Determine the compressed size of testfile.
	compressedSize, err := io.Copy(zw, r)
	// compressedSize, err := zw.Write(uncompressed)
	if err != nil {
		b.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		b.Fatal(err)
	}

	// b.SetBytes(int64(compressedSize))
	b.SetBytes(compressedSize)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Reset(uncompressed)
		// buf.Reset()
		zw.Reset(buf)
		_, _ = io.Copy(zw, r)
		// _, _ = zw.Write(uncompressed)
	}

	fmt.Println(len(uncompressed), len(buf.Bytes()))
	/* -----------------------改造后的方法----------------------- */
	// s := "测试压缩解压时的消耗"
 //    r := []byte{}
 //    s := uncompressed
 //    client := NewClient(string(s))
    
 //    b.SetBytes(int64(len(s)))
	// b.ReportAllocs()
	// b.ResetTimer()

 //    for i := 0; i < b.N; i++ {
 //        r, _ = client.compress(s)
 //    }

 //    fmt.Println(len(s), len(r))
}

func benchmarkLz4Uncompress(b *testing.B, s []byte) {
	/* ------------------官方方法---------------------- */
    client := NewClient(string(s))
    compressed, _ := client.compress(s)
    // br := []byte{}

	r := bytes.NewReader(compressed)
	zr := lz4.NewReader(r)
	
	var buf bytes.Buffer
	// Determine the uncompressed size of testfile.
	uncompressedSize, err := io.Copy(&buf, zr)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(uncompressedSize)
	b.ReportAllocs()
	b.ResetTimer()

	// for i := 0; i < b.N; i++ {
	// 	br, _ = client.uncompress(compressed)
	// }
	// fmt.Println(len(s), len(compressed), len(br))

	for i := 0; i < b.N; i++ {
		r.Reset(compressed)
		zr.Reset(r)
		buf.Reset()
		_, _ = io.Copy(&buf, zr)
	}
	// for i := 0; i < b.N; i++ {
	// 	r.Reset(compressed)
	// 	// zr.Reset(r)
	// 	zr := lz4.NewReader(r)
	// 	buf.Reset()
	// 	_, _ = io.Copy(&buf, zr)
	// }

	fmt.Println(len(s), len(compressed), len(buf.Bytes()))
}

func mustLoadFile(f string) []byte {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return b
}

// var ErrShortWrite = errors.New("short write")
// var errInvalidWrite = errors.New("invalid write result")
// var EOF = errors.New("EOF")

type Client struct {
	compressWriters sync.Pool
	compressReaders sync.Pool
	bytesReaders sync.Pool
}

func NewClient(cluster string) *Client {
	client := &Client{}
	client.compressWriters.New = func() interface{} {
		var buf bytes.Buffer
		// writer, _ := lz4.NewWriterLevel(&buf, 1)
		// return writer
		return lz4.NewWriter(&buf)
	}
	client.bytesReaders.New = func() interface{} {
		// var buf bytes.Buffer
		// writer, _ := lz4.NewWriterLevel(&buf, 1)
		// return writer
		return bytes.NewReader(nil)
	}
	client.compressReaders.New = func() interface{} {
		// var buf bytes.Buffer
		// writer, _ := lz4.NewWriterLevel(&buf, 1)
		// return writer
		return lz4.NewReader(nil)
	}
	return client
}

// lz4 压缩
func (c *Client) compress(src []byte) ([]byte, error) {
	w := c.compressWriters.Get().(*lz4.Writer)
	defer func() {
		c.compressWriters.Put(w)
	}()
	var buf bytes.Buffer
	w.Reset(&buf)
	if _, err := w.Write(src); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// lz4 解缩
func (c *Client) uncompress(src []byte) ([]byte, error) {
	r := c.bytesReaders.Get().(*bytes.Reader)
	zr := c.compressReaders.Get().(*lz4.Reader)
	defer func() {
		c.compressReaders.Put(zr)
		c.bytesReaders.Put(r)
	}()
	// r := bytes.NewReader(src)
	r.Reset(src)
	zr.Reset(r)
	// 	r.Reset(compressed)
	// 	zr.Reset(r)
	// 	buf.Reset()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zr); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

	// r := lz4.NewReader(bytes.NewReader(src))
	// // if err != nil {
	// // 	return nil, err
	// // }
	// var buf bytes.Buffer
	// if _, err := io.Copy(&buf, r); err != nil {
	// 	return nil, err
	// }
	// // if err := r.Close(); err != nil {
	// // 	return nil, err
	// // }
	// return buf.Bytes(), nil
}
