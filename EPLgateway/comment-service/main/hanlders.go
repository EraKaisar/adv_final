package main

import (
	"EPLgateway/comment-service/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	teamID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	var input struct {
		CommentText string `json:"comment_text"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	comment := &model.Comment{
		UserID:      app.contextGetUserID(r),
		TeamID:      teamID,
		CommentText: input.CommentText,
	}

	err = app.models.Comments.Insert(comment)
	if err != nil {
		http.Error(w, "Unable to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getCommentByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := app.models.Comments.GetByID(id)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

func (app *application) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	teamID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	comments, err := app.models.Comments.GetAll(teamID)
	if err != nil {
		http.Error(w, "Unable to fetch comments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (app *application) createRatingHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	teamID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Rating int `json:"rating"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	rating := &model.Rating{
		UserID: app.contextGetUserID(r),
		TeamID: teamID,
		Rating: input.Rating,
	}

	err = app.models.Ratings.Insert(rating)
	if err != nil {
		http.Error(w, "Unable to add rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) listRatingsHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	teamID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	ratings, err := app.models.Ratings.GetAll(teamID)
	if err != nil {
		http.Error(w, "Unable to fetch ratings", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ratings)
}
func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	commentID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var input struct {
		CommentText string `json:"comment_text"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	comment := &model.Comment{
		ID:          commentID,
		CommentText: input.CommentText,
	}

	err = app.models.Comments.Update(comment)
	if err != nil {
		http.Error(w, "Unable to update comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	commentID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = app.models.Comments.Delete(commentID)
	if err != nil {
		http.Error(w, "Unable to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) updateRatingHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	ratingID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid rating ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Rating int `json:"rating"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	rating := &model.Rating{
		ID:     ratingID,
		Rating: input.Rating,
	}

	err = app.models.Ratings.Update(rating)
	if err != nil {
		http.Error(w, "Unable to update rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) deleteRatingHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	ratingID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid rating ID", http.StatusBadRequest)
		return
	}

	err = app.models.Ratings.Delete(ratingID)
	if err != nil {
		http.Error(w, "Unable to delete rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) countCommentsByTeamIDHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	teamID, err := strconv.Atoi(params.ByName("team_id"))
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	count, err := app.models.Comments.CountByTeamID(teamID)
	if err != nil {
		http.Error(w, "Unable to count comments", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func (app *application) getRatingByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid rating ID", http.StatusBadRequest)
		return
	}

	rating, err := app.models.Ratings.GetByID(id)
	if err != nil {
		http.Error(w, "Unable to fetch rating", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rating)
}

func (app *application) countRatingsByTeamIDHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	teamID, err := strconv.Atoi(params.ByName("team_id"))
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	count, err := app.models.Ratings.CountByTeamID(teamID)
	if err != nil {
		http.Error(w, "Unable to count ratings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
