package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	//"github.com/mergermarket/go-pkcs7"
)

// Cipher key must be 32 chars long because block size is 16 bytes
const CIPHER_KEY = "abcdefghijklmnopqrstuvwxyz012345"

//使用PKCS7进行填充，IOS也是7
//PKCS7Padding ..
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS7UnPadding ...
func PKCS7UnPadding(origData []byte, blockSize int) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Encrypt encrypts plain text string into cipher text string
func Encrypt(unencrypted string) (string, error) {
	key := []byte(CIPHER_KEY)
	plainText := []byte(unencrypted)
	plainText = PKCS7Padding(plainText, aes.BlockSize)

	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	if fileObj, err := os.OpenFile("/Users/jf10/pwd2.dat", os.O_RDWR|os.O_CREATE, 0644); err == nil {
		defer fileObj.Close()
		writeObj := bufio.NewWriterSize(fileObj, 4096)

		//使用Write方法,需要使用Writer对象的Flush方法将buffer中的数据刷到磁盘
		//buf := []byte(content)
		if _, err := writeObj.Write(cipherText); err == nil {
			fmt.Println("Successful appending to the buffer with os.OpenFile and bufio's Writer obj Write method.")
			if err := writeObj.Flush(); err != nil {
				panic(err)
			}
			fmt.Println("Successful flush the buffer data to file ")
		}
	}
	return fmt.Sprintf("%x", cipherText), nil
}

// Decrypt decrypts cipher text string into plain text string
func Decrypt(encrypted string) (string, error) {
	key := []byte(CIPHER_KEY)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText = PKCS7UnPadding(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

func main() {
	//plainText := []byte("1234567890")
	encrypted, err := Encrypt("1234567890")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("%x\n", encrypted)
	fmt.Println(encrypted)
	decrypted, err := Decrypt(encrypted)
	fmt.Println(decrypted)
}
