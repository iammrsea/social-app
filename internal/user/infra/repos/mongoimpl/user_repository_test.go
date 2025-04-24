package mongoimpl

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/shared/guards/rbac"
	"github.com/iammrsea/social-app/internal/testutil"
	"github.com/iammrsea/social-app/internal/user/domain"
	"github.com/lucsky/cuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUserRepository(t *testing.T) {
	t.Parallel()
	client := testutil.SetupTestMongoDb(t)
	require.NotNil(t, client)

	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		testRegisterUser(t, client)
	})

	t.Run("ChangeUsername", func(t *testing.T) {
		t.Parallel()
		testChangeUsername(t, client)
	})

	t.Run("AwardBadge", func(t *testing.T) {
		t.Parallel()
		testAwardBadge(t, client)
	})

	t.Run("BanUser", func(t *testing.T) {
		t.Parallel()
		testBanUser(t, client)
	})

	t.Run("UnbanUser", func(t *testing.T) {
		t.Parallel()
		testUnbanUser(t, client)
	})

	t.Run("MakeModerator", func(t *testing.T) {
		t.Parallel()
		testMakeModerator(t, client)
	})

	t.Run("RevokeAwardedBadge", func(t *testing.T) {
		t.Parallel()
		testRevokeAwardedBadge(t, client)
	})

	t.Run("UserExists", func(t *testing.T) {
		t.Parallel()
		testUserExists(t, client)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		t.Parallel()
		testGetUserByEmail(t, client)
	})

	t.Run("GetUserById", func(t *testing.T) {
		t.Parallel()
		testGetUserById(t, client)
	})

	t.Run("GetUsers", func(t *testing.T) {
		t.Parallel()
		testGetUsers(t, client)
	})

}

func testGetUsers(t *testing.T, client *mongo.Client) {
	t.Skip("skipping testGetUsers")
	panic("unimplemented")
}

func testGetUserById(t *testing.T, client *mongo.Client) {
	t.Skip("skipping testGetUserById")
	panic("unimplemented")
}

func testGetUserByEmail(t *testing.T, client *mongo.Client) {
	t.Skip("skipping testGetUserByEmail")
	panic("unimplemented")
}

