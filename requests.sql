CREATE TABLE IF NOT EXISTS customers (
		id      BIGSERIAL PRIMARY KEY,
		name    TEXT  NOT NULL,
		phone   TEXT NOT NULL UNIQUE,
		active  BOOLEAN NOT NULL DEFAULT TRUE,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

SELECT * FROM customers;

SELECT * FROM customers WHERE id = 1;

CREATE TABLE IF NOT EXISTS managers (
	 id BIGSERIAL PRIMARY KEY,
	 name TEXT    NOT NULL,
	 salary       INTEGER NOT NULL CHECK ( salary > 0),
	 plan         INTEGER NOT NULL CHECK ( plan > 0 ),
	 boss_id      BIGINT REFERENCES managers,
	 department   TEXT,
	 active       BOOLEAN NOT NULL DEFAULT TRUE,
	 created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	 login        TEXT,
	 password     TEXT,
);