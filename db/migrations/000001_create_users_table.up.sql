CREATE TABLE users
(
    id SERIAL,
    name TEXT NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone TEXT ,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)