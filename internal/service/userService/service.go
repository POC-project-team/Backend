package userService

import (
	"backend/internal/controller/rest/APIerror"
	"backend/internal/controller/rest/response"
	"backend/internal/dto/request"
	"backend/internal/dto/request/userRequest"
	"backend/internal/dto/responseDto"
	u "backend/internal/entity"
	"backend/internal/repository/postgres"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	db *postgres.Client
}

func NewService(database *postgres.Client) *Service {
	return &Service{
		db: database,
		//*db.NewSQLDataBase(),
	}
}

// GetAllUsers func to return all users in the map
func (s *Service) GetAllUsers(w http.ResponseWriter, _ *http.Request) {
	type response struct {
		Users []string `json:"Users"`
	}

	var (
		resp response
		err  error
	)
	resp.Users, err = s.db.GetAllUsers()
	if err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: err.Error(),
		})
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// CreateUser handler for creating new user
func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest.AuthRequest
	if req.Bind(r) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := s.db.CreateUser(req.Login, req.Password)
	if err != nil {
		switch err.Error() {
		case response.UserAlreadyExists:
			APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
				ErrorCode:   http.StatusBadRequest,
				Description: err.Error(),
			})
		default:
			APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
				ErrorCode:   http.StatusInternalServerError,
				Description: err.Error(),
			})
		}
	}

	if json.NewEncoder(w).Encode(result) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// GetAllUsersTags handler for getting all tags of specific user
func (s *Service) GetAllUsersTags(w http.ResponseWriter, r *http.Request) {
	var (
		tags []responseDto.TagNoUserNotes
		req  request.Request
		err  error
	)
	if req.ParseToken(w, r) != nil {
		return
	}

	if tags, err = s.db.GetUserTags(req.UserID); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	if err = json.NewEncoder(w).Encode(tags); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Service) GetTag(w http.ResponseWriter, r *http.Request) {
	var (
		resp responseDto.TagNoUserNotes
		req  request.Request
		err  error
	)
	if req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	if req.TagID == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag id provided",
		})
		return
	}

	if resp, err = s.db.GetTag(req.UserID, req.TagID); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	if resp.TagName == "" && resp.TagID == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No such tag",
		})
		return
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Service) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var (
		req  request.Request
		resp responseDto.TagNoUserNotes
		err  error
	)
	if req.Bind(w, r) != nil || req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	if req.TagID == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag id provided",
		})
		return
	}

	if req.TagName == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag name provided",
		})
		return
	}

	if resp, err = s.db.UpdateTag(req.UserID, req.TagID, req.TagName); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Service) CreateTag(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Tag responseDto.TagNoUserNotes `json:"tag"`
	}
	var (
		resp response
		req  request.Request
		err  error
	)
	if req.Bind(w, r) != nil || req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	if req.TagID == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag id provided",
		})
		return
	}

	if req.TagName == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag name provided",
		})
		return
	}

	if resp.Tag, err = s.db.CreateTag(req.UserID, req.TagID, req.TagName); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Service) DeleteTag(w http.ResponseWriter, r *http.Request) {
	var (
		req request.Request
		err error
	)
	if req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	if req.TagID == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag id provided",
		})
		return
	}

	if err = s.db.DeleteTag(req.UserID, req.TagID); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	resp := "Tag deleted"

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Service) TransferTag(w http.ResponseWriter, r *http.Request) {
	var (
		req request.Request
		err error
	)
	if req.Bind(w, r) != nil || req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	if req.TagID == "" {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: "No tag id provided",
		})
		return
	}

	if err = s.db.TransferTag(req.UserID, req.TagID, req.Login); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	resp := "Tag was transferred"

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// GetNotes handler for getting notes for specific tag of user
func (s *Service) GetNotes(w http.ResponseWriter, r *http.Request) {
	var req request.Request
	if req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	var notes []u.Note

	notes, err := s.db.GetUserNotes(req.UserID, req.TagID)
	if err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	if err = json.NewEncoder(w).Encode(notes); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// AddNote handler for creating new note for specific tag of user
func (s *Service) AddNote(w http.ResponseWriter, r *http.Request) {
	var req request.Request
	// param checking
	if req.Bind(w, r) != nil || req.ParseToken(w, r) != nil || req.ParseTagID(w, r) != nil {
		return
	}

	response, err := s.db.AddNote(req.UserID, req.TagID, req.Note)
	if err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		APIerror.HTTPErrorHandle(w, APIerror.HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: "Cannot write data to request",
		})
	} else {
		log.Info("New note for userID: ", req.UserID, " tagID: ", req.TagID, " was created")
		w.WriteHeader(http.StatusCreated)
	}
}
