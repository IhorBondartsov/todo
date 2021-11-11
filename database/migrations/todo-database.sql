CREATE
OR REPLACE FUNCTION update_time()
    RETURNS TRIGGER AS $$
BEGIN
   IF
row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at := now();
RETURN NEW;
ELSE
      RETURN OLD;
END IF;
END;
$$ language 'plpgsql';


DROP SCHEMA IF EXISTS todo_app CASCADE;
CREATE SCHEMA todo_app;

CREATE TABLE todo_app.users
(
    user_id  serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE todo_app.todo_list
(
    id  serial PRIMARY KEY,
    user_id    INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    message    VARCHAR(40000),
    FOREIGN KEY (user_id)
        REFERENCES users (user_id)
);

CREATE TRIGGER prevent_timestamp_changes
    BEFORE UPDATE
    ON todo_app.todo_list
    FOR EACH ROW
    EXECUTE PROCEDURE update_time();