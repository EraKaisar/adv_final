package model

import (
	"database/sql"
)

type Models struct {
	Comments CommentModel
	Ratings  RatingModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Comments: CommentModel{DB: db},
		Ratings:  RatingModel{DB: db},
	}
}
