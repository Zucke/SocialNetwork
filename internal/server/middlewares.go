package server

import (
	"context"
	"net/http"

	"github.com/Zucke/SocialNetwork/pkg/authentication"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ValidateMiddleware used to validate tokes
func ValidateMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token *jwt.Token
		token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &authentication.Claim{}, func(token *jwt.Token) (interface{}, error) {
			return authentication.PublicKey, nil
		})

		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		if !token.Valid {
			response.HTTPError(w, r, http.StatusUnauthorized, "Invalid Token")
			return
		}
		id := token.Claims.(*authentication.Claim).ID
		ctx := context.WithValue(r.Context(), primitive.ObjectID{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
	return fn

}
