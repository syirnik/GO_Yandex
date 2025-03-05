package api

import (
	"net/http"

	"github.com/syirnik/GO_Yandex/internal/application"
)

// Server представляет HTTP-сервер
type Server struct {
	Handler *Handler
}

// NewServer создает новый сервер
func NewServer(app *application.Application) *Server {
	return &Server{Handler: NewHandler(app)}
}

// Start запускает сервер
func (s *Server) Start(port string) error {
	http.HandleFunc("/api/v1/calculate", s.Handler.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", s.Handler.HandleExpressions)
	http.HandleFunc("/api/v1/expressions/", s.Handler.HandleGetExpressionByID)
	http.HandleFunc("/api/v1/result/", s.Handler.HandleGetResult) // обработчик для получения результата

	// Обработчики задач
	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.Handler.HandleGetTask(w, r)
		case http.MethodPost:
			s.Handler.HandlePostTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return http.ListenAndServe(port, nil)
}
