CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    genre TEXT,
    budget INT
);

INSERT INTO movies (title, genre, budget) VALUES ('SAW', 'horror', 500000);
INSERT INTO movies (title, genre, budget) VALUES ('TEST', 'Romance', 1000000);