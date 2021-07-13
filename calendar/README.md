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

## 4. Integration with DB

Add persistence layer. Now all information about events and user should be stored inside database.

### TODO

* Organize DB connection (learn what parameters it accept and what of them should be used), make it configurable; (use Docker to start PostgreSQL)
* Use DB migration to setup schema on the start;
* Update logic to store data in the DB;
* Add integration test to cover interation with a DB; 

### Hint:
- http://go-database-sql.org our friend
	