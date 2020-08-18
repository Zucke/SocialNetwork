package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Comment has the information of a comment
type Comment struct {
	ComendID      primitive.ObjectID `json:"comment_id,omitempty" bson:"comment_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	PublicationID primitive.ObjectID `json:"publication_id" bson:"publication_id"`
	Likes         []Like             `json:"likes,omitempty" bson:"likes"`
	Content       string             `json:"content" bson:"content"`
}

//IsValidFields validate the field of this data
func (c *Comment) IsValidFields() bool {
	return c.Content != ""

}

//NewPublicationComment add a comment to a publication
func (up *UserPublications) NewPublicationComment(ctx context.Context, comment *Comment) error {
	publication, err := up.FindPublicationByID(ctx, comment.PublicationID)
	if err != nil {
		return err
	}
	comment.ComendID = primitive.NewObjectID()
	publication.Comments = append(publication.Comments, *comment)
	return up.UpdatePublication(ctx, &publication)

}

//ChangeAPublicationComment a comment from a user made to a publication
func (up *UserPublications) ChangeAPublicationComment(ctx context.Context, newComment *Comment) error {
	publication, err := up.FindPublicationByID(ctx, newComment.PublicationID)
	if err != nil {
		return err
	}
	for i, comment := range publication.Comments {
		if comment.ComendID.Hex() == newComment.ComendID.Hex() {
			if comment.UserID.Hex() != ctx.Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
				return ErrorAccesDenied
			}
			newComment.ComendID = comment.ComendID
			newComment.UserID = comment.UserID
			publication.Comments[i] = *newComment
			return up.UpdatePublication(ctx, &publication)

		}
	}
	return ErrorNotFount

}

//DeletePublicationComment delete a comment from a user made to a publication
func (up *UserPublications) DeletePublicationComment(ctx context.Context, commentToDelete *Comment) error {
	publication, err := up.FindPublicationByID(ctx, commentToDelete.PublicationID)
	if err != nil {
		return err
	}
	for i, comment := range publication.Comments {
		if comment.ComendID == commentToDelete.ComendID {
			if comment.UserID.Hex() != ctx.Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
				return ErrorAccesDenied
			}
			publication.Comments = append(publication.Comments[:i], publication.Comments[i+1:]...)
			return up.UpdatePublication(ctx, &publication)

		}
	}
	return ErrorNotFount

}

//CommentLiked like a publication
func (up *UserPublications) CommentLiked(ctx context.Context, commentLiked Comment) (*Comment, error) {
	publication, err := up.FindPublicationByID(ctx, commentLiked.PublicationID)
	if err != nil {
		return &commentLiked, err
	}

	for i, comment := range publication.Comments {
		if comment.ComendID == commentLiked.ComendID {
			l := Like{}
			comment.Likes = *l.AppendLike(ctx, comment.Likes)
			publication.Comments[i] = comment
			return &comment, up.UpdatePublication(ctx, &publication)
		}
	}
	return &commentLiked, ErrorNotFount

}
