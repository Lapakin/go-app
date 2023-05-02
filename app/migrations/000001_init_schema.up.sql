BEGIN;

CREATE TABLE list_db (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    salary INT NOT NULL,
    date_receive DATE
);

COMMIT;
