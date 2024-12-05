CREATE TABLE IF NOT EXISTS customers (
		id      BIGSERIAL PRIMARY KEY,
		name    TEXT  NOT NULL,
		phone   TEXT NOT NULL UNIQUE,
		active  BOOLEAN NOT NULL DEFAULT TRUE,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	