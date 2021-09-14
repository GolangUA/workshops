# Calendar

## 1. REST, routing

Create a service that will store client events (title, description, point of time (date and time in certain timezone) of start, duration, notes [optional])

### REST CRUD for events;
    
* GET events (with filter by date/data-range, time range, title), 
* GET event by ID, 
* POST create an event, 
* PUT update an event (title, description, data & time, timezone, duration, notes), 
* DELETE delete event.

```yaml
swagger: "3.9"
info:
  title: Calendar
  version: '1.0'
  description: Simple event calendar server
  license:
    name: Open license
host: 'localhost:5000'
basePath: "/"
tags:
  - name: auth
    description: Basic auth operations
  - name: user
    description: Basic user operations
  - name: events
    description: Basic events operations 
  - name: event
    description: Basic event operations
schemas: "http"
paths:
  /login:
    post:
      tags:
       - auth
      summary: Logs user into the system
      description: ''
      consumes:
        - application/json
      parameters:
        - in: body
          name: body
          description: Auth form data as json
          required: true
          schema:
            $ref: '#/definitions/Auth'
      responses:
        '200':
          description: Successfully saved"

  /logout:
    get:
      tags:
        - auth
      summary: Logs out current logged in user session
      description: ''
      operationId: logoutUser
      produces:
        - text/plain
      parameters: []
      responses:
        default:
          description: successful loged out

  /api/user:
    put:
      tags: 
       - user
      summary: Update user's timezone
      description: 'This operation can be done only for loged in users'
      consumes:
        - application/json
      parameters:
        - in: body
          name: body
          description: User's timezone data
          required: true
          schema:
            $ref: '#/definitions/User'
      responses:
       '401':
          description: Unathorized access
       '200':
          description: Successfully saved

  /api/events:
    get:
      tags: 
       - events
      summary: Get events
      description: 'This operation can be done only for loged in users'
      consumes:
        - application/json
      parameters:
        - name: title
          in: query
          description: The title of the event
          required: false
          type: string
          default: birthday
        - name: timezone
          in: query
          description: The timezone of events
          required: false
          type: string
          default: Europe/Riga
        - name: dateFrom
          in: query
          description: Events must be after or on this date
          required: false
          type: string
          format: date
          default: '2021-09-01'
        - name: dateTo
          in: query
          description: Events must be before this date
          required: false
          type: string
          format: date
          default: '2021-09-01'
        - name: timeFrom
          in: query
          description: Events must be after or on this time
          required: false
          type: string
          format: time
          default: 08:00
        - name: timeTo
          in: query
          description: Events must be before this time
          required: false
          type: string
          format: time
          default: 10:00
      responses:
       '401':
          description: Unathorized access
       '200':
          description: Successful operation
          schema:
            type: array
            items:
              $ref: '#/definitions/Event'
    post:
      tags: 
       - event
      summary: Create event
      description: 'This operation can be done only for loged in users'
      consumes:
        - application/json
      parameters:
        - in: body
          name: body
          description: Created event object
          required: true
          schema:
            $ref: '#/definitions/Event'
      responses:
       '401':
          description: Unathorized access
       '201':
          description: Successfully saved

  /api/event/{id}:
    get:
      tags: 
       - event
      summary: Get event by id
      description: 'This operation can be done only for loged in users'
      parameters:
        - name: id
          in: path
          description: 'Event ID'
          required: true
          type: integer
          default: 1
      responses:
       '401':
          description: Unathorized access
       '200':
          description: Successful operation
          schema:
            $ref: '#/definitions/Event'
    put:
      tags:
        - event
      summary: Update event
      description: 'This operation can be done only for loged in users'
      consumes:
        - application/json
      parameters:
        - in: body
          name: body
          description: Updated event object
          required: true
          schema:
            $ref: '#/definitions/Event'
      responses:
        '401':
          description: Unathorized access
        '201':
          description: Successfully saved
definitions:
  Auth:
    type: object
    properties:
      Username:
        type: string
      Password:
        type: string 
  User:
    type: object
    properties:
      login:
        type: string
      timezone:
        type: string       
        
  Event:
    type: object
    properties: 
      id: 
        type: string
      title:
        type: string
      description:
        type: string
      time:
        type: string
      timezone:
        type: string
      duration:
        type: integer
        format: int32
      notes:
        type: array
        items: 
          type: string
```

### Hints

* Start from utilizing standard packages: `net/http`; 
  * [tutorial for simple http server with tests](https://quii.gitbook.io/learn-go-with-tests/build-an-application/app-intro)
  * [http handlers](https://quii.gitbook.io/learn-go-with-tests/questions-and-answers/http-handlers-revisited)
  * [repo with code from learn go with tests](https://github.com/quii/learn-go-with-tests)
* Use `net/http/httptest` for creation tests;
* Fineout how to convert & store time between different timezones with package `time`;
* All data should be stored in memory;
* Use Postman to validate your server;
* Use [is](https://github.com/matryer/is) and [moq](https://github.com/matryer/moq) in tests;

```go
package main

import "fmt"
import "time"

func main() {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		panic(err)
	}

	t := time.Now().In(loc)
	fmt.Println(t)
}
```

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
    - requests per second, 
    - avg. requests per user per second, 
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
- [example of Prometheus metrics](https://prometheus.io/docs/guides/go-application)


## 4. Integration with DB

Add persistence layer. Now all information about events and user should be stored inside database.

### TODO

* Organize DB connection (learn what parameters it accept and what of them should be used), make it configurable; (use Docker to start PostgreSQL)
* Use DB migration to setup schema on the start;
* Update logic to store data in the DB;
* Add integration test to cover interation with a DB; 

### Hint:
- http://go-database-sql.org our friend
- [use docker-compose for DB setup](https://medium.com/analytics-vidhya/getting-started-with-postgresql-using-docker-compose-34d6b808c47c)

```sh
docker-compose up -d db # for start DB, all initial queries should be in initdb.d folder
docker-compose logs -f db #for reading logs 
dc exec db /bin/sh # enter DB container
psql --host=db --username=gouser --dbname=gotest # launch PG client
psql -U postgres # as superuser, alternative variant

```

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