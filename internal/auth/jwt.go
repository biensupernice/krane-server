package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// OneYear : unix time for 1 year
	OneYear = time.Now().Add(time.Minute * 525600).Unix()
)

// AuthClaims : custom claims for user authentication
type AuthClaims struct {
	Phrase string `json:"phrase"`
	jwt.StandardClaims
}

// CustomClaims : custom claims for request
type CustomClaims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

// CreateToken : new jwt token
func CreateToken(SigningKey []byte, data interface{}) (string, error) {
	if SigningKey == nil {
		return "", fmt.Errorf("Cannot create token - signing key not provided")
	}

	customClaims := &CustomClaims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: OneYear,
			Issuer:    "krane-server",
		},
	}

	// Declare the unsigned token using RSA HS256 Algorithm for ecryption
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign the token with secret
	signedTkn, err := tkn.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return signedTkn, nil
}

// ParseJWTToken : parse jwt using signing key
func ParseJWTToken(signKey string, tknStr string) *jwt.Token {
	tkn, err := jwt.ParseWithClaims(tknStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})

	if err != nil {
		return nil
	}

	if !tkn.Valid {
		return nil
	}

	return tkn
}

// ParseAuthTokenWithAuthKey : get the claims of a jwt auth token
func ParseAuthTokenWithAuthKey(pubKey string, tknStr string) (claims jwt.Claims, err error) {
	// Convert ssh format pub key to rsa pub key
	rsaPubKey, err := DecodePublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	// Validate token signed with private key against rsa public key
	tkn, err := jwt.ParseWithClaims(
		tknStr,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return rsaPubKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Verify token is still valid and not expired
	if !tkn.Valid {
		return nil, errors.New("Invalid token")
	}

	return tkn.Claims, nil
}

// DecodePublicKey : decode ssh-rsa string into rsa public key
func DecodePublicKey(str string) (*rsa.PublicKey, error) {
	// comes in as a three part string
	// split into component parts
	tokens := strings.Split(str, " ")

	if len(tokens) < 2 {
		return nil, fmt.Errorf("Invalid key format; must contain at least two fields (keytype data [comment])")
	}

	keyType := tokens[0]
	data, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return nil, err
	}

	format, e, n, err := getRsaValues(data)

	if format != keyType {
		return nil, fmt.Errorf("Key type said %s, but encoded format said %s.  These should match", keyType, format)
	}

	pubKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return pubKey, nil
}

func readLength(data []byte) ([]byte, uint32, error) {
	lBuf := data[0:4]

	buf := bytes.NewBuffer(lBuf)

	var length uint32

	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return nil, 0, err
	}

	return data[4:], length, nil
}

func readBigInt(data []byte, length uint32) ([]byte, *big.Int, error) {
	var bigint = new(big.Int)
	bigint.SetBytes(data[0:length])
	return data[length:], bigint, nil
}

func getRsaValues(data []byte) (format string, e *big.Int, n *big.Int, err error) {
	data, length, err := readLength(data)
	if err != nil {
		return
	}

	format = string(data[0:length])
	data = data[length:]

	data, length, err = readLength(data)
	if err != nil {
		return
	}

	data, e, err = readBigInt(data, length)
	if err != nil {
		return
	}

	data, length, err = readLength(data)
	if err != nil {
		return
	}

	data, n, err = readBigInt(data, length)
	if err != nil {
		return
	}

	return
}
