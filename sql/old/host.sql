-- host определение

CREATE TABLE host (
	id INTEGER DEFAULT (0) PRIMARY KEY AUTOINCREMENT,
	host TEXT(32) NOT NULL UNIQUE,
	status INTEGER DEFAULT (0),
	"descriptor" TEXT(256)
	);

INSERT INTO host (host,status,"descriptor") 
VALUES
	('10.10.10.1',0,'host1'),
	('50.10.1.22',0,'host2'),
	('20.20.1.1',0,'host3'),
	('30.1.2.3',0,'host4'),
	('2.1.3.4',0,'host5'),
	('5.5.5.5',0,'host6'),
	('0.0.0.0',0,'host7');
