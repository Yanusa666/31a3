package handlers

import (
	"crud/internal/config"
	"crud/internal/constants"
	"crud/internal/entities"
	"crud/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

type Handler struct {
	cfg     *config.Config
	lgr     zerolog.Logger
	authors storage.IAuthors
	posts   storage.IPosts
}

func NewHandler(cfg *config.Config, lgr zerolog.Logger, stor *storage.Storage) *Handler {
	return &Handler{
		cfg:     cfg,
		lgr:     lgr,
		authors: stor.Authors,
		posts:   stor.Posts,
	}
}

func (h *Handler) AddAuthor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(AddAuthorReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "AddAuthor").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("name", request.Name)).
		Logger()

	err = h.authors.Add(ctx, &entities.Author{Name: request.Name})
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListAuthors(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "ListAuthors").
		Str(constants.RequestIdKey, requestId).
		Logger()

	listAuthors, err := h.authors.List(ctx)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")

	resp, _ := json.Marshal(ListAuthorsResp{Authors: listAuthors})
	fmt.Fprintf(w, string(resp))
}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(UpdateAuthorReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	idStr := ps.ByName("id")
	lgr := h.lgr.With().
		Str("handler", "UpdateAuthor").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("id", idStr).
			Str("name", request.Name)).
		Logger()

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect id: %s", idStr)})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	err = h.authors.Update(ctx, &entities.Author{
		Id:   id,
		Name: request.Name,
	})
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")
}

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	idStr := ps.ByName("id")
	lgr := h.lgr.With().
		Str("handler", "DeleteAuthor").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("id", idStr)).
		Logger()

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect id: %s", idStr)})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	err = h.authors.Delete(ctx, id)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")
}

func (h *Handler) AddPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(AddPostReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "AddPost").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("author_id", request.AuthorId).
			Str("title", request.Title).
			Str("content", request.Content).
			Time("created_at", request.CreatedAt)).
		Logger()

	err = h.posts.Add(ctx, &entities.Post{
		AuthorId:  request.AuthorId,
		Title:     request.Title,
		Content:   request.Content,
		CreatedAt: request.CreatedAt,
	})
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "ListPosts").
		Str(constants.RequestIdKey, requestId).
		Logger()

	listPosts, err := h.posts.List(ctx)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")

	resp, _ := json.Marshal(ListPostsResp{Posts: listPosts})
	fmt.Fprintf(w, string(resp))
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(UpdatePostReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	idStr := ps.ByName("id")
	lgr := h.lgr.With().
		Str("handler", "UpdatePost").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("id", idStr).
			Uint64("author_id", request.AuthorId).
			Str("title", request.Title).
			Str("content", request.Content).
			Time("created_at", request.CreatedAt)).
		Logger()

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect id: %s", idStr)})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	err = h.posts.Update(ctx, &entities.Post{
		Id:        id,
		AuthorId:  request.AuthorId,
		Title:     request.Title,
		Content:   request.Content,
		CreatedAt: request.CreatedAt,
	})
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	idStr := ps.ByName("id")
	lgr := h.lgr.With().
		Str("handler", "DeletePost").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("id", idStr)).
		Logger()

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect id: %s", idStr)})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	err = h.posts.Delete(ctx, id)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")
}
