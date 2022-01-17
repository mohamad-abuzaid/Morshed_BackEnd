package models

import (
	"database/sql"
	"time"
)

type TransRating struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	TransID   int64      `db:"trans_id" json:"trans_id"`
	Rate      int64      `db:"rate" json:"rate"`
	Comment   string     `db:"comment" json:"comment"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (r TransRating) TableName() string {
	return "transratings"
}

func (r *TransRating) PrimaryKey() string {
	return "id"
}

func (r *TransRating) SortBy() string {
	return "updated_at"
}

func (r *TransRating) ValidateInsert() bool {
	return r.UserID > 0 && r.TransID > 0 && r.Rate > 0 && r.Comment != ""
}

func (r *TransRating) Scan(rows *sql.Rows) error {
	r.CreatedAt = new(time.Time)
	r.UpdatedAt = new(time.Time)
	return rows.Scan(&r.ID, &r.UserID, &r.TransID, &r.Rate, &r.Comment, &r.UpdatedAt)
}

type TransRatings []*TransRating

func (rs *TransRatings) Scan(rows *sql.Rows) (err error) {
	cr := *rs
	for rows.Next() {
		r := new(TransRating)
		if err = r.Scan(rows); err != nil {
			return
		}
		cr = append(cr, r)
	}

	if len(cr) == 0 {
		return sql.ErrNoRows
	}

	*rs = cr

	return rows.Err()
}
