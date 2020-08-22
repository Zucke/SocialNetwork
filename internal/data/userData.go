package data

import (
	"context"

	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"github.com/Zucke/SocialNetwork/pkg/socialuser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//User params for db connection
type User struct {
	data *Data
	coll *mongo.Collection
}

//GetUserByEmail find user by email
func (u *User) GetUserByEmail(ctx context.Context, email string) (socialuser.SocialUser, error) {
	result := socialuser.SocialUser{}
	err := u.coll.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return result, errorstatus.ErrorNotFount
	}
	return result, nil

}

//GetUserByID verify user by ID
func (u *User) GetUserByID(ctx context.Context) (socialuser.SocialUser, error) {
	result := socialuser.SocialUser{}
	err := u.coll.FindOne(ctx, bson.M{"_id": ctx.Value(primitive.ObjectID{})}).Decode(&result)
	if err != nil {
		return result, errorstatus.ErrorNotFount
	}
	return result, nil

}

//UpdateUser update a user
func (u *User) UpdateUser(ctx context.Context, newSocialUserData *socialuser.SocialUser) error {
	_, err := u.coll.UpdateOne(ctx, bson.M{"_id": ctx.Value(primitive.ObjectID{})}, bson.M{"$set": newSocialUserData})
	return err
}

//NewSocialUser add a new SocialUser
func (u *User) NewSocialUser(ctx context.Context, user socialuser.SocialUser) error {
	_, err := u.coll.InsertOne(ctx, user)
	return err

}

//UserToFriendList add a user to friend list if the user exist delete the user from friend list
func (u *User) UserToFriendList(ctx context.Context, user *socialuser.SocialUser) {
	ID := ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)
	i := user.IsUserInFriendList(ID)
	if i == -1 {
		user.Friends = append(user.Friends, ID)
	} else {
		user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
	}

}

//UserToBlockList add a user to block list if the user exist delete the user from block list
func (u *User) UserToBlockList(ctx context.Context, user *socialuser.SocialUser) {
	ID := ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)
	i := user.IsUserInBLockList(ID)
	if i == -1 {
		user.BlokedUser = append(user.BlokedUser, ID)
	} else {
		user.BlokedUser = append(user.BlokedUser[:i], user.BlokedUser[i+1:]...)
	}

}

//NewUser return db info
func NewUser() User {
	return User{
		data: New(),
		coll: data.DBCollection(UserColletion),
	}
}
