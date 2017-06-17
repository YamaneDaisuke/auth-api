package authapi

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/charakoba-com/auth-api/db"
	"github.com/charakoba-com/auth-api/model"
	"github.com/charakoba-com/auth-api/service"
	"github.com/charakoba-com/auth-api/utils"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const defaultAlgorithm = "RS256"

var algorithm string

func init() {
	algorithm = defaultAlgorithm
}

// HealthCheckHandler is a HTTP handler, which path is `/`
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HealthCheckHandler")
	method := r.Method
	if method != `GET` {
		httpError(w, http.StatusMethodNotAllowed, `method GET is expected`, nil)
		return
	}
	httpJSON(w, map[string]string{
		"message": "hello, world",
	})
}

// CreateUserHandler is a HTTP handler, which creates an new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateUserHandler")

	// Verify Request
	method := r.Method
	if method != `POST` {
		httpError(w, http.StatusMethodNotAllowed, `method POST is expected`, nil)
		return
	}
	ctype := r.Header["Content-Type"][0]
	if ctype != `application/json` {
		httpError(w, http.StatusBadRequest, `Content-Type: application/json is expected`, nil)
		return
	}

	// preparation
	var createUserRequest model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		httpError(w, http.StatusBadRequest, `invalid json request`, err)
		return
	}
	newUser := db.User{
		ID:       createUserRequest.ID,
		Name:     createUserRequest.Username,
		Password: createUserRequest.Password,
	}

	// main logic
	tx, err := db.BeginTx()
	if err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	usrSvc := service.UserService{}
	if err := usrSvc.Create(tx, &newUser); err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}

	httpJSON(w, map[string]string{"message": "success"})
}

// LookupUserHandler is a HTTP handler, which search an user by ID
func LookupUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LookupUserHandler")

	// Verify Request
	method := r.Method
	if method != `GET` {
		httpError(w, http.StatusMethodNotAllowed, `method GET is expected`, nil)
		return
	}
	id := mux.Vars(r)["id"]
	tx, err := db.BeginTx()
	if err != nil {
		httpError(w, http.StatusInternalServerError, `database error`, err)
		return
	}
	var usrSvc service.UserService
	user, err := usrSvc.Lookup(tx, id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			httpError(w, http.StatusNotFound, `user not found`, err)
			return
		}
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	user.Password = ""
	httpJSON(w, model.LookupUserResponse{User: *user})
}

// UpdateUserHandler is a HTTP handler, which updates an user
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateUserHandler")

	// Verify Request
	method := r.Method
	if method != `PUT` {
		httpError(w, http.StatusMethodNotAllowed, `method PUT is expected`, nil)
		return
	}
	ctype := r.Header["Content-Type"][0]
	if ctype != `application/json` {
		httpError(w, http.StatusBadRequest, `Content-Type: application/json is expected`, nil)
		return
	}

	// preparation
	var updateUserRequest model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		httpError(w, http.StatusBadRequest, `invalid json request`, err)
		return
	}
	updater := db.User{
		ID:       updateUserRequest.ID,
		Name:     updateUserRequest.Username,
		Password: updateUserRequest.Password,
	}

	// main logic
	tx, err := db.BeginTx()
	if err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	usrSvc := service.UserService{}
	if err := usrSvc.Update(tx, &updater); err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	httpJSON(w, map[string]string{"message": "success"})
}

// DeleteUserHandler is a HTTP handler, which deletes an user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteUserHandler")
	method := r.Method
	if method != `DELETE` {
		httpError(w, http.StatusMethodNotAllowed, `method DELETE is expected`, nil)
		return
	}
	id := mux.Vars(r)["id"]
	tx, err := db.BeginTx()
	if err != nil {
		httpError(w, http.StatusInternalServerError, `database error`, err)
		return
	}
	var usrSvc service.UserService
	if err := usrSvc.Delete(tx, id); err != nil {
		httpError(w, http.StatusInternalServerError, `deleting user`, err)
		return
	}
	if err := tx.Commit(); err != nil {
		httpError(w, http.StatusInternalServerError, `deleting user(commit)`, err)
		return
	}
	httpJSON(w, map[string]string{"message": "success"})
}

// ListupUserHandler is a HTTP handler, which returns all user list
func ListupUserHandler(w http.ResponseWriter, r *http.Request) {
	// NotImplemented
	log.Printf("ListupUserHandler")
	method := r.Method
	if method != `GET` {
		httpError(w, http.StatusMethodNotAllowed, `method GET is expected`, nil)
		return
	}
	tx, err := db.BeginTx()
	if err != nil {
		httpError(w, http.StatusInternalServerError, `database error`, err)
		return
	}
	var usrSvc service.UserService
	users ,err := usrSvc.Listup(tx)
	if err != nil {
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	httpJSON(w, model.ListupUserResponse{Users: users})
}

// AuthHandler is a HTTP handler, which authes with username and password
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// NotImplemented
	log.Printf("AuthHandler")

	method := r.Method
	if method != `POST` {
		httpError(w, http.StatusMethodNotAllowed, `method POST is expected`, nil)
		return
	}
	var authRequest model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		httpError(w, http.StatusBadRequest, `invalid json request`, nil)
		return
	}
	tx, err := db.BeginTx()
	if err != nil {
		httpError(w, http.StatusInternalServerError, `database errorr`, err)
		return
	}
	var usrSvc service.UserService
	user, err := usrSvc.Lookup(tx, authRequest.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			httpError(w, http.StatusUnauthorized, `auth invalid`, nil)
			return
		}
		httpError(w, http.StatusInternalServerError, `internal server error`, err)
		return
	}
	if user.Password != utils.HashPassword(authRequest.Password, authRequest.ID + user.Name) {
		httpError(w, http.StatusUnauthorized, `auth invalid`, nil)
		return
	}

	httpJSON(w, model.AuthResponse{
		Message: "auth valid",
		Token:   "",
	})
}

// GetAlgorithmHandler is a HTTP handler, which returns system signature algorithm
func GetAlgorithmHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAlgorithmHandler")
	method := r.Method
	if method != `GET` {
		httpError(w, http.StatusMethodNotAllowed, `method GET is expected`, nil)
		return
	}
	httpJSON(w, model.GetAlgorithmResponse{Algorithm: algorithm})
}

// VerifyHandler is a HTTP handler, which verifies given token
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// NotImplemented
	log.Printf("VerifyHandler")
}

// GetKeyHandler is a HTTP handler, which returns public key verifying token
func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	// NotImplemented
	log.Printf("GetKeyHandler")
}

// NotFoundHandler is a HTTP handler, which handles 404 Not Found
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("NotFoundHandler")
	httpError(w, http.StatusNotFound, `not found`, nil)
}
