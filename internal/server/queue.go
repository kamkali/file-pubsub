package server

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"pubsub-assignment/internal/domain"
	"pubsub-assignment/internal/server/schema"
	"strings"
)

func (s *Server) writeLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}
		var file schema.File
		if err := json.Unmarshal(body, &file); err != nil {
			s.writeErrResponse(w, err, http.StatusBadRequest, schema.ErrBadRequest)
			return
		}

		if err := s.queueService.WriteFile(ctx, domain.File{
			Name:  file.Name,
			Lines: file.Lines,
		}); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) readLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		file, err := s.queueService.ReadFile(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.writeErrResponse(w, err, http.StatusRequestTimeout, schema.ErrTimedOut)
				return
			}
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(schema.FileResponse{
			Name:    file.Name,
			Content: strings.Join(file.Lines, "\n"),
		})
		if err != nil {
			s.writeErrResponse(w, err, http.StatusInternalServerError, schema.ErrInternal)
			return
		}
		if _, err := w.Write(resp); err != nil {
			log.Println("cannot write response")
			return
		}
	}
}
