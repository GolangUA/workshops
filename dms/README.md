# How to use DMS in Golang
## DMS overview
* Relational Database Management Systems
* NoSQL and NewSQL

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
4. Parse data from request 
   (test request: `curl -X POST -H 'Content-Type: application/json' -d '{"eventType":"session_start","ts":1473837996,"params":{"first":1,"second":"Two"}}' http://127.0.0.1:8080`
5. Prepare data to save (use [sqlite3](https://hub.docker.com/r/spartakode/sqlite3/))
6. Launch docker-compose: `docker-compose up -d`
7. 

