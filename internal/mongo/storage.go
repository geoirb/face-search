package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	search "github.com/geoirb/face-search/internal/face-search"
)

// Storage mongo.
type Storage struct {
	faceSearchCollection *mongo.Collection
}

// NewStorage ...
func NewStorage(
	ctx context.Context,
	connStr, faceSearchDatabaseName, faceSearchCollectionName string,
) (*Storage, error) {
	opts := options.Client().ApplyURI(connStr)
	connect, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err = connect.Ping(ctx, nil); err != nil {
		err = fmt.Errorf("error ping mongo storage %w", err)
	}

	collection := connect.Database(faceSearchDatabaseName).Collection(faceSearchCollectionName)
	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"photo_hash": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return &Storage{
		faceSearchCollection: collection,
	}, err
}

// Save face search result.
func (m *Storage) Save(ctx context.Context, result search.Result) (err error) {
	faceSearch := faceSearch{
		Status:    result.Status,
		Error:     result.Error,
		UUID:      result.UUID,
		PhotoHash: result.PhotoHash,
		Profiles:  make([]profile, 0, len(result.Profiles)),
		UpdateAt:  result.UpdateAt,
		CreateAt:  result.CreateAt,
	}
	for _, p := range result.Profiles {
		faceSearch.Profiles = append(faceSearch.Profiles, profile(p))
	}

	_, err = m.faceSearchCollection.InsertOne(ctx, faceSearch)
	return
}

func (m *Storage) Update(ctx context.Context, result search.Result) (err error) {
	faceSearch := faceSearch{
		Status:    result.Status,
		Error:     result.Error,
		UUID:      result.UUID,
		PhotoHash: result.PhotoHash,
		Profiles:  make([]profile, 0, len(result.Profiles)),
		UpdateAt:  result.UpdateAt,
		CreateAt:  result.CreateAt,
	}
	for _, p := range result.Profiles {
		faceSearch.Profiles = append(faceSearch.Profiles, profile(p))
	}
	filter := bson.M{
		"uuid": result.UUID,
	}
	update := bson.M{
		"$set": faceSearch,
	}
	_, err = m.faceSearchCollection.UpdateOne(ctx, filter, update)
	return
}

// Get face search result by filter.
func (m *Storage) Get(ctx context.Context, filter search.ResultFilter) (search.Result, error) {
	f := make(bson.M)
	if filter.UUID != nil {
		f["uuid"] = *filter.UUID
	}

	if filter.PhotoHash != nil {
		f["photo_hash"] = *filter.PhotoHash
	}

	var faceSearchResult faceSearch
	if err := m.faceSearchCollection.FindOne(ctx, f).Decode(&faceSearchResult); err != nil {
		if mongo.ErrNoDocuments == err {
			err = search.ErrFaceSearchResultNotFound
		}
		return search.Result{}, err
	}

	res := search.Result{
		Status:    faceSearchResult.Status,
		Error:     faceSearchResult.Error,
		UUID:      faceSearchResult.UUID,
		PhotoHash: faceSearchResult.PhotoHash,
		Profiles:  m.convertProfileFromMongo(faceSearchResult.Profiles),
		UpdateAt:  faceSearchResult.UpdateAt,
		CreateAt:  faceSearchResult.CreateAt,
	}
	return res, nil
}

func (m *Storage) convertProfileFromMongo(src []profile) []search.Profile {
	dst := make([]search.Profile, 0, len(src))
	for _, p := range src {
		dst = append(dst, search.Profile(p))
	}
	return dst
}
