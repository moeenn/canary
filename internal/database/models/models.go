package models

type LogEntry struct {
	Id      string `db:"id"`
	Time    string `db:"time"`
	Level   string `db:"level"`
	Payload string `db:"payload"`
}

type LogEntryMeta struct {
	Id    string `db:"id"`
	Time  string `db:"time"`
	Level string `db:"level"`
}
