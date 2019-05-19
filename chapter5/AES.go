/*
描述 :  golang  AES/ECB/PKCS5  加密解密
date : 2016-04-08
https://www.cnblogs.com/lavin/p/5373188.html
*/
package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "fmt"
	"strings"
	"os"
)
func main() {
    /*
    *src 要加密或解密的字符串
    *key 用来加密的密钥 密钥长度可以是128bit、192bit、256bit中的任意一个
    *16位key对应128bit

    aes ae4ead79e2bb40d69f163717b61e7984 85:6CFB29781341DDB2B1DD1E88F8C5B39FEC6CB8FCFADA86B01F506C6F9B45EE68:1
    aes ae4ead79e2bb40d69f163717b61e7984 1
    aes d ae4ead79e2bb40d69f163717b61e7984 YEtumQjPD6Vqq4gn2UT6qg==
    aes d ae4ead79e2bb40d69f163717b61e7984 avsmYXH2kwROgfn713Igddp0ubo3le7JCFzg9lGA5sGiQu3fUwlQn/Ta8Uw9CDarV/COPN+KAYjrt2g5IzNmHCSk2eglQpK7chWI48AIpYI=
     */
     
     switch len(os.Args) {
        case 4 :
            key := os.Args[2]
            crypted := os.Args[3]
           // fmt.Println([]byte(crypted), string([]byte(crypted)));
            //AesDecrypt([]byte(crypted), []byte(key))
            iReturn := saveAesDecrypt(crypted, []byte(key))
            os.Exit(iReturn)
        case 3:
            key := os.Args[1]
            src := os.Args[2]
            iReturn := saveAesEncrypt(src, key)
            os.Exit(iReturn)
        default:
            panic("使用方式：aes.exe key text 或解密 aes.exe d key text " )
     }
}

func saveAesEncrypt(src , key string) int {
    //将加密后的结果写到txt  
	crypted := AesEncrypt(src, key)
	userFile := "AES.txt"
    fout, err := os.Create(userFile)        
    if err != nil {
        fmt.Println(userFile, err)
        return -1
    }
    defer fout.Close()
 
    fout.Write([]byte(base64.StdEncoding.EncodeToString(crypted) ))
    fout.Write([]byte(key))
    //fout.Writeln([]byte(base64.StdEncoding.EncodeToString(crypted) ))
    
    fmt.Println("------end---------")

    // src := AesDecrypt( []byte("YEtumQjPD6Vqq4gn2UT6qg=="), []byte(key))
    // fout.Write([]byte('解密----'))
    // fout.Write(src)
    //  iReturn := saveAesDecrypt( []byte("YEtumQjPD6Vqq4gn2UT6qg==")  , []byte(key))
   //  fmt.Println( iReturn)
    return 100
}

func AesEncrypt(src, key string) []byte {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        fmt.Println("key error1", err)
    }
    if src == "" {
        fmt.Println("plain content empty")
    }
    ecb := NewECBEncrypter(block)
    content := []byte(src)
    content = PKCS5Padding(content, block.BlockSize())
    crypted := make([]byte, len(content))
    ecb.CryptBlocks(crypted, content)
    // 普通base64编码加密 区别于urlsafe base64
    fmt.Println("base64 result:", base64.StdEncoding.EncodeToString(crypted))

    fmt.Println("base64UrlSafe result:", Base64UrlSafeEncode(crypted))
    //测试解密
    //AesDecrypt(crypted, []byte(key)) 
    return crypted
}

func saveAesDecrypt(crypted string, key []byte) int {
    //将解密后的结果写到txt  

    decrypted, _ := base64.StdEncoding.DecodeString(crypted)
    
	src := AesDecrypt(decrypted, key)
	userFile := "AES.txt"
    fout, err := os.Create(userFile)        
    if err != nil {
        fmt.Println(userFile, err)
        return -1
    }
    defer fout.Close()
      
    fout.Write(src)
   
    return 100
}

func Base64URLDecode(data string) ([]byte, error) {
    var missing = (4 - len(data)%4) % 4
    data += strings.Repeat("=", missing)
    res, err := base64.URLEncoding.DecodeString(data)
    fmt.Println("  decodebase64urlsafe is :", string(res), err)
    return base64.URLEncoding.DecodeString(data)
}

func Base64UrlSafeEncode(source []byte) string {
    // Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
    bytearr := base64.StdEncoding.EncodeToString(source)
    safeurl := strings.Replace(string(bytearr), "/", "_", -1)
    safeurl = strings.Replace(safeurl, "+", "-", -1)
    safeurl = strings.Replace(safeurl, "=", "", -1)
    return safeurl
}

func AesDecrypt(crypted, key []byte) []byte {
   // fmt.Print(string(crypted))
   // fmt.Print(string(key))
    
    block, err := aes.NewCipher(key)
    if err != nil {
        fmt.Println("err is:", err)
    }
    blockMode := NewECBDecrypter(block)
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = PKCS5UnPadding(origData)
    fmt.Println("source is :", origData, string(origData))
    return origData
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

type ecb struct {
    b         cipher.Block
    blockSize int
}

func newECB(b cipher.Block) *ecb {
    return &ecb{
        b:         b,
        blockSize: b.BlockSize(),
    }
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
    return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
    if len(src)%x.blockSize != 0 {
        panic("crypto/cipher: input not full blocks")
    }
    if len(dst) < len(src) {
        panic("crypto/cipher: output smaller than input")
    }
    for len(src) > 0 {
        x.b.Encrypt(dst, src[:x.blockSize])
        src = src[x.blockSize:]
        dst = dst[x.blockSize:]
    }
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
    return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
    if len(src)%x.blockSize != 0 {
        panic("crypto/cipher: input not full blocks")
    }
    if len(dst) < len(src) {
        panic("crypto/cipher: output smaller than input")
    }
    for len(src) > 0 {
        x.b.Decrypt(dst, src[:x.blockSize])
        src = src[x.blockSize:]
        dst = dst[x.blockSize:]
    }
} 