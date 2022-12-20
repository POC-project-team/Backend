package handlers

import (
	"backend/internal/controller/rest/APIerror"
	"backend/internal/dto/request"
	"backend/internal/dto/request/noteRequest"
	"backend/internal/dto/request/tagRequest"
	"backend/internal/dto/request/userRequest"
	"backend/internal/dto/responseDto"
	"backend/internal/entity"
	"backend/internal/repository"
	au "backend/internal/service/auth"
	service "backend/internal/service/userService"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	Router      *mux.Router
	Service     *service.Service
	AuthService *au.Service
}

func NewHandler(database repository.IClient) *Handler {
	router := mux.NewRouter()
	return &Handler{
		Router:      router,
		Service:     service.NewService(database),
		AuthService: au.NewAuthService(database),
	}
}

func (h *Handler) Routes() {
	// routers for users and auth
	h.Router.HandleFunc("/users", h.GetAllUsers).Methods("GET")
	h.Router.HandleFunc("/signup", h.CreateUser).Methods("POST")
	// routers for users settings and profile
	h.Router.HandleFunc("/auth", h.Authorise).Methods("POST")
	h.Router.HandleFunc("/{token}/changeLogin", h.ChangeLogin).Methods("POST")
	h.Router.HandleFunc("/{token}/changePassword", h.ChangePassword).Methods("POST")
	// routers for tag handling
	h.Router.HandleFunc("/{token}/tags", h.GetAllUsersTags).Methods("GET")
	h.Router.HandleFunc("/{token}/{tag_id}/tag", h.GetTag).Methods("GET")
	h.Router.HandleFunc("/{token}/{tag_id}/tag", h.CreateTag).Methods("POST")
	h.Router.HandleFunc("/{token}/{tag_id}/tag", h.UpdateTag).Methods("PUT")
	h.Router.HandleFunc("/{token}/{tag_id}/tag", h.DeleteTag).Methods("DELETE")
	h.Router.HandleFunc("/{token}/{tag_id}/transfer", h.TransferTag).Methods("POST")
	// routers for notes handling
	h.Router.HandleFunc("/{token}/{tag_id}/notes", h.GetNotes).Methods("GET")
	h.Router.HandleFunc("/{token}/{tag_id}/note", h.AddNote).Methods("POST")

	// test router
	h.Router.HandleFunc("/test", h.Test)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		APIerror.Error(w, err)
	}
	// return all id's of users
	var usersId []string
	for _, user := range users {
		usersId = append(usersId, fmt.Sprint(user.UserID))
	}

	resp := responseDto.UsersResponse{
		Users: usersId,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest.AuthRequest
	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	user, err := h.Service.CreateUser(entity.User{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if json.NewEncoder(w).Encode(user) != nil {
		APIerror.Error(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Authorise(w http.ResponseWriter, r *http.Request) {
	var response userRequest.AuthRequest

	if err := response.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	token, err := h.AuthService.Auth(entity.User{
		Login:    response.Login,
		Password: response.Password,
	})
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(token); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) ChangeLogin(w http.ResponseWriter, r *http.Request) {
	var response userRequest.ChangeLoginRequest

	if err := response.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if err := h.AuthService.ChangeLogin(entity.User{
		Login:  response.Login,
		UserID: response.Token.UserId,
	}); err != nil {
		APIerror.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("login changed"))
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var response userRequest.ChangePasswdRequest

	if response.Bind(r) != nil {
		return
	}
	if err := h.AuthService.ChangePassword(entity.User{
		Password: response.Password,
		UserID:   response.UserID,
	}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("password changed"))
}

func (h *Handler) GetAllUsersTags(w http.ResponseWriter, r *http.Request) {
	token, err := request.ParseToken(r)
	if err != nil {
		APIerror.Error(w, err)
		return
	}
	tags, err := h.Service.GetAllUsersTags(token.UserId)
	if err != nil {
		APIerror.Error(w, err)
		return
	}
	var tagsNoUserNotes []responseDto.TagNoUserNotes
	for _, tag := range tags {
		tagsNoUserNotes = append(tagsNoUserNotes, responseDto.TagNoUserNotes{
			TagID:   tag.TagID,
			TagName: tag.TagName,
		})
	}

	if err = json.NewEncoder(w).Encode(tagsNoUserNotes); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	var req request.BasicRequest
	if err := req.BindBasicRequest(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	tag, err := h.Service.GetTag(req.Token.UserId, req.TagId)
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	resp := responseDto.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req tagRequest.CreateUpdateTagRequest
	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	tag, err := h.Service.CreateTag(entity.Tag{
		TagName: req.TagName,
		UserID:  req.Token.UserId,
		TagID:   req.TagId,
	})
	if err != nil {
		APIerror.Error(w, err)
		return
	}
	type response struct {
		Tag responseDto.TagNoUserNotes `json:"tag"`
	}
	resp := response{
		Tag: responseDto.TagNoUserNotes{
			TagID:   tag.TagID,
			TagName: tag.TagName,
		},
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var req tagRequest.CreateUpdateTagRequest
	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	tag, err := h.Service.UpdateTag(entity.Tag{
		UserID:  req.Token.UserId,
		TagID:   req.TagId,
		TagName: req.TagName,
	})
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	resp := responseDto.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	var req request.BasicRequest
	if err := req.BindBasicRequest(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	if err := h.Service.DeleteTag(entity.Tag{
		UserID: req.Token.UserId,
		TagID:  req.TagId,
	}); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("tag deleted"))
}

func (h *Handler) TransferTag(w http.ResponseWriter, r *http.Request) {
	var req tagRequest.TransferTagRequest
	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	if err := h.Service.TransferTag(entity.Tag{
		UserID: req.Token.UserId,
		TagID:  req.TagId,
	}, req.Login); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("tag transferred"))
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	var req request.BasicRequest

	if err := req.BindBasicRequest(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	notes, err := h.Service.GetNotes(req.Token.UserId, req.TagId)
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(notes); err != nil {
		APIerror.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add note")
	var req noteRequest.CreateNoteRequest
	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}
	response, err := h.Service.AddNote(entity.Note{
		UserId: req.Token.UserId,
		TagID:  req.TagId,
		Note:   req.Note,
	})
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		APIerror.Error(w, err)
	} else {
		log.Info("New note for userID: ", req.Token.UserId, " tagID: ", req.TagId, " was created")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Test(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte("Hello, I'm working\n")); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "I don't know…",
		})
	}
}
