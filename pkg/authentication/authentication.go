package authentication

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"github.com/Zucke/SocialNetwork/pkg/socialuser"
	jwt "github.com/dgrijalva/jwt-go"
)

//Claim contiaint the claims that use the token
type Claim struct {
	socialuser.SocialUser
	jwt.StandardClaims
}

var (
	privateKey *rsa.PrivateKey
	//PublicKey the public key
	PublicKey *rsa.PublicKey
)

func init() {
	privateBytes, err := ioutil.ReadFile("./cert/private.rsa")
	if err != nil {
		log.Fatal("error reading private key")
	}

	publicBytes, err := ioutil.ReadFile("./cert/public.rsa.pub")

	if err != nil {
		log.Fatal("error reading public key")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("error parsing private key")
	}

	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("error parsing private key")
	}
}

//GenerateJWT generate a JWT token to a user
func GenerateJWT(socialUser socialuser.SocialUser) (string, error) {
	claims := Claim{
		SocialUser: socialUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "log a user",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return result, nil

}
