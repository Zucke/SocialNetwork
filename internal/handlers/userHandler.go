package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/authentication"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//RegisterSocialUser register a new social user
func RegisterSocialUser(w http.ResponseWriter, r *http.Request) {
	var newSocialUser data.SocialUser
	err := json.NewDecoder(r.Body).Decode(&newSocialUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, "Failed to parse user")
		return
	}
	if !newSocialUser.IsValidFields() {
		response.HTTPError(w, r, http.StatusBadRequest, "Bad Information")
		return
	}

	newUser := data.NewUser()
	ctx := context.Background()
	_, err = newUser.GetUserByEmail(ctx, newSocialUser.Email)

	if err != data.ErrorNotFount {
		response.HTTPError(w, r, http.StatusBadRequest, "Email Used")
		return
	}
	newUser.NewSocialUser(ctx, newSocialUser)
	newSocialUser.CleanPassword()
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"user": newSocialUser.Name})

}

//LoginUser login the user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var newSocialUser data.SocialUser
	err := json.NewDecoder(r.Body).Decode(&newSocialUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, "Failed to parse user")
		return
	}

	newUser := data.NewUser()
	ctx := context.Background()
	resultUser, err := newUser.GetUserByEmail(ctx, newSocialUser.Email)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return

	}
	if !newSocialUser.IsPassWordEqualTo(resultUser.Password) {
		response.HTTPError(w, r, http.StatusBadRequest, "Bad Information")
		return

	}

	var token string
	newSocialUser.ID = resultUser.ID
	newSocialUser.CleanPassword()
	token, err = authentication.GenerateJWT(newSocialUser)

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"user":  newSocialUser,
		"token": token,
	})

}

//GetUserFriends return the friend of the currend user
func GetUserFriends(w http.ResponseWriter, r *http.Request) {
	newUser := data.NewUser()
	socialUser, err := newUser.GetUserByID(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"friends": socialUser.Friends,
	})

}

//GetUserBlockedUser return the BlockedUser of the currend user
func GetUserBlockedUser(w http.ResponseWriter, r *http.Request) {
	newUser := data.NewUser()
	socialUser, err := newUser.GetUserByID(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"blokedusers": socialUser.BlokedUser,
	})

}

//GetUserByID get a user bi the header id
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "user_id"))

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.WithValue(context.Background(), primitive.ObjectID{}, userID)

	newUser := data.NewUser()
	socialUser, err := newUser.GetUserByID(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if !(socialUser.IsUserInBLockList(ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)) == -1) {
		return
	}

	render.JSON(w, r, render.M{
		"public_user_data": socialUser.GetPublicUserData(),
	})

}

//UserToFriendList add a user to friend list
func UserToFriendList(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "user_id"))
	if userID.Hex() == r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
		return
	}

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.WithValue(context.Background(), primitive.ObjectID{}, userID)

	newUser := data.NewUser()
	socialUser, err := newUser.GetUserByID(r.Context())
	if socialUser.IsUserInBLockList(userID) != -1 {
		newUser.UserToBlockList(ctx, &socialUser)
	}

	newUser.UserToFriendList(ctx, &socialUser)
	err = newUser.UpdateUser(r.Context(), &socialUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"friend": userID,
	})

}

//UserToBlockList add a user to friend list
func UserToBlockList(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "user_id"))
	if userID.Hex() == r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
		return
	}

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.WithValue(context.Background(), primitive.ObjectID{}, userID)

	newUser := data.NewUser()
	socialUser, err := newUser.GetUserByID(r.Context())
	if socialUser.IsUserInFriendList(userID) != -1 {
		newUser.UserToFriendList(ctx, &socialUser)
	}

	newUser.UserToBlockList(ctx, &socialUser)
	err = newUser.UpdateUser(r.Context(), &socialUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"blockedUser": userID,
	})

}
