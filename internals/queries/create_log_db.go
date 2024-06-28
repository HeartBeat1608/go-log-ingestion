package queries

const CREATE_LOG_TABLE = `
CREATE TABLE logs (
	id INTEGER PRIMARY KEY NOT NULL,
	timestamp INTEGER NOT NULL,
	message TEXT NOT NULL
);
`
