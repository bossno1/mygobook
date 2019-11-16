package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"

	//"encoding/base64"

	"fmt"
)

//DES 和 3DES加密区别
//前者 加密  密钥必须是8byte
//后者加密 解密 再加密  密钥必须是24byte
func main() {
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

	//定义密钥，必须是24byte
	//	key := []byte("469D5D3AC6594B979E0B6262B3B38DBA")  //469D5D3AC6594B979E0B6262B3B38DBA
	//469D5D3AC6594B97
	//9E0B6262B3B38DBA
	//0000000000000000
	key, _ := hex.DecodeString("469D5D3AC6594B979E0B6262B3B38DBA469D5D3AC6594B97") //16 bytes长度
	fmt.Println(key)
	data, _ := hex.DecodeString("54303832343435393600000000000000")
	fmt.Println(data)
	en := ThriDESEnCryptECB(data, key)
	fmt.Println(en)
	s1 := fmt.Sprintf("%x", en) //将[]byte转成16进制
	fmt.Println(s1)
	de := ThriDESDeCryptECB(en, key)
	fmt.Println(fmt.Sprintf("%x", de))

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
