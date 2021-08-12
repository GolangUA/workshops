# Reflect & interfaces

## 1. Theory

<a href="reflect.pdf" class="image fit">Click here to open presentation (PDF)</a>

## 2. Homework part.

You need to implement your own json.Marshal function starting with this code snippet

https://play.golang.org/p/lzGTwz7jgpk

## 3. Add tags as json package

Examples of struct field tags and their meanings:

```go
// Field appears in JSON as key "myName".
Field int `json:"myName"`

// Field appears in JSON as key "myName" and
// the field is omitted from the object if its value is empty,
// as defined above.
Field int `json:"myName,omitempty"`
```

TODO: 
 - Implement the same in the JSONEncode function

## 4. (Optional) Add more field`s types

For now JSONEncode supports string and int64 type. Add more data types

