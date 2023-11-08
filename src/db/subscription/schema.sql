CREATE TABLE IF NOT EXISTS subscription (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    schedule_id INTEGER NOT NULL,
    FOREIGN KEY (schedule_id) REFERENCES schedule (id) ON DELETE CASCADE,
    CONSTRAINT unique_user_schedule UNIQUE (user_id, schedule_id)
) strict;