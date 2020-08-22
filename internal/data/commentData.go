package data

import (
	"context"

	"github.com/Zucke/SocialNetwork/pkg/comments"
	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"github.com/Zucke/SocialNetwork/pkg/likes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//NewPublicationComment add a comment to a publication
func (up *UserPublications) NewPublicationComment(ctx context.Context, comment *comments.Comment) error {
	publication, err := up.FindPublicationByID(ctx, comment.PublicationID)
	if err != nil {
		return err
	}
	comment.ComendID = primitive.NewObjectID()
	publication.Comments = append(publication.Comments, *comment)
	return up.UpdatePublication(ctx, &publication)

}

//ChangeAPublicationComment a comment from a user made to a publication
func (up *UserPublications) ChangeAPublicationComment(ctx context.Context, newComment *comments.Comment) error {
	publication, err := up.FindPublicationByID(ctx, newComment.PublicationID)
	if err != nil {
		return err
	}
	for i, comment := range publication.Comments {
		if comment.ComendID.Hex() == newComment.ComendID.Hex() {
			if comment.UserID.Hex() != ctx.Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
				return errorstatus.ErrorAccesDenied
			}
			newComment.ComendID = comment.ComendID
			newComment.UserID = comment.UserID
			publication.Comments[i] = *newComment
			return up.UpdatePublication(ctx, &publication)

		}
	}
	return errorstatus.ErrorNotFount

}

//DeletePublicationComment delete a comment from a user made to a publication
func (up *UserPublications) DeletePublicationComment(ctx context.Context, commentToDelete *comments.Comment) error {
	publication, err := up.FindPublicationByID(ctx, commentToDelete.PublicationID)
	if err != nil {
		return err
	}
	for i, comment := range publication.Comments {
		if comment.ComendID == commentToDelete.ComendID {
			if comment.UserID.Hex() != ctx.Value(primitive.ObjectID{}).(primitive.ObjectID).Hex() {
				return errorstatus.ErrorAccesDenied
			}
			publication.Comments = append(publication.Comments[:i], publication.Comments[i+1:]...)
			return up.UpdatePublication(ctx, &publication)

		}
	}
	return errorstatus.ErrorNotFount

}

//CommentLiked like a publication
func (up *UserPublications) CommentLiked(ctx context.Context, commentLiked comments.Comment) (*comments.Comment, error) {
	publication, err := up.FindPublicationByID(ctx, commentLiked.PublicationID)
	if err != nil {
		return &commentLiked, err
	}

	for i, comment := range publication.Comments {
		if comment.ComendID == commentLiked.ComendID {
			l := likes.Like{}
			comment.Likes = *l.AppendLike(ctx, comment.Likes)
			publication.Comments[i] = comment
			return &comment, up.UpdatePublication(ctx, &publication)
		}
	}
	return &commentLiked, errorstatus.ErrorNotFount

}
