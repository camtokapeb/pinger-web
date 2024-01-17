CREATE TABLE area
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	areas_id INTEGER,
	host_id INTEGER,
	definition TEXT(256),
	UNIQUE (areas_id, host_id)
	FOREIGN KEY (host_id) REFERENCES host(id)
);

INSERT INTO area (areas_id,host_id,definition) VALUES
	(1,1,'MA4000'),
	(1,4,'MA4000'),
	(1,2,'LTP-16N'),
	(1,3,'kz'),
	(2,5,'lalalala'),
	(2,6,'ugu'),
	(3,7,'group1');
