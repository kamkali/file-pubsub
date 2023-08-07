package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"pubsub-assignment/internal/domain"
	"pubsub-assignment/internal/domain/mocks"
	"pubsub-assignment/internal/server/schema"
	"testing"
)

func TestReadLine(t *testing.T) {
	mockQS := mocks.NewFileQueueService(t)
	mockQS.
		On("ReadFile", mock.Anything).
		Return(domain.File{
			Name:  "test.txt",
			Lines: []string{"two", "lines"},
		}, nil)

	s := &Server{
		queueService: mockQS,
	}
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.readLine().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var resp schema.FileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	if resp.Name != "test.txt" {
		t.Errorf("Expected file name 'test.txt', got '%s'", resp.Name)
	}
	if resp.Content != "two\nlines" {
		t.Errorf("Expected file content 'line1\nline2', got '%s'", resp.Content)
	}
}

func TestWriteLine(t *testing.T) {
	mockQS := mocks.NewFileQueueService(t)
	mockQS.
		On("WriteFile", mock.Anything, domain.File{
			Name: "test.txt",
			Lines: []string{
				"line1",
				"line2",
			},
		}).
		Return(nil)
	s := &Server{
		queueService: mockQS,
	}
	file := schema.File{
		Name: "test.txt",
		Lines: []string{
			"line1",
			"line2",
		},
	}
	fileJSON, _ := json.Marshal(file)
	r, _ := http.NewRequest("POST", "/", bytes.NewReader(fileJSON))
	w := httptest.NewRecorder()
	s.writeLine().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
