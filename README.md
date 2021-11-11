May we ask Igor to do a little project in leu of sample code?
If there is a hire, the hours spent should be invoiced, so please keep track of that.
The project is:

"Build a TodoList with Go (Golang)
Design and implement a backend RESTful service in golang with CRUD  functionality that sends data to frontend clients.

Start by setting up an initial golang project using go modules. Decide and use a web framework (Gin, Echo, etc). Design a small DB schema for your todo list objects in any db of your preference. Organize your application's routes to support the CRUD paradigm. Use any ORM of your liking to connect and execute queries against the chosen DB. Dockerize your application. DB can also be a docker container. Docker compose is also an option."

The idea is to not spend more than a few hours.  You do not need to get too fancy with the capabilities of a todo list; keep that part simple. There is no single "right" way to do this and one need not spend time looking for it.  What is important, is to do an implementation and be able to articulate the reasoning behind the choices made.





TODO:

1. Create id for todo not good idea. Maybe will be good idea:composite primary key postgres
2. Create integration tests (Test routes)
3. Add expired field for todo
4. Init db in docker-compose ??? Is it need?
   1. How to deploy 
5. If obj will be too big we can use https://github.com/jmoiron/sqlx
6. Rework Makefile (I should know CI/CD tool)