func testUserExists(t *testing.T, client *mongo.Client) {
	t.Run("UserExists", func(t *testing.T) {
		t.Parallel()
		db, cleanup := setUpDB(t, client)
		defer cleanup()
		user := addUserToDB(t, db, nil)
		repo := getUserRepo(t, db)

		exists, err := repo.UserExists(context.Background(), user.Email(), user.Username())
		require.NoError(t, err)
		assert.True(t, exists)
	})
	t.Run("UserDoesNotExist", func(t *testing.T) {
		t.Parallel()
		db, cleanup := setUpDB(t, client)
		defer cleanup()

		repo := getUserRepo(t, db)
		exists, err := repo.UserExists(context.Background(), "nonexisting@gmail.com", "non_existing_username")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func testRevokeAwardedBadge(t *testing.T, client *mongo.Client) {
	db, cleanup := setUpDB(t, client)
	defer cleanup()
	badge := "new_badge"
	u, err := domain.NewUser(
		cuid.New(),
		fmt.Sprintf("%s@example.com", cuid.New()),
		fmt.Sprintf("username_%s", cuid.New()),
		rbac.Regular,
		time.Now(),
		time.Now(),
		domain.MustNewUserReputation(0, []string{badge}),
		nil)
	require.NoError(t, err)
	user := addUserToDB(t, db, &u)
	repo := getUserRepo(t, db)

	err = repo.RevokeAwardedBadge(context.Background(), user.Id(), func(u *domain.User) error {
		return u.RevokeAwardedBadge(badge)
	})
	require.NoError(t, err)
	updatedUser := getUser(t, db, user.Id())
	assert.NotContains(t, updatedUser.Reputaion.Badges, badge)
}

func testMakeModerator(t *testing.T, client *mongo.Client) {
	db, cleanup := setUpDB(t, client)
	defer cleanup()

	user := addUserToDB(t, db, nil)
	repo := getUserRepo(t, db)

	err := repo.MakeModerator(context.Background(), user.Id(), func(u *domain.User) error {
		return u.MakeModerator()
	})
	require.NoError(t, err)
	updatedUser := getUser(t, db, user.Id())
	assert.Equal(t, rbac.Moderator.String(), updatedUser.Role)
}

func testUnbanUser(t *testing.T, client *mongo.Client) {
	db, cleanup := setUpDB(t, client)
	defer cleanup()

	u, err := domain.NewUser(
		cuid.New(),
		fmt.Sprintf("%s@example.com", cuid.New()),
		fmt.Sprintf("username_%s", cuid.New()),
		rbac.Regular,
		time.Now(),
		time.Now(),
		nil,
		domain.NewBan(true, "spamming", true, time.Now(), time.Now(), time.Now()))

	require.NoError(t, err)
	user := addUserToDB(t, db, &u)
	repo := getUserRepo(t, db)

	err = repo.UnbanUser(context.Background(), user.Id(), func(u *domain.User) error {
		return u.UnBan()
	})
	require.NoError(t, err)
	updatedUser := getUser(t, db, user.Id())
	assert.False(t, updatedUser.BanStatus.IsBanned)
	assert.Equal(t, "", updatedUser.BanStatus.ReasonForBan)
	assert.False(t, updatedUser.BanStatus.IsBanIndefinite)
}

func testRegisterUser(t *testing.T, client *mongo.Client) {
	db, cleanup := setUpDB(t, client)
	defer cleanup()

	user, err := domain.NewUser(
		cuid.New(),
		fmt.Sprintf("%s@example.com", cuid.New()),
		fmt.Sprintf("username_%s", cuid.New()),
		rbac.Regular,
		time.Now(),
		time.Now(),
		nil,
		nil)

	require.NoError(t, err)

	repo := NewUserRepository(db)
	require.NotNil(t, repo)
	err = repo.Register(context.Background(), user)
	require.NoError(t, err)
	assertUserInDB(t, db, user)
}
func testChangeUsername(t *testing.T, client *mongo.Client) {
	db, cleanup := setUpDB(t, client)
	defer cleanup()

	user := addUserToDB(t, db, nil)
	repo := getUserRepo(t, db)
	newUsername := fmt.Sprintf("new_username_%s", cuid.New())
	err := repo.ChangeUsername(context.Background(), user.Id(), func(u *domain.User) error {
		return u.ChangeUsername(newUsername)
	})
	require.NoError(t, err)
	updatedUser := getUser(t, db, user.Id())
	assert.Equal(t, newUsername, updatedUser.Username)
}
func testAwardBadge(t *testing.T, client *mongo.Client) {
	db, cleanup := setUpDB(t, client)
	defer cleanup()
	user := addUserToDB(t, db, nil)
	repo := getUserRepo(t, db)

	badge := "new_badge"
	err := repo.AwardBadge(context.Background(), user.Id(), func(u *domain.User) error {
		return u.AwardBadge(badge)
	})
	require.NoError(t, err)
	updatedUser := getUser(t, db, user.Id())
	assert.Contains(t, updatedUser.Reputaion.Badges, badge)
}
func testBanUser(t *testing.T, client *mongo.Client) {
	t.Helper()
	db, cleanup := setUpDB(t, client)
	defer cleanup()

	t.Run("Ban user indefinitely", func(t *testing.T) {
		t.Parallel()
		user := addUserToDB(t, db, nil)
		repo := getUserRepo(t, db)

		err := repo.BanUser(context.Background(), user.Id(), func(u *domain.User) error {
			return u.Ban("spamming", true, nil)
		})
		require.NoError(t, err)
		updatedUser := getUser(t, db, user.Id())
		assert.True(t, updatedUser.BanStatus.IsBanned)
		assert.Equal(t, "spamming", updatedUser.BanStatus.ReasonForBan)
		assert.True(t, updatedUser.BanStatus.IsBanIndefinite)
	})
	t.Run("Ban user with start and end date", func(t *testing.T) {
		t.Parallel()
		user := addUserToDB(t, db, nil)
		repo := getUserRepo(t, db)

		startDate := time.Now()
		endDate := startDate.Add(24 * time.Hour)
		banTimeline := domain.NewBanTimeline(startDate, endDate)
		err := repo.BanUser(context.Background(), user.Id(), func(u *domain.User) error {
			return u.Ban("spamming", false, banTimeline)
		})
		require.NoError(t, err)
		updatedUser := getUser(t, db, user.Id())
		assert.True(t, updatedUser.BanStatus.IsBanned)
		assert.Equal(t, "spamming", updatedUser.BanStatus.ReasonForBan)
		assert.False(t, updatedUser.BanStatus.IsBanIndefinite)
		assert.Equal(t, startDate.UTC().Format(time.RFC3339), updatedUser.BanStatus.BanStartDate.UTC().Format(time.RFC3339))
		assert.Equal(t, endDate.UTC().Format(time.RFC3339), updatedUser.BanStatus.BanEndDate.UTC().Format(time.RFC3339))
	})

}

func addUserToDB(t *testing.T, db *mongo.Database, user *domain.User) domain.User {
	t.Helper()
	if user == nil {
		u, err := domain.NewUser(
			cuid.New(),
			fmt.Sprintf("%s@example.com", cuid.New()),
			fmt.Sprintf("username_%s", cuid.New()),
			rbac.Regular,
			time.Now(),
			time.Now(),
			nil,
			nil)
		require.NoError(t, err)
		user = &u
	}
	doc := fromDomain(*user)
	_, err := db.Collection("users").InsertOne(context.Background(), doc)
	require.NoError(t, err)
	return *user
}

func getUserRepo(t *testing.T, db *mongo.Database) *UserRepository {
	t.Helper()
	repo := NewUserRepository(db)
	require.NotNil(t, repo)
	return repo
}

func setUpDB(t *testing.T, client *mongo.Client) (*mongo.Database, func()) {
	t.Helper()
	ctx := context.Background()
	dbName := fmt.Sprintf("testdb_%s", cuid.New())
	db := client.Database(dbName)

	cleanup := func() {
		err := db.Drop(ctx)
		if err != nil {
			log.Default().Println("Error dropping database: ", err)
		}
		log.Default().Println("Cleaned up database: ", dbName)
	}
	return db, cleanup
}

func assertUserInDB(t *testing.T, db *mongo.Database, user domain.User) {
	var foundUser userDocument
	err := db.Collection("users").FindOne(context.Background(), bson.M{"_id": user.Id()}).Decode(&foundUser)
	require.NoError(t, err)
	assert.Equal(t, user.Id(), foundUser.ID)
	assert.Equal(t, user.Email(), foundUser.Email)
	assert.Equal(t, user.Username(), foundUser.Username)
	assert.Equal(t, user.Role().String(), foundUser.Role)
	assert.Equal(t, user.ReputationScore(), foundUser.Reputaion.ReputationScore)
	assert.Equal(t, user.JoinedAt().UTC().Format(time.RFC3339), foundUser.CreatedAt.UTC().Format(time.RFC3339))
	assert.Equal(t, user.UpdatedAt().UTC().Format(time.RFC3339), foundUser.UpdatedAt.UTC().Format(time.RFC3339))
	assert.Equal(t, user.IsBanned(), foundUser.BanStatus.IsBanned)
	assert.Equal(t, user.BannedAt(), foundUser.BanStatus.BannedAt)
	assert.Equal(t, user.BanStartDate(), foundUser.BanStatus.BanStartDate)
	assert.Equal(t, user.BanEndDate(), foundUser.BanStatus.BanEndDate)
	assert.Equal(t, user.IsBanIndefinite(), foundUser.BanStatus.IsBanIndefinite)
	assert.Equal(t, user.ReasonForBan(), foundUser.BanStatus.ReasonForBan)
	assert.Equal(t, user.Badges(), foundUser.Reputaion.Badges)
}

func getUser(t *testing.T, db *mongo.Database, userId string) userDocument {
	var foundUser userDocument
	err := db.Collection("users").FindOne(context.Background(), bson.M{"_id": userId}).Decode(&foundUser)
	require.NoError(t, err)
	return foundUser
}
