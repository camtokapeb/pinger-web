-- areas определение

CREATE TABLE areas
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	area_id INTEGER,
	users_id INTEGER,
	UNIQUE (area_id, users_id)
	FOREIGN KEY (users_id) REFERENCES users(id)
);

INSERT INTO areas (area_id,users_id) VALUES (1,1), (2,2), (3,2);
