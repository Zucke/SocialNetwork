package handlers

import (
	"context"
	"net/http"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/authentication"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//NewPublication create a publication for a user
func NewPublication(w http.ResponseWriter, r *http.Request) {
	publication := data.Publication{}
	publication.ID = primitive.NewObjectID()
	if !authentication.BasicValidations(&publication, w, r) {
		response.HTTPError(w, r, http.StatusBadRequest, "Bad information")
		return
	}

	publication.UserID = (r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID))
	newPublication := data.NewUserPublications()
	err := newPublication.NewPublication(r.Context(), &publication)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"publication": publication,
	})
}

//GetUserPublication get a publication from a user
func GetUserPublication(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	newPublication := data.NewUserPublications()
	ctx := context.Background()
	publication, err := newPublication.FindPublicationByID(ctx, publicationID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"publication": publication,
	})

}

//ChangePublication update a publication of logged user
func ChangePublication(w http.ResponseWriter, r *http.Request) {

	var updatedPublication data.Publication
	if !authentication.BasicValidations(&updatedPublication, w, r) {
		return
	}

	if updatedPublication.UserID != r.Context().Value(primitive.ObjectID{}) {
		response.HTTPError(w, r, http.StatusBadRequest, data.ErrorAccesDenied.Error())
		return
	}

	newPublication := data.NewUserPublications()
	ctx := context.Background()
	newPublication.UpdatePublication(ctx, &updatedPublication)

	render.JSON(w, r, render.M{
		"publication": updatedPublication,
	})

}

//DeletePublication delete a publication of logged user
func DeletePublication(w http.ResponseWriter, r *http.Request) {

	var publication data.Publication
	if !authentication.BasicValidations(&publication, w, r) {
		return
	}

	if publication.UserID != r.Context().Value(primitive.ObjectID{}) {
		response.HTTPError(w, r, http.StatusBadRequest, data.ErrorAccesDenied.Error())
		return
	}

	newPublication := data.NewUserPublications()
	ctx := context.Background()
	newPublication.DeletePublication(ctx, publication)

	render.JSON(w, r, render.M{
		"publication": publication,
	})

}

//NewPublicationLike like or dislike from a user to a publication
func NewPublicationLike(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	newPublication := data.NewUserPublications()

	publication, err := newPublication.PublicationLiked(r.Context(), publicationID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"publication": publication,
	})
}

//GetAllUserPublication get all publications for a user
func GetAllUserPublication(w http.ResponseWriter, r *http.Request) {
	newPublication := data.NewUserPublications()
	publications, err := newPublication.GetAllUserPublications(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"publications": publications,
	})

}
