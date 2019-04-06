//将base64解码（png格式），然后将png格式转jpg 
//2019.4.6 陶
package main
import (
	"os"
	"fmt"
	"encoding/base64"
	"image"
	
	"image/draw"
	"image/jpeg"
	"image/png"
)
func base64Decode(src []byte) ([]byte, error) {
    return base64.StdEncoding.DecodeString(string(src))
}
func main(){
	
	file, err := os.Open("d:\\base64.txt") 
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
	 
	pngfile, err := os.Create("d:\\test.png")      
	defer pngfile.Close()  
    if err != nil {
        fmt.Println(pngfile, err)
        return
	}
	pngfile.Write(pngBlob)

	//png转jpg
	pngfile1, err := os.Open("d:\\test.png")    
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
	fmt.Printf("ok\n")

	f3, err := os.Create("d:\\test.jpg")
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f3, m, &jpeg.Options{90})
     if err != nil {
         panic(err)
     }
}
