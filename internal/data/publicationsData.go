package data

import (
	"context"

	"github.com/Zucke/SocialNetwork/pkg/likes"
	"github.com/Zucke/SocialNetwork/pkg/publications"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//UserPublications params for db connection
type UserPublications struct {
	data *Data
	coll *mongo.Collection
}

//NewPublication add a new publication
func (up *UserPublications) NewPublication(ctx context.Context, publication *publications.Publication) error {
	_, err := up.coll.InsertOne(ctx, *publication)
	return err

}

//FindPublicationByID find a publication by user_id
func (up *UserPublications) FindPublicationByID(ctx context.Context, ID primitive.ObjectID) (publications.Publication, error) {
	var publication publications.Publication
	result := up.coll.FindOne(ctx, bson.M{"_id": ID})
	err := result.Decode(&publication)
	return publication, err

}

//GetAllUserPublications get all publications for a the logged user
func (up *UserPublications) GetAllUserPublications(ctx context.Context) (*[]publications.Publication, error) {
	publications := []publications.Publication{}
	cursor, err := up.coll.Find(ctx, bson.M{"user_id": ctx.Value(primitive.ObjectID{})})
	if err != nil {
		return &publications, err
	}
	cursor.All(ctx, &publications)
	return &publications, err
}

//UpdatePublication update a publication
func (up *UserPublications) UpdatePublication(ctx context.Context, publication *publications.Publication) error {
	publication.ID = primitive.NilObjectID
	_, err := up.coll.UpdateOne(ctx, bson.M{"user_id": &publication.UserID}, bson.M{"$set": &publication})
	return err

}

//DeletePublication delete a publication
func (up *UserPublications) DeletePublication(ctx context.Context, publication publications.Publication) error {
	_, err := up.coll.DeleteOne(ctx, bson.M{"user_id": publication.UserID})
	return err

}

//PublicationLiked like a publication
func (up *UserPublications) PublicationLiked(ctx context.Context, publicationID primitive.ObjectID) (publications.Publication, error) {
	publication, err := up.FindPublicationByID(ctx, publicationID)
	if err != nil {
		return publication, err
	}
	l := likes.Like{}

	publication.Likes = *l.AppendLike(ctx, publication.Likes)
	return publication, up.UpdatePublication(ctx, &publication)

}

//NewUserPublications return db info
func NewUserPublications() UserPublications {
	return UserPublications{
		data: New(),
		coll: data.DBCollection(PublicationsColletion),
	}
}
