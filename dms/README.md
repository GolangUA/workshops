# How to use DMS in Golang

## Test task:
Develop a system able to collect SDK events, count events by type in real-time and upload it to any data sotrage.
The event has next structure:

* EventType string, this refers to the type of event being sent, such as session_start, session_end, link_clicked, etc
* Ts - Unix timestamp in seconds
* Params - A key-value dictionary where the key must be a string but the value can be of any data type

## Preparetion stemp
1. Create API to handle such type of calls
2. Add counter
3. Implement data validation

