package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
		err := rows.Scan(&note.Id, &note.Title, &note.Body)
		if err = rows.Err(); err != nil {
			return nil, nil
		}
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

// yay pointer fuckery
func (s *Store) DeleteNote(notes *[]Note, index int) error {
	sqlDelete := `delete from notes where id=?`

	if _, err := s.conn.Exec(sqlDelete, (*notes)[index].Id); err != nil {
		return err
	}

	// hopefully in-place pointer modification
	st := *notes
	*notes = append(st[:index], st[index+1:]...)

	return nil
}
