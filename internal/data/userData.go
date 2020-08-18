package data

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//SocialUser is a extension from a normal authentication.User
type SocialUser struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Email           string               `json:"email" bson:"email"`
	Name            string               `json:"name,omitempty" bson:"nickname"`
	Surname         string               `json:"surname,omitempty" bson:"surname"`
	Password        string               `json:"password,omitempty" bson:"password"`
	ConfirmPassword string               `json:"confirm_password,omitempty" bson:"confirm_password,omitempty"`
	Friends         []primitive.ObjectID `json:"friends,omitempty" bson:"friends"`
	BlokedUser      []primitive.ObjectID `json:"bloked_user,omitempty" bson:"bloked_user"`
	Biography       string               `json:"biography,omitempty" bson:"biography"`
}

//User params for db connection
type User struct {
	data *Data
	coll *mongo.Collection
}

//GetPublicUserData get the public user data
func (u *SocialUser) GetPublicUserData() SocialUser {
	socialUser := SocialUser{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Surname:   u.Surname,
		Biography: u.Biography,
	}
	return socialUser

}

//IsValidFields validate the fields
func (u *SocialUser) IsValidFields() bool {
	return u.IsValidEmail() && u.IsPassWordEqualTo(u.ConfirmPassword)

}

//IsPassWordEqualTo macth with a password
func (u *SocialUser) IsPassWordEqualTo(password string) bool {
	v := u.Password == password
	return v
}

//CleanPassword clean the password field to evoid send
func (u *SocialUser) CleanPassword() {
	u.Password = ""
	u.ConfirmPassword = ""
}

//IsValidEmail validated email
func (u *SocialUser) IsValidEmail() bool {
	r := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return r.Match([]byte(u.Email))

}

//IsUserInFriendList find if a user in this user friend list return the index if not in retunr -1
func (u *SocialUser) IsUserInFriendList(user primitive.ObjectID) int {
	for i, value := range u.Friends {
		if value.Hex() == user.Hex() {
			return i
		}
	}
	return -1
}

//IsUserInBLockList find if a user is in this user block List and return the index if not in retunr -1
func (u *SocialUser) IsUserInBLockList(user primitive.ObjectID) int {
	for i, value := range u.BlokedUser {
		if value.Hex() == user.Hex() {
			return i
		}
	}
	return -1

}

//GetUserByEmail find user by email
func (u *User) GetUserByEmail(ctx context.Context, email string) (SocialUser, error) {
	result := SocialUser{}
	err := u.coll.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return result, ErrorNotFount
	}
	return result, nil

}

//GetUserByID verify user by ID
func (u *User) GetUserByID(ctx context.Context) (SocialUser, error) {
	result := SocialUser{}
	err := u.coll.FindOne(ctx, bson.M{"_id": ctx.Value(primitive.ObjectID{})}).Decode(&result)
	if err != nil {
		return result, ErrorNotFount
	}
	return result, nil

}

//UpdateUser update a user
func (u *User) UpdateUser(ctx context.Context, newSocialUserData *SocialUser) error {
	_, err := u.coll.UpdateOne(ctx, bson.M{"_id": ctx.Value(primitive.ObjectID{})}, bson.M{"$set": newSocialUserData})
	return err
}

//NewSocialUser add a new SocialUser
func (u *User) NewSocialUser(ctx context.Context, user SocialUser) error {
	_, err := u.coll.InsertOne(ctx, user)
	return err

}

//UserToFriendList add a user to friend list if the user exist delete the user from friend list
func (u *User) UserToFriendList(ctx context.Context, user *SocialUser) {
	ID := ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)
	i := user.IsUserInFriendList(ID)
	if i == -1 {
		user.Friends = append(user.Friends, ID)
	} else {
		user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
	}

}

//UserToBlockList add a user to block list if the user exist delete the user from block list
func (u *User) UserToBlockList(ctx context.Context, user *SocialUser) {
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
