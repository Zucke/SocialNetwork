package likes

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Like info o who made a like
type Like struct {
	// LikeToID primitive.ObjectID `json:"like_to_id" bson:"like_to_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
}

//AppendLike add a like
func (l *Like) AppendLike(ctx context.Context, likes []Like) *[]Like {
	l.UserID = ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)

	finded := -1
	for i, like := range likes {
		if like.UserID == l.UserID {
			finded = i
		}

	}

	if finded == -1 {
		likes = append(likes, *l)
	} else {
		likes = append(likes[:finded], likes[finded+1:]...)
	}
	return &likes

}
