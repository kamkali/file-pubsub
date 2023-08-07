package server

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "golang.org/x/net/context"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "pubsub-assignment/internal/config"
    "pubsub-assignment/internal/domain"
    "pubsub-assignment/internal/server/schema"
    "syscall"
    "time"
)

type Server struct {
	router     *mux.Router
	config     *config.Config
	httpServer *http.Server

	queueService domain.FileQueueService
}

func New(cfg *config.Config, queueService domain.FileQueueService) (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		router: r,
		config: cfg,
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler: r,
		},
		queueService: queueService,
	}

	s.registerRoutes()

	return s, nil
}

func (s *Server) registerRoutes() {
	{ // queue routes
		s.router.HandleFunc("/line",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.writeLine()),
		).Methods("POST")

		s.router.HandleFunc("/line",
			s.withTimeout(s.config.Server.TimeoutSeconds, s.readLine()),
		).Methods("GET")
	}
}

func (s *Server) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server Started on host=%s:%s\n", s.config.Server.Host, s.config.Server.Port)

	<-done
	log.Println("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server Exited Properly")
}

func (s *Server) writeErrResponse(w http.ResponseWriter, err error, code int, desc string) {
	log.Println(fmt.Errorf("error response: %w", err).Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonErr, err := json.Marshal(schema.ServerError{Description: desc})
	if err != nil {
		return
	}
	if _, err := w.Write(jsonErr); err != nil {
		log.Println(fmt.Errorf("cannot write error response: %w", err).Error())
		return
	}
}
