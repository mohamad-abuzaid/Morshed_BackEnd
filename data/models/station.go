package models

import (
	"database/sql"
	"time"
)

type Station struct {
	ID         int64      `db:"id" json:"id"`
	NameEn     string     `db:"name_en" json:"name_en"`
	NameAr     string     `db:"name_ar" json:"name_ar"`
	ImagesURLs []string   `db:"images_urls" json:"images_urls"`
	AddressEn  string     `db:"address_en" json:"address_en"`
	AddressAr  string     `db:"address_ar" json:"address_ar"`
	Latitude   float32    `db:"latitude" json:"latitude"`
	Longitude  float32    `db:"longitude" json:"longitude"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
}

func (s Station) TableName() string {
	return "stations"
}

func (s *Station) PrimaryKey() string {
	return "id"
}

func (s *Station) SortBy() string {
	return "updated_at"
}

func (s *Station) ValidateInsert() bool {
	return s.NameEn != "" && s.NameAr != "" && len(s.ImagesURLs) > 0 &&
		s.AddressEn != "" && s.AddressAr != "" && s.Latitude > 0 && s.Longitude > 0
}

func (s *Station) Scan(rows *sql.Rows) error {
	s.CreatedAt = new(time.Time)
	s.UpdatedAt = new(time.Time)
	return rows.Scan(&s.ID, &s.NameEn, &s.NameAr, &s.ImagesURLs, &s.AddressEn, &s.AddressAr, &s.Latitude, &s.Longitude, &s.CreatedAt, &s.UpdatedAt)
}

type Stations []*Station

func (ss *Stations) Scan(rows *sql.Rows) (err error) {
	cs := *ss
	for rows.Next() {
		s := new(Station)
		if err = s.Scan(rows); err != nil {
			return
		}
		cs = append(cs, s)
	}

	if len(cs) == 0 {
		return sql.ErrNoRows
	}

	*ss = cs

	return rows.Err()
}
