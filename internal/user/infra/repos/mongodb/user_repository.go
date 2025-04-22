package mongodb

import (
	"context"
	"errors"

	"github.com/iammrsea/social-app/internal/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository implements the domain.UserRepository interface
type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Register adds a new user to the database
func (r *UserRepository) Register(ctx context.Context, user domain.User) error {
	doc := fromDomain(user)
	// Check if user with the same email already exists
	existingUser := userDocument{}
	err := r.collection.FindOne(ctx, bson.M{"email": user.Email()}).Decode(&existingUser)
	if err == nil {
		return domain.ErrEmailAlreadyExists
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	_, err = r.collection.InsertOne(ctx, doc)

	return err
}

// MakeModerator updates a user to have moderator role
func (r *UserRepository) MakeModerator(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// AwardBadge adds a badge to a user
func (r *UserRepository) AwardBadge(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// RevokeAwardedBadge removes a badge from a user
func (r *UserRepository) RevokeAwardedBadge(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// ChangeUsername updates a user's username
func (r *UserRepository) ChangeUsername(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// UnbanUser unbans a user
func (r *UserRepository) UnbanUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// BanUser bans a user
func (r *UserRepository) BanUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// GetUserBy finds a user by a field name and value
func (r *UserRepository) GetUserBy(ctx context.Context, fieldName string, value any) (*domain.User, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{fieldName: value}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	user := doc.toDomain()
	return &user, nil
}

// getAndUpdateUser is a helper function for updating user documents
func (r *UserRepository) getAndUpdateUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	// Start a session and transaction
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo.SessionContext) (any, error) {
		// Get the current user
		var doc userDocument
		err := r.collection.FindOne(sessionCtx, bson.M{"_id": userId}).Decode(&doc)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, domain.ErrUserNotFound
			}
			return nil, err
		}
		// Convert to domain model
		user := doc.toDomain()
		//Apply the update function
		err = updateFn(&user)
		if err != nil {
			return nil, err
		}
		//Convert back to document and update
		updatedDoc := fromDomain(user)
		_, err = r.collection.ReplaceOne(sessionCtx, bson.M{"_id": userId}, updatedDoc)
		return nil, err
	}
	_, err = session.WithTransaction(ctx, callback)

	return err
}
