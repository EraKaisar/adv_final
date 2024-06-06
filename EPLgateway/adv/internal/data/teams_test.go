package data_test

import (
	"adv.erakaisar.net/internal/data"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func TestTeamModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new instance of TeamModel with the mock DB.
	model := data.TeamModel{DB: db}

	// Expectations for the mock DB.
	mock.ExpectExec("^INSERT INTO teams").WithArgs("Test Team", "Location", "Stadium", "History").WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new team.
	team := &data.Team{
		Name:     "Test Team",
		Location: "Location",
		Stadium:  "Stadium",
		History:  "History",
	}

	// Insert the team.
	err = model.Insert(team)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestTeamModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new instance of TeamModel with the mock DB.
	model := data.TeamModel{DB: db}

	// Expectations for the mock DB.
	rows := sqlmock.NewRows([]string{"id", "created_at", "name", "location", "stadium", "history"}).
		AddRow(1, time.Now(), "Test Team", "Location", "Stadium", "History")
	mock.ExpectQuery("^SELECT").WithArgs(1).WillReturnRows(rows)

	// Get the team.
	_, err = model.Get(1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestTeamModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new instance of TeamModel with the mock DB.
	model := data.TeamModel{DB: db}

	// Expectations for the mock DB.
	mock.ExpectExec("^UPDATE teams").WithArgs("Updated Team", "Location", "Stadium", "History", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new team.
	team := &data.Team{
		ID:       1,
		Name:     "Updated Team",
		Location: "Location",
		Stadium:  "Stadium",
		History:  "History",
	}

	// Update the team.
	err = model.Update(team)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestTeamModel_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new instance of TeamModel with the mock DB.
	model := data.TeamModel{DB: db}

	// Expectations for the mock DB.
	mock.ExpectExec("^DELETE FROM teams").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	// Delete the team.
	err = model.Delete(1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
