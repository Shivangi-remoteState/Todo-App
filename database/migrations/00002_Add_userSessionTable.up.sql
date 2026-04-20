BEGIN;
ALTER TABLE users
ADD COLUMN archived_at timestamp;

CREATE TABLE user_session (
    id UUID primary key default gen_random_uuid(),
    user_id UUID not null references users(id),
    archived_at timestamp,
    created_at TIMESTAMP default NOW()
);

ALTER  TABLE  todos
ADD COLUMN expiry_at timestamp;

ALTER TABLE todos
ADD COLUMN archived_at timestamp;

commit;