# How to use DMS in Golang
## DMS overview
* Relational Database Management Systems ([SQLite](https://www.sqlite.org/cli.html), [PostgreSQL]())
* NoSQL and NewSQL ([MongoDB](), [Redis](), [Apache Cassanra]())

## Test task:
Develop a system able to collect SDK events, count events by type in real-time and upload it to any data sotrage.
The event has next structure:

* EventType string, this refers to the type of event being sent, such as session_start, session_end, link_clicked, etc
* Ts - Unix timestamp in seconds
* Params - A key-value dictionary where the key must be a string but the value can be of any data type

## Preparetion stemp
1. Prepare enviroment: [Docker](https://www.docker.com/products/docker) [Docker Compose](https://docs.docker.com/compose/overview/)
2. Create API to handle such type of calls (replase comments from app-src.go)
3. Add messages counter
4. Handle the request 
   
   test request: `curl http://127.0.0.1:8080`
5. Prepare data to save (use [SQLite3](https://hub.docker.com/r/spartakode/sqlite3/), [switch from SQLite to PostgreSQL](https://blog.codeship.com/running-rails-development-environment-docker/))

```
touch sqltest.db

#test single sqlite3 docker image
docker run --rm -it -v sqltest.db:/db/sqltest.db spartakode/sqlite3:latest

create table messages (
id integer not null primary key autoincrement,
event_type text,
ts integer not null,
params text
);
 
insert into message (event_type, ts, params) values('test', 1234567, '{"first":1,"second":"Two"}');
 
select * from message;

#build your solution
GOOS=linux go build -o app app.go
```

6. Build & launch docker: 

```
docker build -t testdms .
docker images #find your image hash
docker run --rm -it -v sqlite.db:/app/sqlite.db <image-hash>
```
7. Test request with data:

`curl -X POST -H 'Content-Type: application/json' -d '{"eventType":"session_start","ts":1473837996,"params":{"first":1,"second":"Two"}}' http://127.0.0.1:8080`

