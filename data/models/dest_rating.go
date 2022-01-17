package models

import (
	"database/sql"
	"time"
)

type DestRating struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	DestID    int64      `db:"dest_id" json:"dest_id"`
	Rate      int64      `db:"rate" json:"rate"`
	Comment   string     `db:"comment" json:"comment"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (r DestRating) TableName() string {
	return "destratings"
}

func (r *DestRating) PrimaryKey() string {
	return "id"
}

func (r *DestRating) SortBy() string {
	return "updated_at"
}

func (r *DestRating) ValidateInsert() bool {
	return r.UserID > 0 && r.DestID > 0 && r.Rate > 0 && r.Comment != ""
}

func (r *DestRating) Scan(rows *sql.Rows) error {
	r.CreatedAt = new(time.Time)
	r.UpdatedAt = new(time.Time)
	return rows.Scan(&r.ID, &r.UserID, &r.DestID, &r.Rate, &r.Comment, &r.UpdatedAt)
}

type DestRatings []*DestRating

func (rs *DestRatings) Scan(rows *sql.Rows) (err error) {
	cr := *rs
	for rows.Next() {
		r := new(DestRating)
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
