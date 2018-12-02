CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name TEXT,
        last_name TEXT
);

CREATE TABLE IF NOT EXISTS bankaccounts (
        id SERIAL PRIMARY KEY,
        userid SERIAL,
        account_number SERIAL,
        name TEXT,
        balance NUMBER
);