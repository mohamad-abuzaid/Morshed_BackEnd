package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID            int64      `db:"id" json:"id"`
	ParentID      int64      `db:"parent_id" json:"parent_id"`
	NameEn        string     `db:"name_en" json:"name_en"`
	NameAr        string     `db:"name_ar" json:"name_ar"`
	ImageURL      string     `db:"image_url" json:"image_url"`
	DescriptionEn string     `db:"description_en" json:"description_en"`
	DescriptionAr string     `db:"description_ar" json:"description_ar"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

func (ct Category) TableName() string {
	return "categories"
}

func (ct *Category) PrimaryKey() string {
	return "id"
}

func (ct *Category) SortBy() string {
	return "updated_at"
}

func (ct *Category) ValidateInsert() bool {
	return ct.ParentID > 0 && ct.NameEn != "" && ct.NameAr != "" && ct.ImageURL != "" && ct.DescriptionEn != "" && ct.DescriptionAr != ""
}

func (ct *Category) Scan(rows *sql.Rows) error {
	ct.CreatedAt = new(time.Time)
	ct.UpdatedAt = new(time.Time)
	return rows.Scan(&ct.ID, &ct.ParentID, &ct.NameEn, &ct.NameAr, &ct.ImageURL, &ct.DescriptionEn, &ct.DescriptionAr, &ct.CreatedAt, &ct.UpdatedAt)
}

type Categories []*Category

func (cts *Categories) Scan(rows *sql.Rows) (err error) {
	cct := *cts
	for rows.Next() {
		ct := new(Category)
		if err = ct.Scan(rows); err != nil {
			return
		}
		cct = append(cct, ct)
	}

	if len(cct) == 0 {
		return sql.ErrNoRows
	}

	*cts = cct

	return rows.Err()
}
