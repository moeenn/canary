package logs

import (
	"canary/internal/database/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type getLogsArgs struct {
	Limit  int
	Offset int
	Level  *string
}

func (r *getLogsArgs) validate() error {
	if r.Limit <= 0 {
		return errors.New("limit cannot be zero or negative")
	}

	if r.Limit > 1000 {
		return errors.New("limit must not exceed 1000")
	}

	if r.Offset < 0 {
		return errors.New("offset cannot be less negative")
	}

	if r.Level != nil {
		if *r.Level == "" {
			r.Level = nil
		}
	}

	return nil
}

func getLogsArgsFromRequest(r *http.Request) (*getLogsArgs, error) {
	q := r.URL.Query()
	limitQ := q.Get("limit")
	offsetQ := q.Get("offset")
	level := q.Get("level")

	parsedLimit, err := strconv.Atoi(limitQ)
	if err != nil {
		parsedLimit = 100
	}

	parsedOffset, err := strconv.Atoi(offsetQ)
	if err != nil {
		parsedOffset = 0
	}

	req := getLogsArgs{
		Limit:  parsedLimit,
		Offset: parsedOffset,
		Level:  &level,
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	return &req, nil
}

type logEntryMetaResponse struct {
	Id    string `json:"id"`
	Time  string `json:"time"`
	Level string `json:"level"`
}

func logEntryMetaModelToResponse(m models.LogEntryMeta) logEntryMetaResponse {
	return logEntryMetaResponse{
		Id:    m.Id,
		Time:  m.Time,
		Level: m.Level,
	}
}

type getLogsResponse struct {
	Data []logEntryMetaResponse `json:"data"`
}

func getLogsResponseFromModels(m []models.LogEntryMeta) getLogsResponse {
	res := getLogsResponse{
		Data: make([]logEntryMetaResponse, len(m)),
	}

	for i, me := range m {
		res.Data[i] = logEntryMetaModelToResponse(me)
	}

	return res
}

type logEntryResponse struct {
	Id      string `json:"id"`
	Time    string `json:"time"`
	Level   string `json:"level"`
	Payload string `json:"payload"`
}

func logEntryResponseFromModel(m models.LogEntry) logEntryResponse {
	return logEntryResponse{
		Id:      m.Id,
		Time:    m.Time,
		Level:   m.Level,
		Payload: m.Payload,
	}
}

func addLogEntryPayloadFromRequest(r *http.Request) (map[string]any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	payload := make(map[string]any)
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	return payload, nil
}

type clearLogsArgs struct {
	Start     *string    `json:"start"`
	End       *string    `json:"end"`
	StartTime *time.Time `json:"-"`
	EndTime   *time.Time `json:"-"`
}

func (args *clearLogsArgs) validate() error {
	if args.Start != nil && *args.Start == "" {
		args.Start = nil
	}

	if args.End != nil && *args.End == "" {
		args.End = nil
	}

	if args.Start != nil {
		parsedStart, err := time.Parse(time.DateTime, *args.Start)
		if err != nil {
			return fmt.Errorf("invalid start date: %w", err)
		}
		args.StartTime = &parsedStart
	}

	if args.End != nil {
		parsedEnd, err := time.Parse(time.DateTime, *args.End)
		if err != nil {
			return fmt.Errorf("invalid end date: %w", err)
		}
		args.EndTime = &parsedEnd
	}

	return nil
}

func clearLogsArgsFromRequest(r *http.Request) (*clearLogsArgs, error) {
	q := r.URL.Query()
	start := q.Get("start")
	end := q.Get("end")

	//nolint:exhaustruct
	args := clearLogsArgs{
		Start: &start,
		End:   &end,
	}

	if err := args.validate(); err != nil {
		return nil, err
	}

	return &args, nil
}
