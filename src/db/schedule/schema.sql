CREATE TABLE IF NOT EXISTS schedule (
    id INTEGER PRIMARY KEY,
    activity_id INTEGER NOT NULL UNIQUE,
    datetime TEXT NOT NULL,
    trainer TEXT NOT NULL,
    activity TEXT NOT NULL,
    pre_entry BOOLEAN NOT NULL CHECK (pre_entry IN (0, 1)),
    begin_date TEXT NOT NULL
);