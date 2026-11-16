package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Note struct {
	Id    int64
	Title string
	Body  string
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		return err
	}

	sqlSt := `create table if not exists notes (
		id integer not null primary key,
		title text not null,
		body text not null
	);`

	if _, err = s.conn.Exec(sqlSt); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetNotes() ([]Note, error) {
	rows, err := s.conn.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}

	notes := []Note{}

	defer rows.Close()

	for rows.Next() {
		note := Note{}
		rows.Scan(&note.Id, &note.Title, &note.Body)
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Store) SaveNote(note Note) error {
	if note.Id == 0 {
		note.Id = time.Now().UTC().Unix()
	}

	sqlAdd := `insert into notes (id, title, body)
	values(?, ?, ?)
	on conflict(id) do update
	set title=excluded.title, body=excluded.body;`

	if _, err := s.conn.Exec(sqlAdd, note.Id, note.Title, note.Body); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteNote(note Note, notes []Note, index int) ([]Note, error) {
	sqlDelete := `delete from notes where id=?`

	if _, err := s.conn.Exec(sqlDelete, note.Id); err != nil {
		return nil, err
	}

	// slice deletion
	n := append(notes[:index], notes[index+1:]...)

	return n, nil
}
