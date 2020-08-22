package socialuser

import (
	"regexp"

	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (u *SocialUser) IsValidFields() error {
	if u.IsValidEmail() && u.IsPassWordEqualTo(u.ConfirmPassword) {
		return nil
	}
	return errorstatus.ErrorBadInfo

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
