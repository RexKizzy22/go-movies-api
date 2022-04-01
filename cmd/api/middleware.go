package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pascaldekloe/jwt"
)

// YET TO UNDERSTAND!!!
func (app *application) chainMW(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authentication")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// could set anonymous user here
		}

		bearerAuth := strings.Split(authHeader, "")
		if len(bearerAuth) != 2 {
			app.errorJSON(w, errors.New("Invalid authorization"))
			return
		}

		if bearerAuth[0] != "Bearer" {
			app.errorJSON(w, errors.New("Invalid authorization - Bearer not found"))
		}

		claims, err := jwt.HMACCheck([]byte(bearerAuth[1]), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("Invalid authorization - Failed HMAC check"), http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("Invalid authorization - Token expired"), http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			app.errorJSON(w, errors.New("Invalid authorization - Invalid audience"), http.StatusForbidden)
			return
		}

		if claims.Issuer != "mydomain.com" {
			app.errorJSON(w, errors.New("Invalid authorization - Invalid user"), http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("Unable to parse user id"), http.StatusForbidden)
			return
		}

		log.Println("userid:", userID)

		next.ServeHTTP(w, r)
	})
}
