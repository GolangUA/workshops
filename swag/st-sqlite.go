package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
)

type sqliteDr struct {
	db *sql.DB
	l  chan struct{}
}

func (s *sqliteDr) init() error {
	db, err := sql.Open("sqlite3", "sqltest.db")
	if err != nil {
		log.Fatalf("sqlite.db open fail: %v", err)
		return err
	}
	db.Exec(`create table tasks (
	id integer not null primary key autoincrement,
	alias text,
	desc text,
	category text,
	tags text,
	ts integer not null,
	est_time text,
	real_time text,
	reminders text
	)`)
	s.l = make(chan struct{}, 1)
	s.db = db
	return nil
}

func (s *sqliteDr) Create(t Task) error {
	var res driver.Result
	var err error
	res, err = s.db.Exec(fmt.Sprintf("insert into tasks(alias, desc, category, tags, ts, est_time, real_time, reminders) values('%s','%s','%s','%s','%d','%s','%s','%s')", t.Alias, t.Desc, strings.Join(t.Category, ","), strings.Join(t.Tags, ","), t.Ts, t.EstTime, t.RealTime, strings.Join(t.Reminders, ",")))
	log.Printf("result of insert: %#v of (%#v)\n", res, t)
	return err
}

func (s *sqliteDr) ReadById(id *int64) (TaskList, error) {
	return s.read(id)
}

func (s *sqliteDr) ReadByAlias(alias *string) (TaskList, error) {
	return s.read(alias)
}

func (s *sqliteDr) read(val interface{}) (TaskList, error) {
	var rows *sql.Rows
	var err error
	log.Println("Read from db by %v", val)
	if val == nil {
		rows, err = s.db.Query("select id, alias, desc, category, tags, ts, est_time, real_time, reminders from tasks")
	} else {
		switch v := val.(type) {
		case *int64:
			rows, err = s.db.Query(fmt.Sprintf("select id, alias, desc, category, tags, ts, est_time, real_time, reminders from tasks where id = %d", *v))
		case *string:
			rows, err = s.db.Query(fmt.Sprintf("select id, alias, desc, category, tags, ts, est_time, real_time, reminders from tasks where alias = %s", *v))
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tl := TaskList{}
	var catSet, tagSet, remSet string
	for rows.Next() {
		t := &Task{}
		err = rows.Scan(&t.ID, &t.Alias, &t.Desc, &catSet, &tagSet, &t.Ts, &t.EstTime, &t.RealTime, &remSet)
		if err != nil {
			return tl, err
		}
		t.Category = strings.Split(catSet, ",")
		t.Tags = strings.Split(tagSet, ",")
		t.Reminders = strings.Split(remSet, ",")
		tl = append(tl, *t)
	}
	return tl, rows.Err()
}

func (s *sqliteDr) Update(t Task) error {
	s.l <- struct{}{}
	defer func() { <-s.l }()
	var res driver.Result
	var err error
	res, err = s.db.Exec(fmt.Sprintf("update tasks set alias= '%s', desc='%s', category='%s', tags='%s', ts=%d, est_time='%s', real_time='%s', reminders='%s') where id = %d", t.Alias, t.Desc, strings.Join(t.Category, ","), strings.Join(t.Tags, ","), t.Ts, t.EstTime, t.RealTime, strings.Join(t.Reminders, ","), t.ID))
	log.Printf("result of update: %#v of (%#v)\n", res, t)
	return err
}

func (s *sqliteDr) Delete(t Task) error {
	var res driver.Result
	var err error
	res, err = s.db.Exec(fmt.Sprintf("delete from tasks where id = %d", t.ID))
	log.Printf("result of update: %#v of (%#v)\n", res, t)
	return err
}
