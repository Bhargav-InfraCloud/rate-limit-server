package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bhargav-InfraCloud/rate-limit-server/pkg/service/errors"
)

func WriteResponse(rw http.ResponseWriter, response interface{}) error {
	rw.WriteHeader(http.StatusOK)

	err := json.NewEncoder(rw).Encode(&response)
	if err != nil {
		return fmt.Errorf("failed to encode to response writer: %v", err)
	}

	return nil
}

func WriteErrorResponse(rw http.ResponseWriter, serviceError errors.Error) error {
	rw.WriteHeader(serviceError.StatusCode())

	err := json.NewEncoder(rw).Encode(&serviceError)
	if err != nil {
		return fmt.Errorf("failed to encode service error to response writer: %v", err)
	}

	return nil
}
