package data

import (
	"EPLgateway/auth-service/internal/data"
	"crypto/sha256"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func TestTokenModel_New(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tm := data.TokenModel{DB: db}

	userID := int64(1)
	ttl := 24 * time.Hour
	scope := data.ScopeAuthentication

	mock.ExpectExec("INSERT INTO tokens").
		WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg(), scope).
		WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := tm.New(userID, ttl, scope)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if token == nil {
		t.Error("expected token to be generated, got nil")
	}

	if len(token.Plaintext) != 26 {
		t.Errorf("expected token plaintext length to be 26, got %d", len(token.Plaintext))
	}
}

func TestTokenModel_DeleteAllForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tm := data.TokenModel{DB: db}

	userID := int64(1)
	scope := data.ScopeAuthentication

	mock.ExpectExec("DELETE FROM tokens").
		WithArgs(scope, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = tm.DeleteAllForUser(scope, userID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestPermissionModel_GetAllForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pm := data.PermissionModel{DB: db}

	userID := int64(1)
	expectedPermissions := data.Permissions{"permission1", "permission2"}

	rows := sqlmock.NewRows([]string{"code"}).
		AddRow("permission1").
		AddRow("permission2")

	mock.ExpectQuery("SELECT permissions.code").
		WithArgs(userID).
		WillReturnRows(rows)

	permissions, err := pm.GetAllForUser(userID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(permissions) != len(expectedPermissions) {
		t.Errorf("expected %d permissions, got %d", len(expectedPermissions), len(permissions))
	}

	for i, permission := range permissions {
		if permission != expectedPermissions[i] {
			t.Errorf("expected permission %s, got %s", expectedPermissions[i], permission)
		}
	}
}

func TestPermissionModel_AddForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pm := data.PermissionModel{DB: db}

	userID := int64(1)
	codes := []string{"permission1", "permission2"}

	mock.ExpectExec("INSERT INTO users_permissions").
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 2))

	err = pm.AddForUser(userID, codes...)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
func TestPermissionModel_GetAllForUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pm := data.PermissionModel{DB: db}

	userID := int64(1)

	mock.ExpectQuery("SELECT permissions.code").
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	_, err = pm.GetAllForUser(userID)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}
func TestPermissionModel_AddForUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pm := data.PermissionModel{DB: db}

	userID := int64(1)
	codes := []string{"permission1", "permission2"}

	mock.ExpectExec("INSERT INTO users_permissions").
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnError(errors.New("database error"))

	err = pm.AddForUser(userID, codes...)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}
func TestUserModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userModel := data.UserModel{DB: db}

	user := &data.UserInfo{
		Name:  "Test User",
		Email: "test@example.com",
		Password: data.Password{
			Plaintext: new(string),
			Hash:      []byte("hashedpassword"),
		},
		Activated: true,
	}

	mock.ExpectQuery("INSERT INTO user_info").
		WithArgs(user.Name, user.Email, user.Password.Hash, user.Activated).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}).AddRow(1, time.Now(), 1))

	err = userModel.Insert(user)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", user.ID)
	}
}
func TestUserModel_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userModel := data.UserModel{DB: db}

	email := "test@example.com"
	expectedUser := &data.UserInfo{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Test User",
		Email:     email,
		Password: data.Password{
			Hash: []byte("hashedpassword"),
		},
		Activated: true,
		Version:   1,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "name", "email", "password_hash", "activated", "version"}).
		AddRow(expectedUser.ID, expectedUser.CreatedAt, expectedUser.Name, expectedUser.Email, expectedUser.Password.Hash, expectedUser.Activated, expectedUser.Version)

	mock.ExpectQuery("SELECT id, created_at, name, email, password_hash, activated, version FROM user_info").
		WithArgs(email).
		WillReturnRows(rows)

	user, err := userModel.GetByEmail(email)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.Email != expectedUser.Email {
		t.Errorf("expected email %s, got %s", expectedUser.Email, user.Email)
	}
}
func TestUserModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userModel := data.UserModel{DB: db}

	user := &data.UserInfo{
		ID:        1,
		Name:      "Updated User",
		Email:     "updated@example.com",
		Password:  data.Password{Hash: []byte("updatedhashedpassword")},
		Activated: true,
		Version:   1,
	}

	mock.ExpectQuery("UPDATE user_info").
		WithArgs(user.Name, user.Email, user.Password.Hash, user.Activated, user.ID, user.Version).
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(2))

	err = userModel.Update(user)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.Version != 2 {
		t.Errorf("expected version 2, got %d", user.Version)
	}
}
func TestUserModel_GetForToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userModel := data.UserModel{DB: db}

	tokenScope := data.ScopeAuthentication
	tokenPlaintext := "26-character-long-token"
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	expectedUser := &data.UserInfo{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  data.Password{Hash: []byte("hashedpassword")},
		Activated: true,
		Version:   1,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "name", "email", "password_hash", "activated", "version"}).
		AddRow(expectedUser.ID, expectedUser.CreatedAt, expectedUser.Name, expectedUser.Email, expectedUser.Password.Hash, expectedUser.Activated, expectedUser.Version)

	mock.ExpectQuery("SELECT user_info.id, user_info.created_at, user_info.name, user_info.email, user_info.password_hash, user_info.activated, user_info.version FROM user_info INNER JOIN tokens").
		WithArgs(tokenHash[:], tokenScope, sqlmock.AnyArg()).
		WillReturnRows(rows)

	user, err := userModel.GetForToken(tokenScope, tokenPlaintext)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.Email != expectedUser.Email {
		t.Errorf("expected email %s, got %s", expectedUser.Email, user.Email)
	}
}
