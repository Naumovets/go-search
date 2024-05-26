CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    task_status INTEGER NOT NULL,
    url VARCHAR(2047) NOT NULL UNIQUE
);

-- Установка значения по умолчанию после создания таблицы (опционально)
ALTER TABLE task ALTER COLUMN task_status SET DEFAULT 0; 

CREATE INDEX IF NOT EXISTS task__task_status__idx on task (task_status);