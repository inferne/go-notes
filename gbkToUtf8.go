
package encode

import (
	// "bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"testing"
	"fmt"
	// "golang.org/x/text/transform"
	// "io/ioutil"
	// "io"
	// "strings"
	// "fmt"
)

/* 传统方式 */
// func GbkToUtf8String(s string) string {
// 	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
// 	body, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(body)
// }

/* 下面是三种封装 */
// 这个占用内存小，性能居中
func gbk2Utf8(str []byte) []byte {
	utf8data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(str)
	// fmt.Println(err)
	return utf8data
}
// 这个占用内存稍大，性能最好，gc压力最小
func gbk2Utf8String(str string) string {
	utf8data, _ := simplifiedchinese.GBK.NewDecoder().String(str)
	// fmt.Println(err)
	return utf8data
}
// 这个占用内存稍大，性能最不好，alloc最多
func gbk2Utf8String2(str string) string {
	utf8data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(str))
	// fmt.Println(err)
	return string(utf8data)
}

// // 测试内存：go test -v -bench=Alloc -benchmem gbkToUtf8String_test.go
// func Benchmark_Alloc_Old(b *testing.B) {
// 	s := "月色真美，风也温柔，233333333，~！@#"
//     // r := GbkToUtf8String(s)
//  	//    reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
// 	// body, err := ioutil.ReadAll(reader)
// 	// if err != nil {
// 	// 	return 
// 	// }
// 	// r := string(body)
// 	 //    fmt.Println(r)

// 	r := ""

//     for i := 0; i < b.N; i++ {
//         r = GbkToUtf8String(s)
//     }
//     fmt.Println(r)
// }

// // flat 表示当前函数直接分配
// // cum 表示下游函数内部分配

func Benchmark_Alloc_New2(b *testing.B) {
    str := "月色真美，风也温柔，233333333，~！@#"  //go字符串编码为utf-8
    fmt.Println("before convert:", str)   //打印转换前的字符串
    // fmt.Println("isUtf8:", isUtf8([]byte(str)))   //判断是否是utf-8

    gbkData := []byte{}
    // gbkData := ""
    for i := 0; i < b.N; i++ {
    	gbkData, _ = simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))  //使用官方库将utf-8转换为gbk
    	// gbkData, _ = simplifiedchinese.GBK.NewEncoder().String(str)  //使用官方库将utf-8转换为gbk
    	// gbkData = gbk2Utf8([]byte(str))
    }
    // gbkStr, err := simplifiedchinese.GBK.NewEncoder().String(str)
    // if err != nil {
    //     panic(err)
    // }
    // fmt.Println("gbk直接打印会出现乱码:", gbkStr)   //乱码字符串
    fmt.Println("gbk直接打印会出现乱码:", string(gbkData))   //乱码字符串
    // fmt.Println("gbk直接打印会出现乱码:", gbkData)   //乱码字符串

    // fmt.Println("isGBK:", isGBK(gbkData)) //判断是否是gbk
    utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkData)) //将gbk再转换为utf-8
    // fmt.Println("isUtf8:", isUtf8(utf8Data) )  //判断是否是utf-8
    fmt.Println("after convert:", string(utf8Data))   //打印转换前的字符串
}

func Benchmark_Alloc_New3(b *testing.B) {
    str := "月色真美，风也温柔，233333333，~！@#"  //go字符串编码为utf-8
    fmt.Println("before convert:", str)   //打印转换前的字符串

    gbkData := []byte{}
    for i := 0; i < b.N; i++ {
    	gbkData = gbk2Utf8([]byte(str))
    }

    fmt.Println("gbk直接打印会出现乱码:", string(gbkData))   //乱码字符串

    // fmt.Println("isGBK:", isGBK(gbkData)) //判断是否是gbk
    utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(gbkData) //将gbk再转换为utf-8
    // fmt.Println("isUtf8:", isUtf8(utf8Data) )  //判断是否是utf-8
    fmt.Println("after convert:", string(utf8Data))   //打印转换前的字符串
}

// go test -v -bench=Alloc_New4 -benchmem gbkToUtf8String_test.go gbkToUtf8String.go
func Benchmark_Alloc_New4(b *testing.B) {
    str := "月色真美，风也温柔，233333333，~！@#"  //go字符串编码为utf-8
    fmt.Println("before convert:", str)   //打印转换前的字符串

    gbkData := ""
    for i := 0; i < b.N; i++ {
    	gbkData = gbk2Utf8String(str)
    }

    fmt.Println("gbk直接打印会出现乱码:", gbkData)   //乱码字符串

    // fmt.Println("isGBK:", isGBK(gbkData)) //判断是否是gbk
    utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkData)) //将gbk再转换为utf-8
    // fmt.Println("isUtf8:", isUtf8(utf8Data) )  //判断是否是utf-8
    fmt.Println("after convert:", string(utf8Data))   //打印转换前的字符串
}

func Benchmark_Alloc_New5(b *testing.B) {
    str := "月色真美，风也温柔，233333333，~！@#"  //go字符串编码为utf-8
    fmt.Println("before convert:", str)   //打印转换前的字符串

    gbkData := ""
    for i := 0; i < b.N; i++ {
    	gbkData = gbk2Utf8String2(str)
    }

    fmt.Println("gbk直接打印会出现乱码:", gbkData)   //乱码字符串

    // fmt.Println("isGBK:", isGBK(gbkData)) //判断是否是gbk
    utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkData)) //将gbk再转换为utf-8
    // fmt.Println("isUtf8:", isUtf8(utf8Data) )  //判断是否是utf-8
    fmt.Println("after convert:", string(utf8Data))   //打印转换前的字符串
}
