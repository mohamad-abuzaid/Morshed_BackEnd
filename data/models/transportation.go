package models

import (
	"database/sql"
	"time"
)

type Transportation struct {
	ID            int64      `db:"id" json:"id"`
	CategoryID    int64      `db:"category_id" json:"category_id"`
	NameEn        string     `db:"name_en" json:"name_en"`
	NameAr        string     `db:"name_ar" json:"name_ar"`
	CatNameEn     string     `db:"cat_en" json:"cat_en"`
	CatNameAr     string     `db:"cat_ar" json:"cat_ar"`
	ImagesURLs    []string   `db:"images_urls" json:"images_urls"`
	DescriptionEn string     `db:"description_en" json:"description_en"`
	DescriptionAr string     `db:"description_ar" json:"description_ar"`
	IsStation     bool       `db:"is_station" json:"is_station"`
	StationId     int64      `db:"station_id" json:"station_id"`
	TicketPrice   float32    `db:"ticket_price" json:"ticket_price"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

func (t Transportation) TableName() string {
	return "transportations"
}

func (t *Transportation) PrimaryKey() string {
	return "id"
}

func (t *Transportation) SortBy() string {
	return "updated_at"
}

func (t *Transportation) ValidateInsert() bool {
	return t.CategoryID > 0 && t.NameEn != "" && t.NameAr != "" && t.CatNameEn != "" && t.CatNameAr != "" && len(t.ImagesURLs) > 0 &&
		t.DescriptionEn != "" && t.DescriptionAr != "" && t.StationId > 0 && t.TicketPrice > 0
}

func (t *Transportation) Scan(rows *sql.Rows) error {
	t.CreatedAt = new(time.Time)
	t.UpdatedAt = new(time.Time)
	return rows.Scan(&t.ID, &t.CategoryID, &t.NameEn, &t.NameAr, &t.CatNameEn, &t.CatNameAr, &t.ImagesURLs, &t.DescriptionEn,
		&t.DescriptionAr, &t.IsStation, &t.StationId, &t.TicketPrice)
}

type Transportations []*Transportation

func (ts *Transportations) Scan(rows *sql.Rows) (err error) {
	ct := *ts
	for rows.Next() {
		t := new(Transportation)
		if err = t.Scan(rows); err != nil {
			return
		}
		ct = append(ct, t)
	}

	if len(ct) == 0 {
		return sql.ErrNoRows
	}

	*ts = ct

	return rows.Err()
}
