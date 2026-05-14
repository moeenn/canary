package models

type LogEntry struct {
	Id      int64  `db:"id"`
	Time    string `db:"time"`
	Level   string `db:"level"`
	Payload string `db:"payload"`
}

type LogEntryMeta struct {
	Id    int64  `db:"id"`
	Time  string `db:"time"`
	Level string `db:"level"`
}
