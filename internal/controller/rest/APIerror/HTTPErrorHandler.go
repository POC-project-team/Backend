package APIerror

import (
	"backend/internal/controller/rest/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HTTPErrorHandler struct {
	ErrorCode   int
	Description string
}

func HTTPErrorHandle(w http.ResponseWriter, err HTTPErrorHandler) {
	w.WriteHeader(err.ErrorCode)
	// If the Error is on server, then log it
	if err.ErrorCode == http.StatusInternalServerError {
		log.Error(err.Description)
	}
	_, err1 := w.Write([]byte(err.Description))
	if err1 != nil {
		return
	}
	return
}

func Error(w http.ResponseWriter, err error) {
	var HTTPError HTTPErrorHandler

	log.Error(err)

	switch err.Error() {
	case response.UserAlreadyExists:
		HTTPError = HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		}
	case response.UserNotFound:
		HTTPError = HTTPErrorHandler{
			ErrorCode:   http.StatusNotFound,
			Description: err.Error(),
		}
	case response.UserNotAuthorized:
		HTTPError = HTTPErrorHandler{
			ErrorCode:   http.StatusUnauthorized,
			Description: err.Error(),
		}
	case response.WrongLoginOrPassword:
		HTTPError = HTTPErrorHandler{
			ErrorCode:   http.StatusBadRequest,
			Description: err.Error(),
		}
	case response.NoSuchTag:
		HTTPError = HTTPErrorHandler{
			ErrorCode:   http.StatusNotFound,
			Description: err.Error(),
		}
	default:
		HTTPError = HTTPErrorHandler{
			ErrorCode:   http.StatusInternalServerError,
			Description: err.Error(),
		}
	}
	w.WriteHeader(HTTPError.ErrorCode)
	if _, err = w.Write([]byte(HTTPError.Description)); err != nil {
		log.Error(err.Error())
	}
}
