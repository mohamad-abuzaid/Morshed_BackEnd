package repositories

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"morshed/data/engine/sql"
	"morshed/data/models"

	"github.com/kataras/iris/v12"
)

// ProductRepository represents the product models service.
// Note that the given models (request) should be already validated
// before service's calls.
type ProductRepository struct {
	Ctx iris.Context
	*sql.Repository
	rec sql.Record
}

// NewProductRepository returns a new product service to communicate with the database.
func NewProductRepository(db sql.Database) *ProductRepository {
	return &ProductRepository{Repository: sql.NewRepository(db, new(models.Product))}
}

func (r *ProductRepository) Size(id int64) (int64, error) {
	total, err := r.Count(r.Ctx)
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (r *ProductRepository) SELECT(id int64) (prod models.Product, err error) {
	err = r.GetByID(r.Ctx, prod, id)
	return
}

// Insert stores a product to the database and returns its ID.
func (r *ProductRepository) Insert(e models.Product) (int64, error) {
	if !e.ValidateInsert() {
		return 0, sql.ErrUnprocessable
	}

	q := fmt.Sprintf(`INSERT INTO %s (category_id, title, image_url, price, description)
	VALUES (?,?,?,?,?);`, e.TableName())

	res, err := r.DB().Exec(r.Ctx, q, e.CategoryID, e.Title, e.ImageURL, e.Price, e.Description)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// BatchInsert inserts one or more products at once and returns the total length created.
func (r *ProductRepository) BatchInsert(ctx context.Context, products []models.Product) (int, error) {
	if len(products) == 0 {
		return 0, nil
	}

	var (
		valuesLines []string
		args        []interface{}
	)

	for _, p := range products {
		if !p.ValidateInsert() {
			// all products should be "valid", we don't skip, we cancel.
			return 0, sql.ErrUnprocessable
		}

		valuesLines = append(valuesLines, "(?,?,?,?,?)")
		args = append(args, []interface{}{p.CategoryID, p.Title, p.ImageURL, p.Price, p.Description}...)
	}

	q := fmt.Sprintf("INSERT INTO %s (category_id, title, image_url, price, description) VALUES %s;",
		r.RecordInfo().TableName(),
		strings.Join(valuesLines, ", "))

	res, err := r.DB().Exec(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	n := sql.GetAffectedRows(res)
	return n, nil
}

// Update updates a product based on its `ID` from the database
// and returns the affected numbrer (0 when nothing changed otherwise 1).
func (r *ProductRepository) Update(ctx context.Context, e models.Product) (int, error) {
	q := fmt.Sprintf(`UPDATE %s
    SET
	    category_id = ?,
	    title = ?,
	    image_url = ?,
	    price = ?,
	    description = ?
	WHERE %s = ?;`, e.TableName(), e.PrimaryKey())

	res, err := r.DB().Exec(ctx, q, e.CategoryID, e.Title, e.ImageURL, e.Price, e.Description, e.ID)
	if err != nil {
		return 0, err
	}

	n := sql.GetAffectedRows(res)
	return n, nil
}

var productUpdateSchema = map[string]reflect.Kind{
	"category_id": reflect.Int,
	"title":       reflect.String,
	"image_url":   reflect.String,
	"price":       reflect.Float32,
	"description": reflect.String,
}

// PartialUpdate accepts a key-value map to
// update the record based on the given "id".
func (repo *ProductRepository) PartialUpdate(ctx context.Context, id int64, attrs map[string]interface{}) (int, error) {
	return repo.Repository.PartialUpdate(ctx, id, productUpdateSchema, attrs)
}
