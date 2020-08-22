package publications

import (
	"github.com/Zucke/SocialNetwork/pkg/comments"
	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"github.com/Zucke/SocialNetwork/pkg/likes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Publication has the field of a publication that make a user
type Publication struct {
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content  string             `json:"content" bson:"content,"`
	Likes    []likes.Like       `json:"likes,omitempty" bson:"likes"`
	Comments []comments.Comment `json:"commets,omitempty" bson:"commets"`
}

//IsValidFields validate the field of this data
func (p *Publication) IsValidFields() error {
	if p.Content != "" {
		return nil
	}
	return errorstatus.ErrorBadInfo

}
