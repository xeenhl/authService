package jwt

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/xeenhl/authService/config"
)

type PrivateKeyLoader interface {
	LoadPrivateKey(path string) *rsa.PrivateKey
}

type PublicKeyLoader interface {
	LoadPublicKey(path string) *rsa.PublicKey
}

type KeyLoader interface {
	PrivateKeyLoader
	PublicKeyLoader
	InitializeKeysChain() *KeyChain
}

type KeyChain struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type FileKeyLoader struct {
	privateKey string
	publicKey  string
}

var keys *KeyChain = nil

func NewFileKeyLoader(c config.Configuration) *FileKeyLoader {
	return &FileKeyLoader{
		privateKey: c.Auth.PrivateKey,
		publicKey:  c.Auth.PublicKey,
	}
}

func (f *FileKeyLoader) InitializeKeysChain() *KeyChain {
	if keys == nil {
		keys = &KeyChain{
			PublicKey:  f.LoadPublicKey(f.publicKey),
			PrivateKey: f.LoadPrivateKey(f.privateKey),
		}
	}
	return keys
}

func (f *FileKeyLoader) LoadPrivateKey(path string) *rsa.PrivateKey {

	data := readKey(path)

	pubKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return pubKeyImported
}

func (f *FileKeyLoader) LoadPublicKey(path string) *rsa.PublicKey {

	data := readKey(path)

	pubKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := pubKeyImported.(*rsa.PublicKey)

	if !ok {
		panic("cant validate public key")
	}

	return rsaPub
}

func readKey(path string) *pem.Block {

	p, err := filepath.Abs(path)

	if err != nil {
		panic(err)
	}

	prKey, err := os.Open(p)
	defer prKey.Close()

	if err != nil {
		panic(err)
	}

	keyStat, _ := prKey.Stat()

	var s int64 = keyStat.Size()
	pembytes := make([]byte, s)

	buffer := bufio.NewReader(prKey)
	buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	return data
}
