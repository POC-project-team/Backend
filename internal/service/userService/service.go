package userService

import (
	"backend/internal/controller/rest/APIerror"
	"backend/internal/dto/request"
	"backend/internal/dto/request/noteRequest"
	"backend/internal/dto/request/tagRequest"
	"backend/internal/dto/request/userRequest"
	"backend/internal/dto/responseDto"
	"backend/internal/entity"
	"backend/internal/repository"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	db repository.IClient
}

func NewService(database repository.IClient) *Service {
	return &Service{
		db: database,
	}
}

// GetAllUsers func to return all users in the map
func (s *Service) GetAllUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		APIerror.Error(w, err)
		return
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

// CreateUser handler for creating new user
func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest.AuthRequest

	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	result, err := s.db.CreateUser(req.Login, req.Password)
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if json.NewEncoder(w).Encode(result) != nil {
		APIerror.Error(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

// GetAllUsersTags handler for getting all tags of specific user
func (s *Service) GetAllUsersTags(w http.ResponseWriter, r *http.Request) {
	var (
		tags []entity.Tag
	)

	token, err := request.ParseToken(r)
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if tags, err = s.db.GetUserTags(token.UserId); err != nil {
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

	if err = json.NewEncoder(w).Encode(tags); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Service) GetTag(w http.ResponseWriter, r *http.Request) {
	var (
		tag entity.Tag
		req request.BasicRequest
		err error
	)

	if err := req.BindBasicRequest(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if tag, err = s.db.GetTag(req.Token.UserId, req.TagId); err != nil {
		APIerror.Error(w, err)
		return
	}

	resp := responseDto.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var (
		tag entity.Tag
		req tagRequest.CreateUpdateTagRequest
		err error
	)

	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if tag, err = s.db.UpdateTag(req.Token.UserId, req.TagId, req.TagName); err != nil {
		APIerror.Error(w, err)
		return
	}

	resp := responseDto.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) CreateTag(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Tag responseDto.TagNoUserNotes `json:"tag"`
	}
	var (
		resp response
		tag  entity.Tag
		req  tagRequest.CreateUpdateTagRequest
		err  error
	)

	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if tag, err = s.db.CreateTag(req.Token.UserId, req.TagId, req.TagName); err != nil {
		APIerror.Error(w, err)
		return
	}

	resp.Tag = responseDto.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) DeleteTag(w http.ResponseWriter, r *http.Request) {
	var (
		req request.BasicRequest
	)
	if err := req.BindBasicRequest(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if err := s.db.DeleteTag(req.Token.UserId, req.TagId); err != nil {
		APIerror.Error(w, err)
		return
	}

	resp := "Tag deleted"

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) TransferTag(w http.ResponseWriter, r *http.Request) {
	var req tagRequest.TransferTagRequest
	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	if err := s.db.TransferTag(req.Token.UserId, req.TagId, req.Login); err != nil {
		APIerror.Error(w, err)
		return
	}

	resp := "Tag was transferred"

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// GetNotes handler for getting notes for specific tag of user
func (s *Service) GetNotes(w http.ResponseWriter, r *http.Request) {
	var req request.BasicRequest

	if err := req.BindBasicRequest(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	notes, err := s.db.GetUserNotes(req.Token.UserId, req.TagId)
	if err != nil {
		APIerror.Error(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(notes); err != nil {
		APIerror.Error(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// AddNote handler for creating new note for specific tag of user
func (s *Service) AddNote(w http.ResponseWriter, r *http.Request) {
	var req noteRequest.CreateNoteRequest

	if err := req.Bind(r); err != nil {
		APIerror.Error(w, err)
		return
	}

	response, err := s.db.AddNote(req.Token.UserId, req.TagId, req.Note)
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
