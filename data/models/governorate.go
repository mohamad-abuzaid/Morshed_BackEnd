package models

import (
	"database/sql"
	"time"
)

type Governorate struct {
	ID        int64      `db:"id" json:"id"`
	CountryID int64      `db:"country_id" json:"country_id"`
	NameEn    string     `db:"name_en" json:"name_en"`
	NameAr    string     `db:"name_ar" json:"name_ar"`
	ImageURL  string     `db:"image_url" json:"image_url"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (g Governorate) TableName() string {
	return "governorates"
}

func (g *Governorate) PrimaryKey() string {
	return "id"
}

func (g *Governorate) SortBy() string {
	return "updated_at"
}

func (g *Governorate) ValidateInsert() bool {
	return g.CountryID > 0 && g.NameEn != "" && g.NameAr != "" && g.ImageURL != ""
}

func (g *Governorate) Scan(rows *sql.Rows) error {
	g.CreatedAt = new(time.Time)
	g.UpdatedAt = new(time.Time)
	return rows.Scan(&g.ID, &g.CountryID, &g.NameEn, &g.NameAr, &g.ImageURL, &g.CreatedAt, &g.UpdatedAt)
}

type Governorates []*Governorate

func (gs *Governorates) Scan(rows *sql.Rows) (err error) {
	cg := *gs
	for rows.Next() {
		g := new(Governorate)
		if err = g.Scan(rows); err != nil {
			return
		}
		cg = append(cg, g)
	}

	if len(cg) == 0 {
		return sql.ErrNoRows
	}

	*gs = cg

	return rows.Err()
}
