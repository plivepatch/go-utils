package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"strings"
)

type AesEncrypt struct {
	Key string
	Iv  string
}

func NewEnc() *AesEncrypt {
	return &AesEncrypt{}
}

func (this *AesEncrypt) getKey() ([]byte, error) {

	keyLen := len(this.Key)
	if keyLen < 16 {
		return []byte(""), errors.New("aes key less than 16 chars.")
	}
	arrKey := []byte(this.Key)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32], nil
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24], nil
	}
	//取前16个字节
	return arrKey[:16], nil
}

//加密字符串
func (this *AesEncrypt) Encrypt(plantText []byte) (string, error) {

	key, errs := this.getKey()
	if errs != nil {
		return "", errs
	}
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return "", err
	}
	plantText = this.PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, []byte(this.Iv)[:aes.BlockSize])

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)

	//mlog.Debug(string(ciphertext))

	return strings.ToUpper(hex.EncodeToString(ciphertext)), nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(strSrc string) ([]byte, error) {

	src, errs := hex.DecodeString(strings.ToLower(strSrc))
	if errs != nil {
		return nil, errs
	}

	key, errs := this.getKey()
	if errs != nil {
		return nil, errs
	}
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, []byte(this.Iv)[:aes.BlockSize])
	plantText := make([]byte, len(src))
	blockModel.CryptBlocks(plantText, src)
	plantText = this.PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

//补位
func (this *AesEncrypt) PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//补位
func (this *AesEncrypt) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
