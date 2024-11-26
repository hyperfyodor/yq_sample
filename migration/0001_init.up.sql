CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    type INT CHECK (type BETWEEN 0 AND 9) NOT NULL,
    value INT CHECK (value BETWEEN 0 AND 99) NOT NULL,
    state TEXT CHECK (state IN ('received', 'processing', 'done')) DEFAULT 'received' NOT NULL,
    creation_time FLOAT DEFAULT extract(epoch from now()) NOT NULL,
    last_update_time FLOAT DEFAULT extract(epoch from now()) NOT NULL
);