package entity

import "time"

type Entity struct {
	Id        string     `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
