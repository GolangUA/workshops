# How to use DMS in Golang
## DMS overview
* Relational Database Management Systems ([SQLite](https://www.sqlite.org/cli.html), [PostgreSQL](https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/05.4.html))
* NoSQL and NewSQL ([MongoDB](https://godoc.org/labix.org/v2/mgo), [Redis](https://github.com/go-redis/redis), [Apache Cassandra](https://academy.datastax.com/resources/getting-started-apache-cassandra-and-go))

## Test task:
Develop a TODO system: API with CRUD operations:

- create new task (id (int), alias (string), description (string), type (set of string: [urgent, important, general]), tags (set of srting: [personal, work, vacation]), timestamp (int), estimate-time (string), real-time (string), reminders (set of strings: ["3h", "15m"]))
- read tasks, or one task by id/alias
- update task
- delete task

## Preparation steps
1. Prepare environment: [Docker](https://www.docker.com/products/docker) [Docker Compose](https://docs.docker.com/compose/overview/)
2. Create API to handles such type of calls (replace comments from app.go)
3. Implement of interaction with database

- implement DB driver
- add support for docker container
- add DB schema (if needed)
- implement CRUD operations

4. Handle the request 
   
   test request: `curl http://127.0.0.1:8080`
   
5. Prepare data to save (use [SQLite3](https://hub.docker.com/r/spartakode/sqlite3/), [switch from SQLite to PostgreSQL](https://blog.codeship.com/running-rails-development-environment-docker/))

```
#create DB file for SQLite
touch sqltest.db

#build your solution
GOOS=linux go build -o app app.go
```

6. Build & launch docker: 

```
docker build -t testdms:latest .
docker images #find your image hash
docker run --rm -it -v sqlite.db:/app/sqlite.db <image-hash>
```
7. Test request with data:
`curl -X POST -H 'Content-Type: application/json' -d '{"alias":"go-dms-workshop","desc":"Create app and try it with different DMS", "type":"important", "ts":1473837996,"tags":["Golang","Workshop","DMS"],"etime":"4h","rtime":"8h","reminders":["3h", "15m"]}' http://127.0.0.1:8080`

