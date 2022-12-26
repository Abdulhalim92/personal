package server

import (
	"context"
	"moneytracker/internal/server/helpers"
	"moneytracker/pkg/logging"
	"net/http"
)

const (
	UserId     = "user_id"
	TypeId     = "type_id"
	CategoryId = "category_id"
	StartDate  = "start_date"
	EndDate    = "end_date"
)

var log = logging.GetLogger()

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// перехват паники todo
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		token := r.Header.Get("token")

		userId, err := s.Service.ValidateToken(token)
		if err != nil {
			log.Println("Non-existent or expired token", err)
			helpers.NewErrorResponse(w, 500, err.Error())
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserId, userId))

		next.ServeHTTP(w, r)
	})
}
