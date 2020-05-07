package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"os"

	//"encoding/base64"

	"fmt"
)

//DES 和 3DES加密区别
//前者 加密  密钥必须是8byte
//后者加密 解密 再加密  密钥必须是24byte
/*
des3 469D5D3AC6594B979E0B6262B3B38DBA C694E8E2031F7665F8DF78B3C1D4FA88 0 des3.txt 88613818
des3 469D5D3AC6594B979E0B6262B3B38DBA C694E8E2031F7665F8DF78B3C1D4FA88 1 des3.txt 88613818
*/
func main() {
	switch len(os.Args) {
	case 6:
		key := os.Args[1]
		crypted := os.Args[2]
		type1 := os.Args[3] //0-解密  1-将解密结果转为 TEXT
		saveFile := os.Args[4]
		pass := os.Args[5] //88613818  固定值
		//fmt.Println(pass)
		if pass != "88613818" {
			os.Exit(-200)
		}
		de1 := DESDeCrypt3(key, crypted, saveFile, type1)
		fmt.Println(de1)

		os.Exit(100)

	default:
		panic("使用方式：des3.exe key text  savefile type1 xxxx  ")
		os.Exit(-100)
	}

	//定义密钥，必须是24byte
	/*
		//des3加解密
		key := []byte("123456789012345678901234")
		//定义明文
		origData := []byte("hello world")
		//加密
		en := ThriDESEnCrypt(origData, key)
		//解密
		de := ThriDESDeCrypt(en, key)
		fmt.Println(string(de))
	*/
	/*
		des3 469D5D3AC6594B979E0B6262B3B38DBA 54303832343435393600000000000000
			key, _ := hex.DecodeString("469D5D3AC6594B979E0B6262B3B38DBA469D5D3AC6594B97") //16 bytes长度
			//fmt.Println(key)
			data, _ := hex.DecodeString("54303832343435393600000000000000")
			//fmt.Println(data)
			en := ThriDESEnCryptECB(data, key)
			fmt.Println(en)
			s1 := fmt.Sprintf("%x", en) //将[]byte转成16进制
			fmt.Println(s1)
			de := ThriDESDeCryptECB(en, key)
			fmt.Println(fmt.Sprintf("%x", de))

			de1 := DESDeCrypt3("469D5D3AC6594B979E0B6262B3B38DBA", "C694E8E2031F7665F8DF78B3C1D4FA88")
			fmt.Println(de1)
			// ---------------------------
			// 解密结果
			// ---------------------------
			// 54303832343435393600000000000000
			// ---------------------------
			// 确定
			// ---------------------------

			XXXX := "54303832343435393600000000000000"

			fmt.Println(XXXX)

			var dst []byte
			fmt.Sscanf(XXXX, "%X", &dst)
			fmt.Println(string(dst)) //T08244596
	*/
	/*
		sText := "中文"
		textQuoted := strconv.QuoteToASCII(sText)

		textUnquoted := textQuoted[1 : len(textQuoted)-1]
		fmt.Println(textUnquoted)

		sUnicodev := strings.Split(textUnquoted, "\\u")
		var context string
		for _, v := range sUnicodev {
			if len(v) < 1 {
				continue
			}
			temp, err := strconv.ParseInt(v, 16, 32)
			if err != nil {
				panic(err)
			}
			context += fmt.Sprintf("%c", temp)
		}
		fmt.Println(context)
	*/
}

//解密
func ThriDESDeCrypt(crypted, key []byte) []byte {
	//获取block块
	block, _ := des.NewTripleDESCipher(key)
	//创建切片
	context := make([]byte, len(crypted))
	//设置解密方式
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	//解密密文到数组
	blockMode.CryptBlocks(context, crypted)
	//去补码
	context = PKCSUnPadding(context)
	//context = PKCSUnPadding(context)

	//context = ZeroUnPadding(context)
	return context
}

//解密   ECB方式, padding=ZERO
func ThriDESDeCryptECB(crypted, key []byte) []byte {
	//获取block块
	block, _ := des.NewTripleDESCipher(key)
	//创建切片
	context := make([]byte, len(crypted))
	//设置解密方式 ECB
	blockMode := NewECBDecrypter(block) //cipher.NewCBCEncrypter(block, key[:8])
	//blockMode := cipher.NewCBCDecrypter(block, key[:8])
	//解密密文到数组
	blockMode.CryptBlocks(context, crypted)
	//去补码 ZERO不用补码
	//context = PKCSUnPadding(context)
	//context = PKCSUnPadding(context)

	//context = ZeroUnPadding(context)
	return context
}
func DESDeCrypt3(pkey string, pdata string, saveFile string, type1 string) string {
	if len(pkey) == 32 {
		pkey = pkey + Substr(pkey, 0, 16)
	} else {
		return ""
	}

	key, _ := hex.DecodeString(pkey) //16 bytes长度

	data, _ := hex.DecodeString(pdata)

	de := ThriDESDeCryptECB(data, key)
	s1 := fmt.Sprintf("%x", de)

	if type1 == "1" {
		var dst []byte
		fmt.Sscanf(s1, "%X", &dst)
		s1 = string(dst)
		fmt.Println(s1)

	}

	//将结果回写到指定的文件
	fout, err := os.Create(saveFile)
	if err != nil {
		fmt.Println("写文件失败")
		fmt.Println(saveFile, err)
		return "-1"

	}
	defer fout.Close()

	fout.Write([]byte(s1))

	return s1
}

//去补码
func PKCSUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

//加密 CBC方式
func ThriDESEnCrypt(origData, key []byte) []byte {
	//获取block块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	//补码
	origData = PKCSPadding(origData, block.BlockSize())
	//设置加密方式为 3DES  使用3条56位的密钥对数据进行三次加密
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	//创建明文长度的数组
	crypted := make([]byte, len(origData))
	//加密明文
	blockMode.CryptBlocks(crypted, origData)
	return crypted
}

//加密 ECB方式, padding=ZERO
func ThriDESEnCryptECB(origData, key []byte) []byte {
	//获取block块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	//补码 (ZERO不要执行补码)
	//origData = PKCSPadding(origData, block.BlockSize())
	//设置加密方式为 3DES  使用3条56位的密钥对数据进行三次加密
	blockMode := NewECBEncrypter(block) //cipher.NewCBCEncrypter(block, key[:8])
	//创建明文长度的数组
	crypted := make([]byte, len(origData))
	//加密明文
	blockMode.CryptBlocks(crypted, origData)
	return crypted
}

//补码
func PKCSPadding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}

//----------------
//des加密
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//des解密
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	//origData := make([]byte, len(crypted))
	origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	//origData = PKCSUnPadding(origData)

	origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//-------AES.GO也用到
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
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
