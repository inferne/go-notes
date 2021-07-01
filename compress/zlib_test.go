package zlib

import (
	"testing"
	"bytes"
	"io"
	"io/ioutil"
	"compress/zlib"
	"fmt"
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
func BenchmarkZlibCompress128(b *testing.B) { benchmarkZlibCompress(b, s128) }
func BenchmarkZlibCompress256(b *testing.B) { benchmarkZlibCompress(b, s256) }
func BenchmarkZlibCompress512(b *testing.B)  { benchmarkZlibCompress(b, s512) }
func BenchmarkZlibCompress1024(b *testing.B)   { benchmarkZlibCompress(b, s1024) }
func BenchmarkZlibCompress2048(b *testing.B)   { benchmarkZlibCompress(b, s2048) }
func BenchmarkZlibCompress4096(b *testing.B)   { benchmarkZlibCompress(b, s4096) }
func BenchmarkZlibCompress8192(b *testing.B)   { benchmarkZlibCompress(b, s8192) }
func BenchmarkZlibCompress16384(b *testing.B)   { benchmarkZlibCompress(b, s16384) }
func BenchmarkZlibCompress32768(b *testing.B)   { benchmarkZlibCompress(b, s32768) }

func BenchmarkZlibUnCompress128(b *testing.B) { benchmarkZlibUnCompress(b, s128) }
func BenchmarkZlibUnCompress256(b *testing.B) { benchmarkZlibUnCompress(b, s256) }
func BenchmarkZlibUnCompress512(b *testing.B) { benchmarkZlibUnCompress(b, s512) }
func BenchmarkZlibUnCompress1024(b *testing.B) { benchmarkZlibUnCompress(b, s1024) }
func BenchmarkZlibUnCompress2048(b *testing.B) { benchmarkZlibUnCompress(b, s2048) }
func BenchmarkZlibUnCompress4096(b *testing.B) { benchmarkZlibUnCompress(b, s4096) }
func BenchmarkZlibUnCompress8192(b *testing.B) { benchmarkZlibUnCompress(b, s8192) }
func BenchmarkZlibUnCompress16384(b *testing.B) { benchmarkZlibUnCompress(b, s16384) }
func BenchmarkZlibUnCompress32768(b *testing.B) { benchmarkZlibUnCompress(b, s32768) }

// 测试内存：go test -v -bench=Alloc -benchmem gbkToUtf8String_test.go
// func benchmarkZlibCompress(b *testing.B, uncompressed []byte) {

func benchmarkZlibCompress(b *testing.B, s []byte) {
	// s := "测试压缩解压时的消耗"
 //    bs := []byte(s)
	// r := []byte{}
 //    client := NewClient(string(s))

 //    b.SetBytes(int64(len(s)))
	// b.ReportAllocs()
	// b.ResetTimer()

 //    for i := 0; i < b.N; i++ {
 //        r, _ = client.compress(bs)
 //    }

 //    fmt.Println(len(s), len(r))

    /*--------------------------分割线，只测试压缩解压的效率-----------------------------*/
    var buf bytes.Buffer
	w, _ := zlib.NewWriterLevel(&buf, 1)
	// w := zlib.NewWriter(&buf)

    b.SetBytes(int64(len(s)))
	b.ReportAllocs()
	b.ResetTimer()
	
	_, _ = w.Write(s)
	w.Close()
	// fmt.Println(len(s), len(buf.Bytes()))

	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Reset(&buf)
        _, _ = w.Write(s)
        w.Close()
    }
    fmt.Println(len(s), len(buf.Bytes()))
}

func benchmarkZlibUnCompress(b *testing.B, s []byte) {
	// s = "测试压缩解压时的消耗"

    client := NewClient(string(s))
    bs, _ := client.compress([]byte(s))
 //    r := []byte{}

 //    b.SetBytes(int64(len(s)))
	// b.ReportAllocs()
	// b.ResetTimer()

 //    for i := 0; i < b.N; i++ {
 //        // client := NewClient(s)
 //        r, _ = client.uncompress(bs)
 //    }
 //    fmt.Println(len(bs), len(r))

	/*--------------------------分割线，只测试压缩解压的效率-----------------------------*/
    r := bytes.NewReader(bs)
    // zr, _ := zlib.NewReader(r)
    
    b.SetBytes(int64(len(s)))
	b.ReportAllocs()
	b.ResetTimer()

    var buf bytes.Buffer
    for i := 0; i < b.N; i++ {
		r.Reset(bs)
    	zr, _ := zlib.NewReader(r)
		buf.Reset()
    	_, _ = io.Copy(&buf, zr)
        // r, _ = client.uncompress(bs)
    }
    fmt.Println(len(bs), len(buf.Bytes()))
}

func mustLoadFile(f string) []byte {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return b
}
