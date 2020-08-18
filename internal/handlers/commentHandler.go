package handlers

import (
	"net/http"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/authentication"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//NewPublicationComment add a comment of the current logged user
func NewPublicationComment(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	var comment data.Comment
	if !authentication.BasicValidations(&comment, w, r) {
		response.HTTPError(w, r, http.StatusBadRequest, data.ErrorBadInfo.Error())
		return
	}

	newPublication := data.NewUserPublications()
	comment.PublicationID = publicationID
	comment.UserID = r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID)
	newPublication.NewPublicationComment(r.Context(), &comment)

	render.JSON(w, r, render.M{
		"comment": comment,
	})

}

//ChangePublicationComment add a comment of the current logged user
func ChangePublicationComment(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	var comment data.Comment
	if !authentication.BasicValidations(&comment, w, r) {
		return
	}

	comment.PublicationID = publicationID
	newPublication := data.NewUserPublications()
	err = newPublication.ChangeAPublicationComment(r.Context(), &comment)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"comment": comment,
	})

}

//DeletePublicationComment delete a comment of a publication made by the current logged user
func DeletePublicationComment(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var comment data.Comment
	if !authentication.BasicValidations(&comment, w, r) {
		return
	}

	comment.PublicationID = publicationID
	newPublication := data.NewUserPublications()
	err = newPublication.DeletePublicationComment(r.Context(), &comment)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"comment": comment,
	})

}

//CommendLiked like a commend if is not yet in other case unlike
func CommendLiked(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	var comment data.Comment
	if !authentication.BasicValidations(&comment, w, r) {
		return
	}

	comment.PublicationID = publicationID
	newPublication := data.NewUserPublications()
	newComment, err := newPublication.CommentLiked(r.Context(), comment)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"comment": newComment,
	})
}
