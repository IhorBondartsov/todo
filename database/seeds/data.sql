INSERT INTO todo_app.users (user_id, username)
    VALUES (DEFAULT, 'Legolas');
INSERT INTO todo_app.users (user_id, username)
    VALUES (DEFAULT, 'Dúnadan');

INSERT INTO todo_app.todo_list (id, user_id, created_at, updated_at, message)
    VALUES (DEFAULT, 1, DEFAULT, DEFAULT, 'Kill more orcs than Gimli');
INSERT INTO todo_app.todo_list (id, user_id, created_at, updated_at, message)
    VALUES (DEFAULT, 1, DEFAULT, DEFAULT, 'Help Minas Tirit');
INSERT INTO todo_app.todo_list (id, user_id, created_at, updated_at, message)
    VALUES (DEFAULT, 1, DEFAULT, DEFAULT, 'To be handsome');