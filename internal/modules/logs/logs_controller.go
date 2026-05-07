package logs

import "net/http"

type logsController struct {
	logsService logsService
}

func (c logsController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/logs", c.getLogs)
	mux.HandleFunc("GET /api/logs/{id}", c.getLogEntryById)
	mux.HandleFunc("POST /api/logs", c.addLogEntry)
	mux.HandleFunc("DELETE /api/logs", c.clearLogs)
}

// TODO: implement.
func (c logsController) getLogs(w http.ResponseWriter, r *http.Request) {
}

// TODO: implement.
func (c logsController) getLogEntryById(w http.ResponseWriter, r *http.Request) {
	// id := r.PathValue("id")
}

// TODO: implement.
func (c logsController) addLogEntry(w http.ResponseWriter, r *http.Request) {
}

// TODO: implement.
func (c logsController) clearLogs(w http.ResponseWriter, r *http.Request) {
}
