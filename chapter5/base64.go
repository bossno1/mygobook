package main
import (
	"os"
	"fmt"
	"encoding/base64"
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
	jpgBase64    := []byte("")
	buf := make([]byte, 1024)
    for {
        n, _ := file.Read(buf)
        if 0 == n {
            break
		}
		jpgBase64 = append(jpgBase64,  buf[:n]...)
       // jpg.Write(buf[:n])
    }
	// decode
	jpgBlob, err2 := base64Decode(jpgBase64)
	if err2 !=nil {
		fmt.Println(err)
	}
	//jpg文件
	 
	jpg, err := os.Create("d:\\test.jpg")      
	defer jpg.Close()  
    if err != nil {
        fmt.Println(jpg, err)
        return
	}
	jpg.Write(jpgBlob)
}
