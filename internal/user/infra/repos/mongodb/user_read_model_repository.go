package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iammrsea/social-app/internal/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserReadModelMongoRepository implements the domain.UserReadModelRepository interface
type UserReadModelRepository struct {
	collection *mongo.Collection
}

func NewUserReadModelRepository(db *mongo.Database) *UserReadModelRepository {
	return &UserReadModelRepository{
		collection: db.Collection("users"),
	}
}

// GetUserByEmail finds a user by their email
func (r *UserReadModelRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserReadModel, error) {
	doc, err := r.getUserBy(ctx, "email", email)
	if err != nil {
		return nil, err
	}
	return documentToReadModel(*doc), nil
}

// GetUserByEmail finds a user by their id
func (r *UserReadModelRepository) GetUserById(ctx context.Context, id string) (*domain.UserReadModel, error) {
	doc, err := r.getUserBy(ctx, "_id", id)
	if err != nil {
		return nil, err
	}
	return documentToReadModel(*doc), nil
}

// GetUsers retrieves paginated users sorted by createdAt
func (r *UserReadModelRepository) GetUsers(ctx context.Context, opts domain.GetUsersOptions) ([]*domain.UserReadModel, bool, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.First + 1))                //Fetch one more to determine if there are more results
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}}) //Sort by createdAt descending (newest first)

	var filter bson.M = bson.M{}
	if opts.After != "" {
		// Parse the after cursor (which is a timestamp) for cursor-based pagination
		createdAt, err := time.Parse(time.RFC3339Nano, opts.After)
		if err != nil {
			return nil, false, err
		}
		filter = bson.M{"createdAt": bson.M{"$lt": createdAt}}
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)

	if err != nil {
		return nil, false, err
	}
	defer cursor.Close(ctx)

	var users []*domain.UserReadModel
	var docs []userDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, false, err
	}
	// Check if there are more results than requested
	hasNex := false
	if len(docs) > int(opts.First) {
		hasNex = true
		docs = docs[:opts.First]
	}
	// Convert documents to domain read models
	for _, doc := range docs {
		users = append(users, documentToReadModel(doc))
	}
	return users, hasNex, nil
}

// GetUserBy finds a user by the specified field and value
func (r *UserReadModelRepository) getUserBy(ctx context.Context, fieldName string, value any) (*userDocument, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{fmt.Sprintf("%s", fieldName): value}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &doc, nil
}
