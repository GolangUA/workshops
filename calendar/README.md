# Calendar

## 1. REST, routing

Create a service that will store client events (title, description, point of time (date and time in certain timezone) of start, duration, notes [optional])

### REST CRUD for events;
    
* GET events (with filter by date/data-range, time range, title), 
* GET event by ID, 
* POST create an event, 
* PUT update an event (title, description, data & time, timezone, duration, notes), 
* DELETE delete event.

### Hints

* Start from utilizing standard packages: `net/http`; 
  * [tutorial for simple http server with tests](https://quii.gitbook.io/learn-go-with-tests/build-an-application/app-intro)
  * [http handlers](https://quii.gitbook.io/learn-go-with-tests/questions-and-answers/http-handlers-revisited)
  * [repo with code from learn go with tests](https://github.com/quii/learn-go-with-tests)
* Use `net/http/httptest` for creation tests;
* Fineout how to convert & store time between different timezones with package `time`;
* All data should be stored in memory;
* Use Postman to validate your server;

## 2. Serve multiple users

Store events per user. For this you require to have authentication/authorization; User should have own time zone and all events should respect it by default.

### TODO

* Add login/logout enpoints (credantials should be stored in memory as well);
* Use JWT for authorization at other endpoints;
* Update logic to Events endpoints to authorize access only to related to user events.
* PUT for update user timezone; Time of events should be returned in a new timezone.

### Links
- [golang jwt](https://github.com/golang-jwt/jwt)
- [Hands-on with JWT in Go](https://betterprogramming.pub/hands-on-with-jwt-in-golang-8c986d1bb4c0)
- [writing middleware in Go](https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81)


## 3. Pre-pord preparation

For preparing service to the production lunch it should have:
- good test coverage: 60-70% of unit tests;
- collect metrics for observability of it health;
- structured logs for investigation; information from logs should be enough for troubleshooting;

### TODO

* Add system endpoints for return metrics; 
  - served in a separate port;
  - metrics: 
    - total number of events, 
    - number of users, 
    - requests per seconds, 
    - requests per seconds per user, 
    - goroutines used, 
    - memory used, 
    - cpu used;
* Organize load testing with `vegetta`;
* Add support of gracefull shutdown;
* Check tests coverage, add tests if it low;
* Use structured logs (with tabs; key: value format);
* Add enough context information to the error message/logs;

### Links
- [Graceful shutdown](https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97)
- [tests coverage](https://blog.golang.org/cover)
- [structured logging](https://www.client9.com/structured-logging-in-golang/)
- [logs for HTTP service](https://ribice.medium.com/http-logging-in-go-344e6fca057c)
- [load testing with vegeta](https://geshan.com.np/blog/2020/09/vegeta-load-testing-primer-with-examples/)


## 4. Integration with DB

Add persistence layer. Now all information about events and user should be stored inside database.

### TODO

* Organize DB connection (learn what parameters it accept and what of them should be used), make it configurable; (use Docker to start PostgreSQL)
* Use DB migration to setup schema on the start;
* Update logic to store data in the DB;
* Add integration test to cover interation with a DB; 

### Hint:
- http://go-database-sql.org our friend

## 5. Add gRPC

Provide same API with gRPC.

### TODO
   * define protobuff for your task and user models
   * define protobuff definition for your Task APIs  
   * use protoc tool to generate grpc stubs
   * define auth interceptor (gRPC middleware) that will reuse same logic from http auth md
   * add second main file that will serve only gRPC requests
   * write gRPC client that will call one of your APIs
   * install BloomRPC (postman analog for gRPC) to test your APIs

### Links

- https://developers.google.com/protocol-buffers/docs/gotutorial
- https://grpc.io/docs/languages/go/basics/