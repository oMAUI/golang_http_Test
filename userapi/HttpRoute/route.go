package HttpRoute

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"refactoring/ErrorProcessing"
	"refactoring/Models/User"
	"strconv"
	"time"
)

const Data = `users.json`

func NewRoute() *chi.Mux {
	route := chi.NewRouter()

	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(middleware.Timeout(60 * time.Second))

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	route.Route("/api/v1/users", func(route chi.Router) {
		route.Get("/", getAllUsers)
		route.Post("/", createUser)

		route.Route("/{id}", func(route chi.Router) {
			route.Get("/", getUser)
			route.Patch("/", updateUser)
			route.Delete("/", deleteUser)
		})
	})

	return route
}

func UnmarshalBody(r io.Reader, v interface{}) error {
	resp, errResp := ioutil.ReadAll(r)
	if errResp != nil {

		return fmt.Errorf("server error: %w", errResp)
	}

	if errUnmarshalJson := json.Unmarshal(resp, v); errUnmarshalJson != nil {
		return fmt.Errorf("server error: %w", errUnmarshalJson)
	}

	return nil
}

func WriteFile (w http.ResponseWriter, v interface{}, DataPath string) bool {
	usersJsonByte, errMarshal := json.Marshal(&v)
	if errMarshal != nil {
		ErrorProcessing.HttpError(w, errMarshal, "failed to Marshal users",
			"server error", http.StatusInternalServerError)

		return false
	}
	if errWriteFile := ioutil.WriteFile(DataPath, usersJsonByte, fs.ModePerm); errWriteFile != nil {
		ErrorProcessing.HttpError(w, errWriteFile, "failed to write file",
			"server error", http.StatusInternalServerError)

		return false
	}

	return true
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, errGetUsersInData := User.GetUsersInData(Data)
	if errGetUsersInData != nil {
		ErrorProcessing.HttpError(w, errGetUsersInData, "failed to get users in data",
			"server error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	users, errGetUsersInData := User.GetUsersInData(Data)
	if errGetUsersInData != nil {
		ErrorProcessing.HttpError(w, errGetUsersInData, "failed to get users in data",
			"server error", http.StatusInternalServerError)
		return
	}

	userFromBody := User.User{}

	if errUnmarshalBody := UnmarshalBody(r.Body, &userFromBody); errUnmarshalBody != nil {
		ErrorProcessing.HttpError(w, errUnmarshalBody, "failed to unmarshal body",
			"server error", http.StatusInternalServerError)

		return
	}

	user := User.User{
		CreatedAt:   time.Now(),
		DisplayName: userFromBody.DisplayName,
		Email:       userFromBody.Email,
	}

	users.Increment++
	users.List[strconv.Itoa(users.Increment)] = user

	if WriteFile(w, users, Data) == false {
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": users.Increment,
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	users, errGetUsersInData := User.GetUsersInData(Data)
	if errGetUsersInData != nil {
		ErrorProcessing.HttpError(w, errGetUsersInData, "failed to get users in data",
			"server error", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := users.List[id]; !ok {
		ErrorProcessing.HttpError(w, errGetUsersInData, "user not found",
			"Not Found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, users.List[id])
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	users, errGetUsersInData := User.GetUsersInData(Data)
	if errGetUsersInData != nil {
		ErrorProcessing.HttpError(w, errGetUsersInData, "failed to get users in data",
			"server error", http.StatusInternalServerError)
		return
	}

	userFromBody := User.User{}

	if errUnmarshalBody := UnmarshalBody(r.Body, &userFromBody); errUnmarshalBody != nil {
		ErrorProcessing.HttpError(w, errUnmarshalBody, "failed to get rendering",
			"server error", http.StatusInternalServerError)

		return
	}

	id := chi.URLParam(r, "id")
	updatedUser, ok := users.List[id]
	if !ok {
		ErrorProcessing.HttpError(w, fmt.Errorf(""), "user not found",
			"Not Found", http.StatusNotFound)
		return
	}
	User.UpdateUser(userFromBody, &updatedUser)
	users.List[id] = updatedUser

	if WriteFile(w, users, Data) == false {
		return
	}

	render.Status(r, http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	users, errGetUsersInData := User.GetUsersInData(Data)
	if errGetUsersInData != nil {
		ErrorProcessing.HttpError(w, errGetUsersInData, "failed to get users in data",
			"server error", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := users.List[id]; !ok {
		ErrorProcessing.HttpError(w, fmt.Errorf(""), "user not found",
			"Not Found", http.StatusNotFound)
		return
	}

	delete(users.List, id)

	if WriteFile(w, users, Data) == false {
		return
	}
	render.Status(r, http.StatusNoContent)
}
