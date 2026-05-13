package logs

import (
	"html/template"
	"log/slog"
	"net/http"
)

type logsController struct {
	logsService logsService
	views       *template.Template
}

func (c logsController) render(w http.ResponseWriter, templateName string, args any) {
	err := c.views.ExecuteTemplate(w, templateName, args)
	if err != nil {
		slog.Error("failed to execute template", "name", templateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c logsController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.getLogs)
}

func (c logsController) getLogs(w http.ResponseWriter, r *http.Request) {
	c.render(w, "home_page.html", nil)
}
