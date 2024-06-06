package model

import (
	"database/sql"
)

type Rating struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	TeamID    int    `json:"team_id"`
	Rating    int    `json:"rating"`
	CreatedAt string `json:"created_at"`
}

type RatingModel struct {
	DB *sql.DB
}

func (m *RatingModel) Insert(rating *Rating) error {
	query := `INSERT INTO ratings (user_id, team_id, rating) VALUES ($1, $2, $3) RETURNING id, created_at`
	return m.DB.QueryRow(query, rating.UserID, rating.TeamID, rating.Rating).Scan(&rating.ID, &rating.CreatedAt)
}

func (m *RatingModel) GetAll(teamID int) ([]*Rating, error) {
	query := `SELECT id, user_id, team_id, rating, created_at FROM ratings WHERE team_id = $1`
	rows, err := m.DB.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*Rating
	for rows.Next() {
		var rating Rating
		err := rows.Scan(&rating.ID, &rating.UserID, &rating.TeamID, &rating.Rating, &rating.CreatedAt)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	return ratings, nil
}
func (m *RatingModel) Update(rating *Rating) error {
	query := `UPDATE ratings SET rating = $1 WHERE id = $2`
	_, err := m.DB.Exec(query, rating.Rating, rating.ID)
	return err
}

func (m *RatingModel) Delete(id int) error {
	query := `DELETE FROM ratings WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *RatingModel) GetByID(id int) (*Rating, error) {
	query := `SELECT id, user_id, team_id, rating, created_at FROM ratings WHERE id = $1`
	row := m.DB.QueryRow(query, id)
	var rating Rating
	err := row.Scan(&rating.ID, &rating.UserID, &rating.TeamID, &rating.Rating, &rating.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (m *RatingModel) CountByTeamID(teamID int) (int, error) {
	query := `SELECT COUNT(*) FROM ratings WHERE team_id = $1`
	row := m.DB.QueryRow(query, teamID)
	var count int
	err := row.Scan(&count)
	return count, err
}
