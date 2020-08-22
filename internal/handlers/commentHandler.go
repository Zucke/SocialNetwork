package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/comments"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CommentsRouter has the data for a comment and the db connection
type CommentsRouter struct {
	Comment         comments.Comment
	PublicationData data.UserPublications
}

func newCommentsRouter() *CommentsRouter {
	return &CommentsRouter{
		PublicationData: data.NewUserPublications(),
	}
}

//DecodeAndValidateComment the request body to the comment, validate the information and return a error if exist
func (cr *CommentsRouter) DecodeAndValidateComment(w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&cr.Comment)
	if err != nil {
		return err
	}
	cr.Comment.PublicationID, err = primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		return err
	}
	return cr.Comment.IsValidFields()

}

//NewPublicationComment add a comment of the current logged user
func NewPublicationComment(w http.ResponseWriter, r *http.Request) {
	commentRouter := newCommentsRouter()
	err := commentRouter.DecodeAndValidateComment(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	commentRouter.Comment.UserID = r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID)
	commentRouter.PublicationData.NewPublicationComment(r.Context(), &commentRouter.Comment)

	render.JSON(w, r, render.M{
		"comment": commentRouter.Comment,
	})

}

//ChangePublicationComment add a comment of the current logged user
func ChangePublicationComment(w http.ResponseWriter, r *http.Request) {
	commentRouter := newCommentsRouter()
	err := commentRouter.DecodeAndValidateComment(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = commentRouter.PublicationData.ChangeAPublicationComment(r.Context(), &commentRouter.Comment)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"comment": commentRouter.Comment,
	})

}

//DeletePublicationComment delete a comment of a publication made by the current logged user
func DeletePublicationComment(w http.ResponseWriter, r *http.Request) {
	commentRouter := newCommentsRouter()
	err := commentRouter.DecodeAndValidateComment(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	err = commentRouter.PublicationData.DeletePublicationComment(r.Context(), &commentRouter.Comment)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"comment": commentRouter.Comment,
	})

}

//CommendLiked like a commend if is not yet in other case unlike
func CommendLiked(w http.ResponseWriter, r *http.Request) {
	commentRouter := newCommentsRouter()
	err := commentRouter.DecodeAndValidateComment(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	newComment, err := commentRouter.PublicationData.CommentLiked(r.Context(), commentRouter.Comment)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"comment": newComment,
	})
}
