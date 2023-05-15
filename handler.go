package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type HTTPError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type ListResponse struct {
	Items []GenericItem `json:"items"`
}

type Handler struct {
	r *chi.Mux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(svc Service) *Handler {
	h := new(Handler)

	h.r = chi.NewRouter()

	h.r.Get("/{kind}", ListHandler(svc))
	h.r.Post("/{kind}", CreateHandler(svc))
	h.r.Get("/{kind}/{id}", ReadHandler(svc))
	h.r.Put("/{kind}/{id}", ReplaceHandler(svc))
	h.r.Delete("/{kind}/{id}", DeleteHandler(svc))

	return h
}

func ListHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")

		res, err := svc.List(r.Context(), kind)
		if err != nil {
			switch {
			case errors.As(err, &KindNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Invalid kind",
					Error:   err.Error(),
				})
			default:
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Unexpected error occurred",
					Error:   err.Error(),
				})
			}

			return
		}

		_ = json.NewEncoder(w).Encode(ListResponse{Items: res})
	}
}

func CreateHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")

		var req GenericItem

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			// TODO: check whether it's a client error or not
			w.WriteHeader(http.StatusBadRequest)

			_ = json.NewEncoder(w).Encode(HTTPError{
				Message: "Invalid request",
				Error:   err.Error(),
			})

			return
		}

		err = svc.Create(r.Context(), kind, req)
		if err != nil {
			switch {
			case errors.As(err, &ItemExistsError{}):
				w.WriteHeader(http.StatusConflict)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Item exists",
					Error:   err.Error(),
				})
			default:
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Unexpected error occurred",
					Error:   err.Error(),
				})
			}

			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(req)
	}
}

func ReadHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")
		id := chi.URLParam(r, "id")

		res, err := svc.Read(r.Context(), kind, id)
		if err != nil {
			switch {
			case errors.As(err, &KindNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Invalid kind",
					Error:   err.Error(),
				})
			case errors.As(err, &ItemNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Item not found",
					Error:   err.Error(),
				})
			default:
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Unexpected error occurred",
					Error:   err.Error(),
				})
			}

			return
		}

		_ = json.NewEncoder(w).Encode(res)
	}
}

func ReplaceHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")
		id := chi.URLParam(r, "id")

		var req GenericItem

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			// TODO: check whether it's a client error or not
			w.WriteHeader(http.StatusBadRequest)

			_ = json.NewEncoder(w).Encode(HTTPError{
				Message: "Invalid request",
				Error:   err.Error(),
			})

			return
		}

		err = svc.Replace(r.Context(), kind, id, req)
		if err != nil {
			switch {
			case errors.As(err, &KindNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Invalid kind",
					Error:   err.Error(),
				})
			case errors.As(err, &ItemNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Item not found",
					Error:   err.Error(),
				})
			default:
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Unexpected error occurred",
					Error:   err.Error(),
				})
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")
		id := chi.URLParam(r, "id")

		err := svc.Delete(r.Context(), kind, id)
		if err != nil {
			switch {
			case errors.As(err, &KindNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Invalid kind",
					Error:   err.Error(),
				})
			case errors.As(err, &ItemNotFoundError{}):
				w.WriteHeader(http.StatusNotFound)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Item not found",
					Error:   err.Error(),
				})
			default:
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(HTTPError{
					Message: "Unexpected error occurred",
					Error:   err.Error(),
				})
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
