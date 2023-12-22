package utils

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func AddValueInRequestContext(r *http.Request, key string, value interface{}) {
	*r = *r.Clone(context.WithValue(r.Context(), key, value))
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}
