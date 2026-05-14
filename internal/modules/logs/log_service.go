package logs

import (
	"canary/internal/database/models"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type logsService struct {
	logsRepo logsRepo
}

// this can be made more flexible for accept different variations of the raw payload.
func (s logsService) logEntryFromRawPayload(payload map[string]any) (*addLogEntryArgs, error) {
	var time string
	if t, ok := payload["time"]; ok {
		converted, ok := t.(string)
		if !ok {
			return nil, errors.New("value for 'time' is provided but it is not a valid string")
		}
		time = converted
		delete(payload, "time")
	}

	var level string
	if l, ok := payload["level"]; ok {
		converted, ok := l.(string)
		if !ok {
			return nil, errors.New("value for 'level' is provided but it is not a valid string")
		}
		level = converted
		delete(payload, "level")
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON encode payload: %w", err)
	}

	entry := addLogEntryArgs{
		Level:   level,
		Time:    time,
		Payload: string(encoded),
	}

	return &entry, nil
}

func (s logsService) AddLogEntry(payload map[string]any) error {
	logEntry, err := s.logEntryFromRawPayload(payload)
	if err != nil {
		return fmt.Errorf("failed to convert raw payload into log entry: %w", err)
	}
	return s.logsRepo.AddLogEntry(*logEntry)
}

type GetLogsArgs struct {
	Limit  *int
	Offset *int
	Level  *string
}

func (s logsService) GetLogs(args GetLogsArgs) ([]models.LogEntryMeta, error) {
	limit := 100
	if args.Limit != nil {
		limit = *args.Limit
	}

	offset := 0
	if args.Offset != nil {
		offset = *args.Offset
	}

	if args.Level == nil {
		return s.logsRepo.ListRecentLogs(limit, offset)
	}

	return s.logsRepo.ListRecentLogsByLevel(*args.Level, limit, offset)
}

func (s logsService) GetLogEntryById(id int64) (*models.LogEntry, error) {
	return s.logsRepo.GetLogById(id)
}

type ClearLogsArgs struct {
	Start *time.Time
	End   *time.Time
}

func (s logsService) ClearLogs(args ClearLogsArgs) error {
	if args.Start == nil && args.End == nil {
		return s.logsRepo.ClearAllLogs()
	}

	if args.Start != nil && args.End != nil {
		return s.logsRepo.ClearLogsByDateRange(*args.Start, *args.End)
	}

	if args.End != nil {
		return s.logsRepo.ClearLogsOlderThanDate(*args.End)
	}

	return errors.New("invalid arguments combination")
}
