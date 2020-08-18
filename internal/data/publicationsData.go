package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Publication has the field of a publication that make a user
type Publication struct {
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content  string             `json:"content" bson:"content,"`
	Likes    []Like             `json:"likes,omitempty" bson:"likes"`
	Comments []Comment          `json:"commets,omitempty" bson:"commets"`
}

//UserPublications params for db connection
type UserPublications struct {
	data *Data
	coll *mongo.Collection
}

//IsValidFields validate the field of this data
func (p *Publication) IsValidFields() bool {
	return p.Content != ""

}

//NewPublication add a new publication
func (up *UserPublications) NewPublication(ctx context.Context, publication *Publication) error {
	_, err := up.coll.InsertOne(ctx, *publication)
	return err

}

//FindPublicationByID find a publication by user_id
func (up *UserPublications) FindPublicationByID(ctx context.Context, ID primitive.ObjectID) (Publication, error) {
	var publication Publication
	result := up.coll.FindOne(ctx, bson.M{"_id": ID})
	err := result.Decode(&publication)
	return publication, err

}

//GetAllUserPublications get all publications for a the logged user
func (up *UserPublications) GetAllUserPublications(ctx context.Context) (*[]Publication, error) {
	publications := []Publication{}
	cursor, err := up.coll.Find(ctx, bson.M{"user_id": ctx.Value(primitive.ObjectID{})})
	if err != nil {
		return &publications, err
	}
	cursor.All(ctx, &publications)
	return &publications, err
}

//UpdatePublication update a publication
func (up *UserPublications) UpdatePublication(ctx context.Context, publication *Publication) error {
	publication.ID = primitive.NilObjectID
	_, err := up.coll.UpdateOne(ctx, bson.M{"user_id": &publication.UserID}, bson.M{"$set": &publication})
	return err

}

//DeletePublication delete a publication
func (up *UserPublications) DeletePublication(ctx context.Context, publication Publication) error {
	_, err := up.coll.DeleteOne(ctx, bson.M{"user_id": publication.UserID})
	return err

}

//PublicationLiked like a publication
func (up *UserPublications) PublicationLiked(ctx context.Context, publicationID primitive.ObjectID) (Publication, error) {
	publication, err := up.FindPublicationByID(ctx, publicationID)
	if err != nil {
		return publication, err
	}
	l := Like{}

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
