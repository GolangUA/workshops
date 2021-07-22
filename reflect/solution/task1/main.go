package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	refObjVal := reflect.ValueOf(v)
	refObjTyp := reflect.TypeOf(v)
	buf := bytes.Buffer{}
	if refObjVal.Kind() != reflect.Struct {
		return buf.Bytes(), fmt.Errorf(
			"val of kind %s is not supported",
			refObjVal.Kind(),
		)
	}
	buf.WriteString("{")
	pairs := []string{}
	for i := 0; i < refObjVal.NumField(); i++ {
		structFieldRefObj := refObjVal.Field(i)
		structFieldRefObjTyp := refObjTyp.Field(i)

		name := structFieldRefObjTyp.Name
		switch structFieldRefObj.Kind() {
		case reflect.String:
			strVal := structFieldRefObj.Interface().(string)
			pairs = append(pairs, `"`+name+`":"`+strVal+`"`)
		case reflect.Int64:
			intVal := structFieldRefObj.Interface().(int64)
			pairs = append(pairs, `"`+name+`":`+strconv.FormatInt(intVal, 10))
		default:
			return buf.Bytes(), fmt.Errorf(
				"struct field with name %s and kind %s is not supprted",
				structFieldRefObjTyp.Name,
				structFieldRefObj.Kind(),
			)
		}
	}

	buf.WriteString(strings.Join(pairs, ","))
	buf.WriteString("}")

	return buf.Bytes(), nil
}
