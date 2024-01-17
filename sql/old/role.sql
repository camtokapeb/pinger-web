
CREATE TABLE "role" (
id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT(128), description TEXT(250));

INSERT INTO "role" (name,description) VALUES
	('admin','админ'),
	('user','юзер');
