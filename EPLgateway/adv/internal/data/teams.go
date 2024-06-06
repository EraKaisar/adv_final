package data

import (
	"adv.erakaisar.net/internal/validator"
	"database/sql"
	"errors"
	"time"
)

type Team struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Stadium   string    `json:"stadium"`
	History   string    `json:"history"`
}

func ValidateTeam(v *validator.Validator, team *Team) {
	v.Check(team.Name != "", "name", "must be provided")
	v.Check(len(team.Name) <= 100, "name", "must not be more than 500 bytes long")

	v.Check(team.Location != "", "location", "must be provided")
	v.Check(len(team.Location) <= 100, "location", "must not be more than 500 bytes long")

	v.Check(team.Stadium != "", "stadium", "must be provided")
	v.Check(len(team.Stadium) <= 100, "stadium", "must not be more than 500 bytes long")

	v.Check(team.History != "", "history", "must be provided")
	v.Check(len(team.History) <= 1000, "history", "must not be more than 2000 bytes long")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type TeamModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m TeamModel) Insert(team *Team) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
        INSERT INTO teams (name, location, stadium, history) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at `
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{team.Name, team.Location, team.Stadium, team.History}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&team.ID, &team.CreatedAt)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m TeamModel) Get(id int64) (*Team, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
        SELECT id, created_at, name, location, stadium, history
        FROM teams
        WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var team Team
	err := m.DB.QueryRow(query, id).Scan(
		&team.ID,
		&team.CreatedAt,
		&team.Name,
		&team.Location,
		&team.Stadium,
		&team.History,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &team, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m TeamModel) Update(team *Team) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
        UPDATE teams 
        SET name = $1, location = $2, stadium = $3, history = $4
        WHERE id = $5
        RETURNING id`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		team.Name,
		team.Location,
		team.Stadium,
		team.History,
		team.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&team.ID)
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m TeamModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
        DELETE FROM teams
        WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

type MockTeamModel struct{}

func (m MockTeamModel) Insert(team *Team) error {
	// Mock the action...
	return nil
}
func (m MockTeamModel) Get(id int64) (*Team, error) {
	return nil, nil
}
func (m MockTeamModel) Update(team *Team) error {
	// Mock the action...
	return nil
}
func (m MockTeamModel) Delete(id int64) error {
	return nil
}

//func (m TeamModel) GetAll(name string, filters Filters) ([]*Team, error) {
//	// Update the SQL query to include the filter conditions.
//	query := `
//        SELECT id, created_at, name, location, stadium, history
//        FROM teams
//        WHERE (LOWER(name) = LOWER($1) OR $1 = '')
//        ORDER BY id`
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//	// Pass the title and genres as the placeholder parameter values.
//	rows, err := m.DB.QueryContext(ctx, query, name)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	teams := []*Team{}
//	for rows.Next() {
//		var team Team
//		err := rows.Scan(
//			&team.ID,
//			&team.CreatedAt,
//			&team.Name,
//			&team.Location,
//			&team.Stadium,
//			&team.History,
//		)
//		if err != nil {
//			return nil, err
//		}
//		teams = append(teams, &team)
//	}
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//	return teams, nil
//}
