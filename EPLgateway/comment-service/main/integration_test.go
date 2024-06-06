package main

import (
	"EPLgateway/comment-service/internal/model"
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApplication(t *testing.T) (*application, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	app := &application{
		models: model.Models{
			Comments: model.CommentModel{DB: db},
			Ratings:  model.RatingModel{DB: db},
		},
	}

	return app, mock
}

func TestCreateCommentAndGetByID(t *testing.T) {
	app, mock := newTestApplication(t)

	// Mock database expectations for Insert
	comment := &model.Comment{
		UserID:      1,
		TeamID:      1,
		CommentText: "Test comment",
	}
	rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, "2024-06-06T00:00:00Z")

	mock.ExpectQuery(`INSERT INTO comments (user_id, team_id, comment_text) VALUES ($1, $2, $3) RETURNING id, created_at`).
		WithArgs(comment.UserID, comment.TeamID, comment.CommentText).
		WillReturnRows(rows)

	// Mock database expectations for GetByID
	mock.ExpectQuery(`SELECT id, user_id, team_id, comment_text, created_at FROM comments WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "team_id", "comment_text", "created_at"}).
			AddRow(1, 1, 1, "Test comment", "2024-06-06T00:00:00Z"))

	// Create comment
	body := `{"comment_text": "Test comment"}`
	req := httptest.NewRequest(http.MethodPost, "/comments/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}))
	rr := httptest.NewRecorder()
	app.createCommentHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", rr.Code)
		t.Errorf("Response body: %s", rr.Body.String())
	}

	// Fetch comment by ID
	req = httptest.NewRequest(http.MethodGet, "/comments/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}))
	rr = httptest.NewRecorder()
	app.getCommentByIDHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
		t.Errorf("Response body: %s", rr.Body.String())
	}

	var fetchedComment model.Comment
	err := json.NewDecoder(rr.Body).Decode(&fetchedComment)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
	if fetchedComment.ID != 1 {
		t.Errorf("Expected comment ID 1, got %d", fetchedComment.ID)
	}
}
