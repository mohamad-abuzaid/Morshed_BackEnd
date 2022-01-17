package models

import (
	"database/sql"
	"time"
)

type Destination struct {
	ID            int64      `db:"id" json:"id"`
	CategoryID    int64      `db:"category_id" json:"category_id"`
	NameEn        string     `db:"name_en" json:"name_en"`
	NameAr        string     `db:"name_ar" json:"name_ar"`
	CatNameEn     string     `db:"cat_en" json:"cat_en"`
	CatNameAr     string     `db:"cat_ar" json:"cat_ar"`
	ImagesURLs    []string   `db:"images_urls" json:"images_urls"`
	DescriptionEn string     `db:"description_en" json:"description_en"`
	DescriptionAr string     `db:"description_ar" json:"description_ar"`
	AddressEn     string     `db:"address_en" json:"address_en"`
	AddressAr     string     `db:"address_ar" json:"address_ar"`
	Latitude      float32    `db:"latitude" json:"latitude"`
	Longitude     float32    `db:"longitude" json:"longitude"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

func (d Destination) TableName() string {
	return "destinations"
}

func (d *Destination) PrimaryKey() string {
	return "id"
}

func (d *Destination) SortBy() string {
	return "updated_at"
}

func (d *Destination) ValidateInsert() bool {
	return d.CategoryID > 0 && d.NameEn != "" && d.NameAr != "" && d.CatNameEn != "" && d.CatNameAr != "" && len(d.ImagesURLs) > 0 &&
		d.DescriptionEn != "" && d.DescriptionAr != "" && d.AddressEn != "" && d.AddressAr != "" && d.Latitude > 0 && d.Longitude > 0
}

func (d *Destination) Scan(rows *sql.Rows) error {
	d.CreatedAt = new(time.Time)
	d.UpdatedAt = new(time.Time)
	return rows.Scan(&d.ID, &d.CategoryID, &d.NameEn, &d.NameAr, &d.CatNameEn, &d.CatNameAr, &d.ImagesURLs, &d.DescriptionEn,
		&d.DescriptionAr, &d.AddressEn, &d.AddressAr, &d.Latitude, &d.Longitude, &d.CreatedAt, &d.UpdatedAt)
}

type Destinations []*Destination

func (ds *Destinations) Scan(rows *sql.Rows) (err error) {
	cd := *ds
	for rows.Next() {
		d := new(Destination)
		if err = d.Scan(rows); err != nil {
			return
		}
		cd = append(cd, d)
	}

	if len(cd) == 0 {
		return sql.ErrNoRows
	}

	*ds = cd

	return rows.Err()
}
