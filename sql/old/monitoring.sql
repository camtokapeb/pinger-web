CREATE TABLE monitoring (
	id INTEGER DEFAULT (0) PRIMARY KEY AUTOINCREMENT,
	date_time TEXT,
	host_id INTEGER,
	status INTEGER DEFAULT (0),
	time_response NUMERIC DEFAULT (0.0),
	FOREIGN KEY (host_id) REFERENCES host(id)
);