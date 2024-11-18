package broker

import (
	"errors"
	"net/http"

	"github.com/pivotal-cf/brokerapi/v11/domain/apiresponses"
)

var (
	ErrNotFound         = apiresponses.NewFailureResponse(errors.New("not found"), http.StatusNotFound, "not-found")
	ErrConcurrencyError = apiresponses.NewFailureResponse(errors.New("ConcurrencyError"), http.StatusUnprocessableEntity, "concurrency-error")
	ErrBadRequest       = apiresponses.NewFailureResponse(errors.New("bad request"), http.StatusBadRequest, "bad-request")
)
