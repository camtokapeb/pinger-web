
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT(32) NOT NULL UNIQUE,
	password TEXT(32),
	role_id INTEGER NOT NULL,
	status INTEGER NOT NULL,
	email TEXT(128),
	description TEXT(128),
	FOREIGN KEY (role_id) REFERENCES role(id)
	);

INSERT INTO users (login,password,role_id,status,email,description) VALUES
	('muxomor','1111',1,0,'muxomor@mail.com','sssss'),
	('ivan','2222',2,0,'ivan@rambler.com','ddddd');
