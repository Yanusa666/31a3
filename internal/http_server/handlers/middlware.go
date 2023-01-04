package handlers

import (
	"context"
	"crud/internal/constants"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) Middlware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		requestId := uuid.New().String()
		ctx = context.WithValue(ctx, constants.RequestIdKey, requestId)
		r = r.WithContext(ctx)

		w.Header().Add(constants.RequestIdKey, requestId)
		w.Header().Add("Content-Type", "application/json")

		handle(w, r, ps)
	}
}
