package server

import (
	"github.com/gorilla/mux"
	"moneytracker/internal/services"
	"net/http"
)

type Server struct {
	Mux     *mux.Router
	Service *services.Service
}

// NewServer конструктор структуры
func NewServer(mux *mux.Router, service *services.Service) *Server {
	return &Server{
		Mux:     mux,
		Service: service,
	}
}

// Обслуживание роутеров
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

// Init инициализация роутеров
func (s *Server) Init() {

	logAuth := mux.MiddlewareFunc(s.ValidateToken)

	generalRout := s.Mux.PathPrefix("/api/v1").Subrouter()

	authRout := generalRout.PathPrefix("/auth").Subrouter()

	authRout.HandleFunc("/registration", s.Registration)
	authRout.HandleFunc("/login", s.Login)

	moneyTracker := generalRout.PathPrefix("/money_tracker").Subrouter()
	moneyTracker.Use(logAuth)

	moneyTracker.HandleFunc("/account", s.CreateAccount).Methods("POST")
	moneyTracker.HandleFunc("/account", s.UpdateAccount).Methods("PUT")
	moneyTracker.HandleFunc("/account", s.DeleteAccount).Methods("DELETE")
	moneyTracker.HandleFunc("/operation", s.CreateOperation).Methods("POST")
	moneyTracker.HandleFunc("/types", s.GetTypes).Methods("GET")
	moneyTracker.HandleFunc("/categories", s.GetCategories).Methods("GET")
	moneyTracker.HandleFunc("/accounts", s.GetAccounts).Methods("GET")
	moneyTracker.HandleFunc("/total-balance", s.GetTotalBalance).Methods("GET")
	moneyTracker.HandleFunc("/active-accounts", s.GetActiveAccounts).Methods("GET")
	moneyTracker.HandleFunc("/reports", s.GetReports).Methods("POST")
	moneyTracker.HandleFunc("/excel-reports", s.GetExcelReports).Methods("POST")
}
