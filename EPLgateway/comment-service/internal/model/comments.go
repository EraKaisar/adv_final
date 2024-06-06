package model

import (
	"database/sql"
)

type Comment struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	TeamID      int    `json:"team_id"`
	CommentText string `json:"comment_text"`
	CreatedAt   string `json:"created_at"`
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(comment *Comment) error {
	query := `INSERT INTO comments (user_id, team_id, comment_text) VALUES ($1, $2, $3) RETURNING id, created_at`
	return m.DB.QueryRow(query, comment.UserID, comment.TeamID, comment.CommentText).Scan(&comment.ID, &comment.CreatedAt)
}

func (m *CommentModel) GetAll(teamID int) ([]*Comment, error) {
	query := `SELECT id, user_id, team_id, comment_text, created_at FROM comments WHERE team_id = $1`
	rows, err := m.DB.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.TeamID, &comment.CommentText, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}
func (m *CommentModel) Update(comment *Comment) error {
	query := `UPDATE comments SET comment_text = $1 WHERE id = $2`
	_, err := m.DB.Exec(query, comment.CommentText, comment.ID)
	return err
}

func (m *CommentModel) Delete(id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
func (m *CommentModel) GetByID(id int) (*Comment, error) {
	query := `SELECT id, user_id, team_id, comment_text, created_at FROM comments WHERE id = $1`
	row := m.DB.QueryRow(query, id)
	var comment Comment
	err := row.Scan(&comment.ID, &comment.UserID, &comment.TeamID, &comment.CommentText, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (m *CommentModel) CountByTeamID(teamID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE team_id = $1`
	row := m.DB.QueryRow(query, teamID)
	var count int
	err := row.Scan(&count)
	return count, err
}
