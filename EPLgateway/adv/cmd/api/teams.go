package main

import (
	"adv.erakaisar.net/internal/data"
	"adv.erakaisar.net/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createTeamsHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Stadium  string `json:"stadium"`
		History  string `json:"history"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the input struct to a new Movie struct.
	team := &data.Team{
		Name:     input.Name,
		Location: input.Location,
		Stadium:  input.Stadium,
		History:  input.History,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateTeam(v, team); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Teams.Insert(team)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/teams/%d", team.ID))

	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"team": team}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showTeamsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}
	team, err := app.models.Teams.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"team": team}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTeamsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	team, err := app.models.Teams.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Stadium  string `json:"stadium"`
		History  string `json:"history"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the movie
	// record.
	team.Name = input.Name
	team.Location = input.Location
	team.Stadium = input.Stadium
	team.History = input.History
	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateTeam(v, team); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated movie record to our new Update() method.
	err = app.models.Teams.Update(team)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"team": team}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTeamsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Teams.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "team successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//func (app *application) listTeamsHandler(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Name string
//		data.Filters
//	}
//	v := validator.New()
//	qs := r.URL.Query()
//	input.Name = app.readString(qs, "name", "")
//	input.Filters.Page = app.readInt(qs, "page", 1, v)
//	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
//	input.Filters.Sort = app.readString(qs, "sort", "id")
//	input.Filters.SortSafelist = []string{"id", "name", "location", "stadium", "history", "-id", "-name", "-location", "-stadium", "-history"}
//	if data.ValidateFilters(v, input.Filters); !v.Valid() {
//		app.failedValidationResponse(w, r, v.Errors)
//		return
//	}
//	// Call the GetAll() method to retrieve the movies, passing in the various filter
//	// parameters.
//	teams, err := app.models.Teams.GetAll(input.Name, input.Filters)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//		return
//	}
//
//	// Send a JSON response containing the movie data.
//	err = app.writeJSON(w, http.StatusOK, envelope{"teams": teams}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//
//}
