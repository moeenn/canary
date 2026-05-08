package logs

import (
	"canary/internal/database/models"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type logsRepo struct {
	db *sqlx.DB
}

const addLogEntryQuery string = `
	insert into logs (id, time, level, payload)
	values (?, ?, ?, ?)
`

func (r logsRepo) AddLogEntry(entry *models.LogEntry) error {
	_, err := r.db.Exec(addLogEntryQuery, entry.Id, entry.Time, entry.Level, entry.Payload)
	if err != nil {
		slog.Error("failed to add log entry", "error", err.Error())
		return errors.New("failed to add log entry")
	}
	return nil
}

const listRecentLogsQuery string = `
	select
		id,
		time,
		level
	from logs
	order by time desc
	limit ?
	offset ?
`

func (r logsRepo) ListRecentLogs(limit, offset int) ([]models.LogEntryMeta, error) {
	rows, err := r.db.Queryx(listRecentLogsQuery, limit, offset)
	if err != nil {
		slog.Error("failed to list recent logs", "error", err.Error())
		return nil, errors.New("failed to list logs")
	}

	//nolint:errcheck
	defer rows.Close()

	result := []models.LogEntryMeta{}
	for rows.Next() {
		var row models.LogEntryMeta
		if err := rows.StructScan(&row); err != nil {
			slog.Error("failed to scan log entry into struct", "error", err.Error())
			return nil, errors.New("failed to read logs")
		}
		result = append(result, row)
	}

	return result, nil
}

const listRecentLogsByLevelQuery string = `
	select
		id,
		time,
		level
	from logs
	where level = ?
	order by time desc
	limit ?
	offset ?
`

func (r logsRepo) ListRecentLogsByLevel(level string, limit, offset int) ([]models.LogEntryMeta, error) {
	rows, err := r.db.Queryx(listRecentLogsByLevelQuery, level, limit, offset)
	if err != nil {
		slog.Error("failed to list recent logs by level", "error", err.Error())
		return nil, errors.New("failed to list logs")
	}

	//nolint:errcheck
	defer rows.Close()

	result := []models.LogEntryMeta{}
	for rows.Next() {
		var row models.LogEntryMeta
		if err := rows.StructScan(&row); err != nil {
			slog.Error("failed to scan log entry into struct", "error", err.Error())
			return nil, errors.New("failed to read logs")
		}
		result = append(result, row)
	}

	return result, nil
}

const getLogByIdQuery string = `
	select * from logs
	where id = ?
`

func (r logsRepo) GetLogById(id string) (*models.LogEntry, error) {
	var entry models.LogEntry
	if err := r.db.Get(&entry, getLogByIdQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("log entry not found")
		}

		//nolint:gosec
		slog.Error("failed to get log entry by id", "error", err.Error())
		return nil, errors.New("failed to find log entry")
	}

	return &entry, nil
}

const clearAllLogsQuery string = `
	delete from logs
`

func (r logsRepo) ClearAllLogs() error {
	_, err := r.db.Exec(clearAllLogsQuery)
	if err != nil {
		slog.Error("failed to clear all logs", "error", err.Error())
		return errors.New("failed to clear all logs")
	}
	return nil
}

const clearLogsByDateRangeQuery string = `
	delete from logs
	where time between ? and ?
`

func (r logsRepo) ClearLogsByDateRange(start, end time.Time) error {
	_, err := r.db.Exec(clearLogsByDateRangeQuery, start, end)
	if err != nil {
		slog.Error("failed to clear logs by date range", "error", err.Error())
		return errors.New("failed to clear all logs")
	}
	return nil
}

const clearLogsOlderThanDateQuery string = `
	delete from logs
	where data < ?
`

func (r logsRepo) ClearLogsOlderThanDate(date time.Time) error {
	_, err := r.db.Exec(clearLogsOlderThanDateQuery, date)
	if err != nil {
		slog.Error("failed to clear logs older than date", "error", err.Error())
		return errors.New("failed to clear logs older than date")
	}
	return nil
}
