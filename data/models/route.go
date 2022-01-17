package models

import (
	"database/sql"
	"time"
)

type Route struct {
	ID            int64      `db:"id" json:"id"`
	TransId       int64      `db:"trans_id" json:"trans_id"`
	DestLat       float32    `db:"dest_lat" json:"dest_lat"`
	DestLong      float32    `db:"dest_long" json:"dest_long"`
	Eta           float32    `db:"eta" json:"eta"`
	Price         float32    `db:"price" json:"price"`
	DescriptionEn string     `db:"description_en" json:"description_en"`
	DescriptionAr string     `db:"description_ar" json:"description_ar"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

func (r Route) TableName() string {
	return "routes"
}

func (r *Route) PrimaryKey() string {
	return "id"
}

func (r *Route) SortBy() string {
	return "updated_at"
}

func (r *Route) ValidateInsert() bool {
	return r.TransId > 0 && r.DestLat > 0 && r.DestLong > 0 && r.Eta > 0 && r.Price > 0 && r.DescriptionEn != "" && r.DescriptionAr != ""
}

func (r *Route) Scan(rows *sql.Rows) error {
	r.CreatedAt = new(time.Time)
	r.UpdatedAt = new(time.Time)
	return rows.Scan(&r.ID, &r.TransId, &r.DestLat, &r.DestLong, &r.Eta, &r.Price, &r.DescriptionEn, &r.DescriptionAr, &r.CreatedAt, &r.UpdatedAt)
}

type Routes []*Route

func (rs *Routes) Scan(rows *sql.Rows) (err error) {
	cr := *rs
	for rows.Next() {
		r := new(Route)
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
