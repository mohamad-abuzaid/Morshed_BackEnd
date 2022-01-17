package models

import (
	"database/sql"
	"time"
)

type Country struct {
	ID        int64      `db:"id" json:"id"`
	NameEn    string     `db:"name_en" json:"name_en"`
	NameAr    string     `db:"name_ar" json:"name_ar"`
	ImageURL  string     `db:"image_url" json:"image_url"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (c Country) TableName() string {
	return "counttries"
}

func (c *Country) PrimaryKey() string {
	return "id"
}

func (c *Country) SortBy() string {
	return "updated_at"
}

func (c *Country) ValidateInsert() bool {
	return c.NameEn != "" && c.NameAr != "" && c.ImageURL != ""
}

func (c *Country) Scan(rows *sql.Rows) error {
	c.CreatedAt = new(time.Time)
	c.UpdatedAt = new(time.Time)
	return rows.Scan(&c.ID, &c.NameEn, &c.NameAr, &c.ImageURL, &c.CreatedAt, &c.UpdatedAt)
}

type Countries []*Country

func (cs *Countries) Scan(rows *sql.Rows) (err error) {
	cc := *cs
	for rows.Next() {
		c := new(Country)
		if err = c.Scan(rows); err != nil {
			return
		}
		cc = append(cc, c)
	}

	if len(cc) == 0 {
		return sql.ErrNoRows
	}

	*cs = cc

	return rows.Err()
}
