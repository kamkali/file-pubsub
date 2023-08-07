package schema

const (
    ErrInternal   = "Internal server error"
    ErrBadRequest = "Bad request"
    ErrTimedOut   = "Timed out"
)

type ServerError struct {
    Description string `json:"description"`
}

type File struct {
    Name  string   `json:"name,omitempty"`
    Lines []string `json:"lines,omitempty"`
}

type FileResponse struct {
    Name    string `json:"name,omitempty"`
    Content string `json:"content,omitempty"`
}
