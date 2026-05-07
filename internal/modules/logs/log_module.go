package logs

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func InitModule(db *sqlx.DB, mux *http.ServeMux) {
	logsRepo := logsRepo{db}
	logsService := logsService{logsRepo}
	logsController := logsController{logsService}
	logsController.RegisterRoutes(mux)
}
