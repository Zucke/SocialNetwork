package authentication

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the user data
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Name     string             `json:"name" bson:"nickname"`
	Surname  string             `json:"surname" bson:"surname"`
	Password string             `json:"password" bson:"password"`
}

//Claim contiaint the claims that use the token
type Claim struct {
	data.SocialUser
	jwt.StandardClaims
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

//UID the user id vinculate to some resource
type UID struct {
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
}

//UserInfo vinculated to the user info
type UserInfo interface {
	SetUserID(ID primitive.ObjectID)
}

//SetUserID set the UserID
func (ui *UID) SetUserID(ID primitive.ObjectID) {
	ui.UserID = ID
}

//ComparePassword macth with a password
func (u *User) ComparePassword(password string) bool {
	return u.Password == password
}

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

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("error parsing private key")
	}
}

//GenerateJWT generate a JWT token to a user
func GenerateJWT(socialUser data.SocialUser) (string, error) {
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

//ValidateMiddleware used to validate tokes
func ValidateMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token *jwt.Token
		token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &Claim{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		if !token.Valid {
			response.HTTPError(w, r, http.StatusUnauthorized, "Invalid Token")
			return
		}
		id := token.Claims.(*Claim).ID
		ctx := context.WithValue(r.Context(), primitive.ObjectID{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
	return fn

}

//BasicValidations commons userinfo validation
func BasicValidations(Fields data.FieldValidation, w http.ResponseWriter, r *http.Request) bool {
	err := json.NewDecoder(r.Body).Decode(Fields)
	if err != nil {
		return false

	}
	return Fields.IsValidFields()
}
