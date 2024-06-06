package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// Initialize a new httprouter router instance.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//router.HandlerFunc(http.MethodGet, "/v1/teams", app.listTeamsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/teams", app.createTeamsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/teams/:id", app.showTeamsHandler)
	router.HandlerFunc(http.MethodPut, "/v1/teams/:id", app.updateTeamsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/teams/:id", app.deleteTeamsHandler)
	return router
}
