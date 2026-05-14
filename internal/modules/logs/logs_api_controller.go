package logs

import (
	"canary/internal/libs/server"
	"net/http"
	"strconv"
)

type logsApiController struct {
	logsService logsService
}

func (c logsApiController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/logs", c.getLogs)
	mux.HandleFunc("GET /api/logs/{id}", c.getLogEntryById)
	mux.HandleFunc("POST /api/logs", c.addLogEntry)
	mux.HandleFunc("DELETE /api/logs", c.clearLogs)
}

func (c logsApiController) getLogs(w http.ResponseWriter, r *http.Request) {
	args, err := getLogsArgsFromRequest(r)
	if err != nil {
		server.SendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	logs, err := c.logsService.GetLogs(GetLogsArgs{
		Limit:  &args.Limit,
		Offset: &args.Offset,
		Level:  args.Level,
	})

	if err != nil {
		server.SendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	res := getLogsResponseFromModels(logs)
	server.SendOkResponse(w, res)
}

func (c logsApiController) getLogEntryById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry, err := c.logsService.GetLogEntryById(parsedId)
	if err != nil {
		server.SendErrorResponse(w, http.StatusNotFound, err)
		return
	}

	res := logEntryResponseFromModel(*entry)
	server.SendOkResponse(w, res)
}

func (c logsApiController) addLogEntry(w http.ResponseWriter, r *http.Request) {
	payload, err := addLogEntryPayloadFromRequest(r)
	if err != nil {
		server.SendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := c.logsService.AddLogEntry(payload); err != nil {
		server.SendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c logsApiController) clearLogs(w http.ResponseWriter, r *http.Request) {
	args, err := clearLogsArgsFromRequest(r)
	if err != nil {
		server.SendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = c.logsService.ClearLogs(ClearLogsArgs{
		Start: args.StartTime,
		End:   args.EndTime,
	})

	if err != nil {
		server.SendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusGone)
}
