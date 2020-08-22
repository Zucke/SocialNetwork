package comments

import (
	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"github.com/Zucke/SocialNetwork/pkg/likes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Comment has the information of a comment
type Comment struct {
	ComendID      primitive.ObjectID `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	PublicationID primitive.ObjectID `json:"publication_id" bson:"publication_id"`
	Likes         []likes.Like       `json:"likes,omitempty" bson:"likes"`
	Content       string             `json:"content" bson:"content"`
}

//IsValidFields validate the field of this data
func (c *Comment) IsValidFields() error {
	if c.Content != "" {
		return nil
	}
	return errorstatus.ErrorBadInfo

}
