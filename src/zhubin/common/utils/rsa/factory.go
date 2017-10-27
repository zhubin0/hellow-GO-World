package rsa

func NewRSA(key Key) (Cipher, error) {
	padding := NewPKCS1Padding(key.Modulus())
	cipherMode := NewPKCS1v15Cipher()
	signMode := NewPKCS1v15Sign()
	return NewCipher(key, padding, cipherMode, signMode), nil
}

func NewRSAWith(key Key, padding Padding, cipherMode CipherMode, signMode SignMode) (Cipher, error) {
	return NewCipher(key, padding, cipherMode, signMode), nil
}
