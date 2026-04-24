BEGIN;

-- user table
CREATE TABLE users (
    id UUID primary key default gen_random_uuid(),
    name TEXT not null ,
    email TEXT unique not null ,
    password TEXT not null ,
    created_at TIMESTAMP default NOW()
);

-- todos table
CREATE TABLE todos (
    id UUID primary key default gen_random_uuid(),
    user_id UUID  references users(id),
    name TEXT not null ,
    description TEXT,
    complete BOOLEAN default false, // is_complete
    created_at TIMESTAMP default NOW()

);
COMMIT ;