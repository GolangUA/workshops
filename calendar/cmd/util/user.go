package main

import (
	"flag"
	"fmt"
	"github.com/Roma7-7-7/workshops/calendar/internal/config"
	"github.com/Roma7-7-7/workshops/calendar/internal/repository/postgre"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var user = flag.String("user", "", "user to be created")
var password = flag.String("password", "", "user password")
var timezone = flag.String("timezone", "Europe/Kiev", "user timezone")

func main() {
	flag.Parse()
	if *user == "" || *password == "" || *timezone == "" {
		panic("user, password and timezone flags are required")
	}
	loc, err := time.LoadLocation(*timezone)
	if err != nil {
		panic(fmt.Sprintf("invalid location \"%s\"", *timezone))
	} else if loc == time.Local {
		panic("location must not be Local")
	}

	cfg := config.GetConfig()
	repo := postgre.NewRepository(cfg.DSN())

	encrypted, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	if u, err := repo.CreateUser(*user, string(encrypted), *timezone); err != nil {
		panic(err)
	} else {
		log.Println(*u)
	}
}
