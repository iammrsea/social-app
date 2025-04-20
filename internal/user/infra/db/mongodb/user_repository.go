package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/iammrsea/social-app/internal/shared/rbac"
	"github.com/iammrsea/social-app/internal/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// userDocument represents how a user is stored in MongoDB
type userDocument struct {
	ID        string         `bson:"_id"`
	Email     string         `bson:"email"`
	Username  string         `bson:"username"`
	Role      string         `bson:"role"`
	Reputaion userReputation `bson:"reputation"`
	CreatedAt time.Time      `bson:"createdAt"`
	UpdatedAt time.Time      `bson:"updatedAt"`
	BanStatus userBanStatus  `bson:"banStatus"`
}

type userReputation struct {
	ReputationScore int      `bson:"reputationScore"`
	Badges          []string `bson:"badges"`
}

type userBanStatus struct {
	IsBanned        bool      `bson:"isBanned"`
	BannedAt        time.Time `bson:"bannedAt"`
	BanStartDate    time.Time `bson:"banStartDate"`
	BanEndDate      time.Time `bson:"banEndDate"`
	ReasonForBan    string    `bson:"reasonForBan"`
	IsBanIndefinite bool      `bson:"isBanIndefinite"`
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

// BanUser bans a user
func (r *UserRepository) BanUser(ctx context.Context, userId string, updateFn func(user *domain.User) error) error {
	return r.getAndUpdateUser(ctx, userId, updateFn)
}

// UserReadModelMongoRepository implements the domain.UserReadModelRepository interface
type UserReadModelRepository struct {
	collection *mongo.Collection
}

// NewUserReadModelRepository creates a new MongoDB user read model repository
func NewUserReadModelRepository(db *mongo.Database) *UserReadModelRepository {
	return &UserReadModelRepository{
		collection: db.Collection("users"),
	}
}

// GetUserById finds a user by their ID
func (r *UserReadModelRepository) GetUserById(ctx context.Context, userId string) (*domain.UserReadModel, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return documentToReadModel(doc), nil
}

// GetUserByEmail finds a user by their email
func (r *UserReadModelRepository) GetUserByEmail(ctx context.Context, userId string) (*domain.UserReadModel, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{"email": userId}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return documentToReadModel(doc), nil
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

// documentToReadModel converts userDocument to UserReadModel
func documentToReadModel(doc userDocument) *domain.UserReadModel {
	return &domain.UserReadModel{
		Username:  doc.Username,
		Email:     doc.Email,
		Id:        doc.ID,
		Role:      rbac.UserRole(doc.Role),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Reputation: domain.UserReputation{
			ReputationScore: doc.Reputaion.ReputationScore,
			Badges:          doc.Reputaion.Badges,
		},
		BanStatus: domain.BanStatus{
			IsBanned:        doc.BanStatus.IsBanned,
			BannedAt:        doc.BanStatus.BannedAt,
			BanStartDate:    doc.BanStatus.BanStartDate,
			BanEndDate:      doc.BanStatus.BanEndDate,
			ReasonForBan:    doc.BanStatus.ReasonForBan,
			IsBanIndefinite: doc.BanStatus.IsBanIndefinite,
		},
	}
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

// fromDomain converts a domain User to userDocument
func fromDomain(user domain.User) userDocument {
	return userDocument{
		ID:       user.Id(),
		Email:    user.Email(),
		Username: user.Username(),
		Role:     user.Role().String(),
		Reputaion: userReputation{
			Badges:          user.Badges(),
			ReputationScore: user.ReputationScore(),
		},
		CreatedAt: user.JoinedAt(),
		UpdatedAt: user.UpdatedAt(),
		BanStatus: userBanStatus{
			IsBanned:        user.IsBanned(),
			BannedAt:        user.BannedAt(),
			BanStartDate:    user.BanStartDate(),
			BanEndDate:      user.BanEndDate(),
			IsBanIndefinite: user.IsBanIndefinitely(),
		},
	}
}

// toDomain converts a userDocument to domain User
func (u userDocument) toDomain() domain.User {
	return domain.MustNewUser(
		u.ID,
		u.Email,
		u.Username,
		rbac.UserRole(u.Role),
		u.CreatedAt,
		u.UpdatedAt,
		domain.MustNewUserReputation(u.Reputaion.ReputationScore, u.Reputaion.Badges),
	)
}
