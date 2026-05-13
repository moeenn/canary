package logs

import (
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func InitModule(db *sqlx.DB, views *template.Template, mux *http.ServeMux) {
	logsRepo := logsRepo{db}
	logsService := logsService{logsRepo}
	logsController := logsController{logsService, views}
	logsApiController := logsApiController{logsService}
	logsController.RegisterRoutes(mux)
	logsApiController.RegisterRoutes(mux)
}
