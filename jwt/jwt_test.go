package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/xeenhl/authService/model"
)

type TestKeyLoader struct{}

func (TestKeyLoader) LoadPrivateKey(path string) *rsa.PrivateKey {

	panic("implement me")
}

func (TestKeyLoader) LoadPublicKey(path string) *rsa.PublicKey {
	panic("implement me")
}

func (TestKeyLoader) InitializeKeysChain() *KeyChain {

	public := `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCAjgwGMYDFdUoXIfp9K4ZHhI3+
jNVGPggqjeufZdZwK8ayKMShC9BDA9tIRSqoG5l2gvouvwni27bUHAMpDlV1Rjf5
opadP6i0NSRRB7BwJAk6Iy6Z0vAZuvEuMRa/PK7nLxbmYqQbcidtHi1fPysCR4j1
e19yhPo3VXS+uR8PLQIDAQAB
-----END PUBLIC KEY-----
`

	private := `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCAjgwGMYDFdUoXIfp9K4ZHhI3+jNVGPggqjeufZdZwK8ayKMSh
C9BDA9tIRSqoG5l2gvouvwni27bUHAMpDlV1Rjf5opadP6i0NSRRB7BwJAk6Iy6Z
0vAZuvEuMRa/PK7nLxbmYqQbcidtHi1fPysCR4j1e19yhPo3VXS+uR8PLQIDAQAB
AoGAIsxQnOyReuHA6HoeH/vEIV/UP+9HW/g2pa489azPWxW+d0Np1l4oRbupg+qV
HWQ7KkVSC41S08G9v7TFdjuXDgC98XRewifX1eWSPH8nyo+qywK4QSJK9IlfXUH7
kQfRnzUgXreC6NkPceZttt7UUPJYyLiSpwnC2trCOkYZoIECQQDM9kELUQqTubm9
HpOPGGzgmoCRrmQ8k8v3o5cHsgYVHJmiuIOJbMWfM4IYFB6bMTSfUVuPldIRzzFO
Dj+q9oRVAkEAoJEO2lepUa8rdH/Z90p3aMJxGsaiEbZs/HMwDLGCNiMV3S+/nEuX
OWgDadlVs5KZNNzCXXPbQ2dzRk1uBbt3eQJACdn0OmUEyyDsKojjsscLxKfochgd
vUOlVBvK0JXf8PfU8ptHxz0xKnvBTwL4jaEJ1HaGnhonZK++wO+yY7dBmQJAZGdw
xz3rxgVogff0v3sUQjDccybkb3kIm7AXysgxKVM1N9PE2KI4FRCimczql1jDbtfg
vnlVEcgdwEdo1jLM2QJBAKU+/XybS+Z/tiQLpZ6meuP5jS01Lws2i4bkBwSrJl1d
BlpueQG3OxV6lqK8EpO9k9C4drpT5r0WrvFRE+9YLBc=
-----END RSA PRIVATE KEY-----
`

	keyPb, _ := pem.Decode([]byte(public))
	keyPr, _ := pem.Decode([]byte(private))

	rsaPublicRaw, _ := x509.ParsePKIXPublicKey(keyPb.Bytes)
	rsaPrivate, _ := x509.ParsePKCS1PrivateKey(keyPr.Bytes)

	rsaPublic, _ := rsaPublicRaw.(*rsa.PublicKey)

	return &KeyChain{
		PublicKey:  rsaPublic,
		PrivateKey: rsaPrivate,
	}
}

func TestNewTokenizer(t *testing.T) {
	k := TestKeyLoader{}
	tok := NewTokenizer(k)

	if tok == nil {
		t.Errorf("New Tokenizer mast return JWT tokenizer but got nil")
	}

	if tok.keyLoader != k {
		t.Errorf("Provided keyloader expecte to be in tokenizer but got %v", tok.keyLoader)
	}
}

func TestJwt_GenerateTokenShouldGenerateValidToken(t *testing.T) {
	k := TestKeyLoader{}
	tok := &Jwt{
		keyLoader: k,
	}

	u := &model.User{Id: 101.0}

	token, err := tok.GenerateToken(u)

	if err != nil {
		t.Errorf("GenerateToken mast generate token but return error: %v", err)
	}

	pToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return k.InitializeKeysChain().PrivateKey, nil
	})

	if claims, ok := pToken.Claims.(jwt.MapClaims); ok {

		if id, ok := claims["sub"].(float64); !ok || id != 101.0 {
			t.Errorf("User id mast be 'sub' claim and with expected value: %v", id)
		}

	} else {
		t.Errorf("GenerateToken mast generate token with MapClaims but got: %v", pToken.Claims)
	}
}

func TestJwt_ParceAndVerifyToken(t *testing.T) {

	k := TestKeyLoader{}
	tok := &Jwt{
		keyLoader: k,
	}

	u := &model.User{Id: 101.0}

	token, _ := tok.GenerateToken(u)

	pToken, err := tok.ParceAndVerifyToken(token)

	if err != nil {
		t.Errorf("Parse token mast parse it but return error: %v", err)
	}

	if claims, ok := pToken.Claims.(jwt.MapClaims); ok {

		if id, ok := claims["sub"].(float64); !ok || id != 101.0 {
			t.Errorf("User id mast be 'sub' claim and with expected value: %v", id)
		}

	} else {
		t.Errorf("GenerateToken mast generate token with MapClaims but got: %v", pToken.Claims)
	}

}
