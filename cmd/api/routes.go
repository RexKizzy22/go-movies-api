package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodPost, "/v1/graphql/list", app.moviesGraphQL)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getAllMoviesByGenre)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)
	router.POST("/v1/admin/edit-movie", app.chainMW(secure.ThenFunc(app.editMovie)))
	router.DELETE("/v1/admin/delete-movie/:id", app.chainMW(secure.ThenFunc(app.deleteMovie)))

	return app.enableCORS(router)
}
