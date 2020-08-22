package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Zucke/SocialNetwork/internal/data"
	"github.com/Zucke/SocialNetwork/pkg/errorstatus"
	"github.com/Zucke/SocialNetwork/pkg/publications"
	"github.com/Zucke/SocialNetwork/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//PublicationsRouter has the data for a publication and the db connection
type PublicationsRouter struct {
	Publication     publications.Publication
	PublicationData data.UserPublications
}

func newPublicationsRouter() *PublicationsRouter {
	return &PublicationsRouter{
		PublicationData: data.NewUserPublications(),
	}
}

//DecodeAndValidatePublication the request body to the publication, validate the information and return a error if exist
func (p *PublicationsRouter) DecodeAndValidatePublication(w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&p.Publication)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return p.Publication.IsValidFields()

}

//NewPublication create a publication for a user
func NewPublication(w http.ResponseWriter, r *http.Request) {
	publicationRouter := newPublicationsRouter()
	publicationRouter.Publication.ID = primitive.NewObjectID()
	err := publicationRouter.DecodeAndValidatePublication(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	publicationRouter.Publication.UserID = (r.Context().Value(primitive.ObjectID{}).(primitive.ObjectID))
	err = publicationRouter.PublicationData.NewPublication(r.Context(), &publicationRouter.Publication)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"publication": publicationRouter.Publication,
	})
}

//GetUserPublication get a publication from a user
func GetUserPublication(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	publicationRouter := newPublicationsRouter()
	ctx := context.Background()
	publicationRouter.Publication, err = publicationRouter.PublicationData.FindPublicationByID(ctx, publicationID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"publication": publicationRouter.Publication,
	})

}

//ChangePublication update a publication of logged user
func ChangePublication(w http.ResponseWriter, r *http.Request) {
	publicationRouter := newPublicationsRouter()
	err := publicationRouter.DecodeAndValidatePublication(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	if publicationRouter.Publication.UserID != r.Context().Value(primitive.ObjectID{}) {
		response.HTTPError(w, r, http.StatusBadRequest, errorstatus.ErrorAccesDenied.Error())
		return
	}
	ctx := context.Background()
	publicationRouter.PublicationData.UpdatePublication(ctx, &publicationRouter.Publication)

	render.JSON(w, r, render.M{
		"publication": publicationRouter.Publication,
	})

}

//DeletePublication delete a publication of logged user
func DeletePublication(w http.ResponseWriter, r *http.Request) {
	publicationRouter := newPublicationsRouter()
	err := publicationRouter.DecodeAndValidatePublication(w, r)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if publicationRouter.Publication.UserID != r.Context().Value(primitive.ObjectID{}) {
		response.HTTPError(w, r, http.StatusBadRequest, errorstatus.ErrorAccesDenied.Error())
		return
	}

	ctx := context.Background()
	publicationRouter.PublicationData.DeletePublication(ctx, publicationRouter.Publication)

	render.JSON(w, r, render.M{
		"publication": publicationRouter.Publication,
	})

}

//NewPublicationLike like or dislike from a user to a publication
func NewPublicationLike(w http.ResponseWriter, r *http.Request) {
	publicationID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "publication_id"))
	publicationRouter := newPublicationsRouter()
	publicationRouter.Publication, err = publicationRouter.PublicationData.PublicationLiked(r.Context(), publicationID)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, render.M{
		"publication": publicationRouter.Publication,
	})
}

//GetAllUserPublication get all publications for a user
func GetAllUserPublication(w http.ResponseWriter, r *http.Request) {
	publicationRouter := newPublicationsRouter()
	publications, err := publicationRouter.PublicationData.GetAllUserPublications(r.Context())
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.JSON(w, r, render.M{
		"publications": publications,
	})

}
