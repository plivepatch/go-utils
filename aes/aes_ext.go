package aes

func StrAesEncrypt(iv string, key string, source []byte) (string, error) {
	aesEnc := NewEnc()
	aesEnc.Iv = iv
	aesEnc.Key = key
	return aesEnc.Encrypt(source)
}

func StrAesDecrypt(iv string, key string, source string) ([]byte, error) {
	aesEnc := NewEnc()
	aesEnc.Iv = iv
	aesEnc.Key = key
	return aesEnc.Decrypt(source)
}
