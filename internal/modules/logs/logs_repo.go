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
	insert into logs (time, level, payload)
	values (?, ?, ?)
`

type addLogEntryArgs struct {
	Time    string
	Level   string
	Payload string
}

func (r logsRepo) AddLogEntry(args addLogEntryArgs) error {
	_, err := r.db.Exec(addLogEntryQuery, args.Time, args.Level, args.Payload)
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
	from logs l
	order by l.time desc
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
	from logs l
	where l.level = ?
	order by l.time desc
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
	select * from logs l
	where l.id = ?
`

func (r logsRepo) GetLogById(id int64) (*models.LogEntry, error) {
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
	delete from logs l
	where l.time between ? and ?
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
	delete from logs l
	where l.date < ?
`

func (r logsRepo) ClearLogsOlderThanDate(date time.Time) error {
	_, err := r.db.Exec(clearLogsOlderThanDateQuery, date)
	if err != nil {
		slog.Error("failed to clear logs older than date", "error", err.Error())
		return errors.New("failed to clear logs older than date")
	}
	return nil
}
