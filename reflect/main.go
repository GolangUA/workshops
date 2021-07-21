package main

import (
	"bytes"
	"fmt"
)

type User struct {
	Name string
	Age  int64
}

type City struct {
	Name       string
	Population int64
	GDP        int64
	Mayor      string
}

func main() {
	var u User = User{"bob", 10}

	res, err := JSONEncode(u)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	c := City{"sf", 5000000, 567896, "mr jones"}
	res, err = JSONEncode(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func JSONEncode(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")

	// TODO: check if v is a struct else return error
	// TODO: iterate over v`s reflect value using NumField()
	// use type switch to create string result of "{field}" + ": " + "{value}"
	// start with just 2 types - reflect.String and reflect.Int64

	buf.WriteString("}")
	return buf.Bytes(), nil
}
