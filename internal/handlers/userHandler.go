package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/authentication"
	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/Zucke/SocialNetwork/pkg/socialuser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//UsersRouter has the data for a user and the db connection
type UsersRouter struct {
	SocialUser socialuser.SocialUser
	UserData   data.User
}

func newUserRouter() *UsersRouter {
	return &UsersRouter{
		UserData: data.NewUser(),
	}
}

//DecodeAndValidateUser the request body to the user, validate the information and return a error if exist
func (u *UsersRouter) DecodeAndValidateUser(w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&u.SocialUser)
	return err
}

//RegisterSocialUser register a new social user
func RegisterSocialUser(w http.ResponseWriter, r *http.Request) {
	userRouter := newUserRouter()
	err := userRouter.DecodeAndValidateUser(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	err = userRouter.SocialUser.IsValidFields()
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	_, err = userRouter.UserData.GetUserByEmail(ctx, userRouter.SocialUser.Email)

	if err != errorstatus.ErrorNotFount {
		response.HTTPError(w, r, http.StatusBadRequest, "Email Used")
		return
	}
	userRouter.UserData.NewSocialUser(ctx, userRouter.SocialUser)
	userRouter.SocialUser.CleanPassword()
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"user": userRouter.SocialUser})

}

//LoginUser login the user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	userRouter := newUserRouter()
	err := userRouter.DecodeAndValidateUser(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	resultUser, err := userRouter.UserData.GetUserByEmail(ctx, userRouter.SocialUser.Email)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return

	}
	if !userRouter.SocialUser.IsPassWordEqualTo(resultUser.Password) {
		response.HTTPError(w, r, http.StatusBadRequest, "Bad Information")
		return

	}

	var token string
	userRouter.SocialUser.ID = resultUser.ID
	userRouter.SocialUser.CleanPassword()
	token, err = authentication.GenerateJWT(userRouter.SocialUser)

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"user":  userRouter.SocialUser,
		"token": token,
	})

}

//GetUserFriends return the friend of the currend user
func GetUserFriends(w http.ResponseWriter, r *http.Request) {
	var err error
	userRouter := newUserRouter()
	userRouter.SocialUser, err = userRouter.UserData.GetUserByID(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"friends": userRouter.SocialUser.Friends,
	})

}

//GetUserBlockedUser return the BlockedUser of the currend user
func GetUserBlockedUser(w http.ResponseWriter, r *http.Request) {
	var err error
	userRouter := newUserRouter()
	userRouter.SocialUser, err = userRouter.UserData.GetUserByID(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"blokedusers": userRouter.SocialUser.BlokedUser,
	})

}

//GetUserByID get a user bi the header id
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "user_id"))
	userRouter := newUserRouter()
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.WithValue(context.Background(), primitive.ObjectID{}, userID)

	userRouter.SocialUser, err = userRouter.UserData.GetUserByID(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if !(userRouter.SocialUser.IsUserInBLockList(ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)) == -1) {
		return
	}

	render.JSON(w, r, render.M{
		"public_user_data": userRouter.SocialUser.GetPublicUserData(),
	})

}

//UserToFriendList add a user to friend list
func UserToFriendList(w http.ResponseWriter, r *http.Request) {
	userRouter := newUserRouter()
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "user_id"))
	if userID.Hex() == r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
		return
	}

	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.WithValue(context.Background(), primitive.ObjectID{}, userID)

	userRouter.SocialUser, err = userRouter.UserData.GetUserByID(r.Context())
	if userRouter.SocialUser.IsUserInBLockList(userID) != -1 {
		userRouter.UserData.UserToBlockList(ctx, &userRouter.SocialUser)
	}

	userRouter.UserData.UserToFriendList(ctx, &userRouter.SocialUser)
	err = userRouter.UserData.UpdateUser(r.Context(), &userRouter.SocialUser)
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
	userRouter := newUserRouter()
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "user_id"))
	if userID.Hex() == r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
		return
	}
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := context.WithValue(context.Background(), primitive.ObjectID{}, userID)

	userRouter.SocialUser, err = userRouter.UserData.GetUserByID(r.Context())
	if userRouter.SocialUser.IsUserInFriendList(userID) != -1 {
		userRouter.UserData.UserToFriendList(ctx, &userRouter.SocialUser)
	}

	userRouter.UserData.UserToBlockList(ctx, &userRouter.SocialUser)
	err = userRouter.UserData.UpdateUser(r.Context(), &userRouter.SocialUser)
	if err != nil {
		response.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"blockedUser": userID,
	})

}
