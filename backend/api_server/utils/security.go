package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"os"
	"strings"
)

const DATA_SIZE = 64
const READ_OFFSET = 8

type Security struct {
	KEY []byte
	IV  []byte
}

type Cipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(cipherIvKey string) (string, error)
	EncryptFile(d string, filePath string) error
	DecryptFile(f string) string
}

func (s Security) Encrypt(plaintext string) (string, error) {
	if strings.TrimSpace(plaintext) == "" {
		return plaintext, nil
	}

	block, err := aes.NewCipher(s.KEY)
	if err != nil {
		fmt.Println("Encrypt Error : ", err.Error())
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(block, s.IV)
	paddedPlainText := padPKCS([]byte(plaintext), encrypter.BlockSize())

	cipherText := make([]byte, len(paddedPlainText))

	encrypter.CryptBlocks(cipherText, paddedPlainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (s Security) Decrypt(cipherText string) (string, error) {
	if strings.TrimSpace(cipherText) == "" {
		return cipherText, nil
	}

	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		fmt.Println("Decrypt Error : ", err.Error())
		return "", err
	}

	block, err := aes.NewCipher(s.KEY)
	if err != nil {
		fmt.Println("Decrypt Error : ", err.Error())
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, s.IV)
	plainText := make([]byte, len(decodedCipherText))

	decrypter.CryptBlocks(plainText, decodedCipherText)
	trimmedPlainText := unpadPKCS(plainText)

	return string(trimmedPlainText), nil
}

// Byte to Byte
// Text to Text

func NewCrypton(cipherKey string, cipherIvKey string) (Cipher, error) {
	if ck := len(cipherKey); ck != 16 {
		return nil, aes.KeySizeError(ck)
	}

	if cik := len(cipherIvKey); cik != 16 {
		return nil, aes.KeySizeError(cik)
	}

	return &Security{[]byte(cipherKey), []byte(cipherIvKey)}, nil
}

func CreateSecurity() Cipher {
	return &Security{[]byte("First___KaierEnc"), []byte("Second___KaierIV")}
}

func padPKCS(p []byte, blockSize int) []byte {
	pad := " "
	padding := blockSize - len(p)%blockSize
	padText := bytes.Repeat([]byte(pad), padding)
	return append(p, padText...)
}

func unpadPKCS(b []byte) []byte {
	padding := len(string(b[len(b)-1]))
	return b[:len(b)-int(padding)]
}

func (s Security) EncryptFile(data string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("File Open Fail - ", filePath)
		return err
	}
	defer file.Close()

	block, err := aes.NewCipher(s.KEY)
	if err != nil {
		return err
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(len(data)))
	err = binary.Write(file, binary.LittleEndian, b)
	if err != nil {
		return err
	}

	_, errW := file.Write(s.IV)
	if errW != nil {
		return errW
	}

	trimData := strings.TrimLeft(data, " ")
	trimData = strings.TrimRight(trimData, " ")
	plaintext := []byte(trimData)

	encrypter := cipher.NewCBCEncrypter(block, s.IV)
	paddedPlainText := padPKCS(plaintext, encrypter.BlockSize())

	cipherText := make([]byte, len(paddedPlainText))

	encrypter.CryptBlocks(cipherText, paddedPlainText)

	if _, err := file.Write(cipherText); err != nil {
		return err
	}

	return nil
}

func (s Security) DecryptFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("File Open Failed - ", filePath, " : ", err)
		return ""
	}
	defer file.Close()

	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("File read failed - ", filePath, " : ", err)
		return ""
	}

	if len(encryptedData) < aes.BlockSize {
		fmt.Println("Data to short.")
		return ""
	}

	ciphertext := encryptedData

	// AES cipher 블록 생성
	block, err := aes.NewCipher(s.KEY)
	if err != nil {
		fmt.Println("Failed to create new cipher : ", err.Error())
		return ""
	}

	// CBC 복호화기 생성
	mode := cipher.NewCBCDecrypter(block, s.IV)
	if len(ciphertext)%aes.BlockSize != 0 {
		fmt.Println("invalid encryption")
		return ""
	}

	// decrypt
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted = unpadPKCS(decrypted)

	return string(decrypted)
}
