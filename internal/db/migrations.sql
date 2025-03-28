CREATE TABLE
    accounts (
        id SERIAL PRIMARY KEY,
        owner VARCHAR(255) NOT NULL,
        balance NUMERIC(15, 2) NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT NOW ()
    );

CREATE TABLE
    transactions (
        id SERIAL PRIMARY KEY,
        account_id INT REFERENCES accounts (id) ON DELETE CASCADE,
        amount NUMERIC(15, 2) NOT NULL,
        type VARCHAR(10) NOT NULL CHECK (type IN ('deposit', 'withdrawal')),
        created_at TIMESTAMP DEFAULT NOW ()
    );