//将base64解码（png格式），然后将png格式转jpg 
//2019.4.6 陶
package main
import (
	"os"
	"fmt"
	"encoding/base64"
	"image"
	"os/exec"
	"image/draw"
	"image/jpeg"
	"image/png"
	//"strconv"
	"strings"
)
func base64Decode(src []byte) ([]byte, error) {
    return base64.StdEncoding.DecodeString(string(src))
}

func getPath(fileAndPath string) string {
    s, err := exec.LookPath(fileAndPath)
    if err != nil {
        panic(err)
    }
    i := strings.LastIndex(s, "\\")
    path := string(s[0 : i+1])
    return path
}
func PathExists(path string) (bool) {
	/*
	golang判断文件或文件夹是否存在的方法为使用os.Stat()函数返回的错误值进行判断:

	如果返回的错误为nil,说明文件或文件夹存在
	如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
	如果返回的错误为其它类型,则不确定是否在存在
	*/

	_, err := os.Stat(path)
	if err == nil {
		return true 
	}
	if os.IsNotExist(err) {
		return false 
	}
	return false 
}

func main(){
	 
	//path := getCurrentPath()
	//fmt.Println(path)
	if len(os.Args) != 3 {
		fmt.Println("参数不正确:" , len(os.Args))
		panic("使用方式：base64PngToJpg.exe  xxxxxx.txt  xxxxxx.jpg"  + string(len(os.Args)))
		os.Exit(-3)
	}
	parmBase64txt := os.Args[1]
	if ! PathExists(parmBase64txt){
		fmt.Println("文件不存在：" + parmBase64txt)
		os.Exit(-2)
	}
	tempJpgfile := os.Args[2]
	 
	tempPngfile := getPath(parmBase64txt) + "\\test.png"
 
	fmt.Println(tempPngfile)
	// for idx, args := range os.Args {
    //     fmt.Println("参数" + strconv.Itoa(idx) + ":", args)
	// }
	file, err := os.Open(parmBase64txt) //"d:\\base64.txt") 
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//read base64文件
	pngBase64    := []byte("")
	buf := make([]byte, 1024)
    for {
        n, _ := file.Read(buf)
        if 0 == n {
            break
		}
		pngBase64 = append(pngBase64,  buf[:n]...)
       // jpg.Write(buf[:n])
    }
	// decode
	pngBlob, err2 := base64Decode(pngBase64)
	if err2 !=nil {
		fmt.Println(err)
	}
	//png文件
	 
	pngfile, err := os.Create(tempPngfile)      
	defer pngfile.Close()  
    if err != nil {
        fmt.Println(pngfile, err)
        return
	}
	pngfile.Write(pngBlob)

	//png转jpg
	pngfile1, err := os.Open(tempPngfile)    
	if err != nil {
        fmt.Println(pngfile1, err)
        return
	}  
	m1, err := png.Decode(pngfile1)
	if err != nil {
		panic(err)
	}
	bounds := m1.Bounds()
	m := image.NewRGBA(bounds)
	draw.Draw(m, bounds, m1, bounds.Min, draw.Src)
	

	f3, err := os.Create(tempJpgfile)
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f3, m, &jpeg.Options{90})
     if err != nil {
         panic(err)
	 }
	 fmt.Printf("ok\n")
	 os.Exit(100)
}
