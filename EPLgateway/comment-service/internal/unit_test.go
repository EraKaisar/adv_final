package internal

import (
	"EPLgateway/comment-service/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	cm := &model.CommentModel{DB: db}

	// Mock database expectations
	comment := &model.Comment{
		UserID:      1,
		TeamID:      1,
		CommentText: "Test comment",
	}
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, time.Now())

	mock.ExpectQuery(`INSERT INTO comments (.+) RETURNING id, created_at`).
		WithArgs(comment.UserID, comment.TeamID, comment.CommentText).
		WillReturnRows(rows)

	// Call the method and check for expected result
	err = cm.Insert(comment)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	cm := &model.CommentModel{DB: db}

	// Mock database expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "team_id", "comment_text", "created_at"}).
		AddRow(1, 1, 1, "Test comment", time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE team_id = (.+)`).
		WithArgs(1).
		WillReturnRows(rows)

	// Call the method and check for expected result
	comments, err := cm.GetAll(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(comments))
	}
}
func TestInsertRating(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rm := &model.RatingModel{DB: db}

	// Mock database expectations
	rating := &model.Rating{
		UserID: 1,
		TeamID: 1,
		Rating: 5,
	}
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, time.Now())

	mock.ExpectQuery(`INSERT INTO ratings (.+) RETURNING id, created_at`).
		WithArgs(rating.UserID, rating.TeamID, rating.Rating).
		WillReturnRows(rows)

	// Call the method and check for expected result
	err = rm.Insert(rating)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetAllRatings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rm := &model.RatingModel{DB: db}

	// Mock database expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "team_id", "rating", "created_at"}).
		AddRow(1, 1, 1, 5, time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM ratings WHERE team_id = (.+)`).
		WithArgs(1).
		WillReturnRows(rows)

	// Call the method and check for expected result
	ratings, err := rm.GetAll(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ratings) != 1 {
		t.Errorf("Expected 1 rating, got %d", len(ratings))
	}
}
func TestCommentModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	cm := &model.CommentModel{DB: db}

	// Mock database expectations
	comment := &model.Comment{
		ID:          1,
		CommentText: "Updated comment text",
	}

	mock.ExpectExec(`UPDATE comments SET comment_text = (.+) WHERE id = (.+)`).
		WithArgs(comment.CommentText, comment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method and check for expected result
	err = cm.Update(comment)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCommentModel_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	cm := &model.CommentModel{DB: db}

	// Mock database expectations
	commentID := 1
	mock.ExpectExec(`DELETE FROM comments WHERE id = (.+)`).
		WithArgs(commentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method and check for expected result
	err = cm.Delete(commentID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRatingModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rm := &model.RatingModel{DB: db}

	// Mock database expectations
	rating := &model.Rating{
		ID:     1,
		Rating: 4,
	}

	mock.ExpectExec(`UPDATE ratings SET rating = (.+) WHERE id = (.+)`).
		WithArgs(rating.Rating, rating.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method and check for expected result
	err = rm.Update(rating)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRatingModel_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rm := &model.RatingModel{DB: db}

	// Mock database expectations
	ratingID := 1
	mock.ExpectExec(`DELETE FROM ratings WHERE id = (.+)`).
		WithArgs(ratingID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the method and check for expected result
	err = rm.Delete(ratingID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetCommentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	cm := &model.CommentModel{DB: db}

	// Mock database expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "team_id", "comment_text", "created_at"}).
		AddRow(1, 1, 1, "Test comment", time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	// Call the method and check for expected result
	comment, err := cm.GetByID(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if comment == nil || comment.ID != 1 {
		t.Errorf("Expected comment ID 1, got %+v", comment)
	}
}

func TestCountCommentsByTeamID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	cm := &model.CommentModel{DB: db}

	// Mock database expectations
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM comments WHERE team_id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// Call the method and check for expected result
	count, err := cm.CountByTeamID(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}
