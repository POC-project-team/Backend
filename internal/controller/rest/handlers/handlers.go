package handlers

import (
	"backend/internal/controller/rest/APIerror"
	db "backend/internal/repository/sqlite"
	au "backend/internal/service/auth"
	service "backend/internal/service/userService"
	"github.com/gorilla/mux"
	"net/http"
)

// MyHandler defines the routes, returns router
func MyHandler(database *db.SQL) *mux.Router {
	srv := service.NewService(database)
	router := mux.NewRouter()

	Routes(router, srv)
	router.Handle("/", router)
	return router
}

func Routes(router *mux.Router, srv *service.Service) {
	// routers for users and auth
	router.HandleFunc("/users", srv.GetAllUsers).Methods("GET")
	router.HandleFunc("/signup", srv.CreateUser).Methods("POST")
	router.HandleFunc("/auth", au.Auth).Methods("POST")

	// routers for users settings and profile
	router.HandleFunc("/{token}/changeLogin", au.ChangeLogin).Methods("POST")
	router.HandleFunc("/{token}/changePassword", au.ChangePassword).Methods("POST")

	// routers for tag handling
	router.HandleFunc("/{token}/tags", srv.GetAllUsersTags).Methods("GET")
	router.HandleFunc("/{token}/{tag_id}/tag", srv.CreateTag).Methods("POST")
	router.HandleFunc("/{token}/{tag_id}/tag", srv.GetTag).Methods("GET")
	router.HandleFunc("/{token}/{tag_id}/tag", srv.UpdateTag).Methods("PUT")
	router.HandleFunc("/{token}/{tag_id}/tag", srv.DeleteTag).Methods("DELETE")
	router.HandleFunc("/{token}/{tag_id}/send", srv.TransferTag).Methods("POST")

	// routers for notes handling
	router.HandleFunc("/{token}/{tag_id}/notes", srv.GetNotes).Methods("GET")
	router.HandleFunc("/{token}/{tag_id}/note", srv.AddNote).Methods("POST")

	router.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusCreated)
		if _, err := writer.Write([]byte("Hello, I'm working\n")); err != nil {
			APIerror.HTTPErrorHandle(writer, APIerror.HTTPErrorHandler{
				ErrorCode:   http.StatusBadRequest,
				Description: "I don't know…",
			})
		}
	}).Methods("GET")
}
