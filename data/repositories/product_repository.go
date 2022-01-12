package repositories

import (
	"fmt"
	"reflect"
	"strings"

	"morshed/data/engine/sql"
	"morshed/data/models"
	"morshed/domain/repositories"

	"github.com/kataras/iris/v12"
)

// ProductRepository represents the product models service.
// Note that the given models (request) should be already validated
// before service's calls.
type productRepository struct {
	Ctx iris.Context
	*sql.Repository
	rec sql.Record
}

// NewProductRepository returns a new product service to communicate with the database.
func NewProductRepository(db sql.Database) repositories.DataRepository {
	return &productRepository{Repository: sql.NewRepository(db, new(models.Product))}
}

func (r *productRepository) Size(id int64) (int64, error) {
	total, err := r.Count(r.Ctx)
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (r *productRepository) Select(id int64) (prod interface{}, err error) {
	err = r.GetByID(r.Ctx, prod, id)
	return
}

func (r *productRepository) SelectByAttrs(attrs map[string]interface{}) (prod interface{}, err error) {
	err = r.GetByAttrs(r.Ctx, prod, attrs)
	return
}

func (r *productRepository) SelectAll() (prods []interface{}, err error) {
	err = r.GetAll(r.Ctx, prods)
	return
}

func (r *productRepository) Delete(id int64) (int, error) {
	rows, err := r.DeleteByID(r.Ctx, id)
	return rows, err
}

// Insert stores a product to the database and returns its ID.
func (r *productRepository) Insert(p interface{}) (interface{}, error) {
	e := p.(models.Product)
	if !e.ValidateInsert() {
		return models.Product{}, sql.ErrUnprocessable
	}

	q := fmt.Sprintf(`INSERT INTO %s (category_id, title, image_url, price, description)
	VALUES (?,?,?,?,?);`, e.TableName())

	_, err := r.DB().Exec(r.Ctx, q, e.CategoryID, e.Title, e.ImageURL, e.Price, e.Description)
	if err != nil {
		return models.Product{}, err
	}

	return e, nil
}

// BatchInsert inserts one or more products at once and returns the total length created.
func (r *productRepository) BatchInsert(products []interface{}) (int, error) {
	if len(products) == 0 {
		return 0, nil
	}

	var (
		valuesLines []string
		args        []interface{}
	)

	for _, v := range products {
		p := v.(models.Product)
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

	res, err := r.DB().Exec(r.Ctx, q, args...)
	if err != nil {
		return 0, err
	}

	n := sql.GetAffectedRows(res)
	return n, nil
}

// Update updates a product based on its `ID` from the database
// and returns the affected numbrer (0 when nothing changed otherwise 1).
func (r *productRepository) Update(p interface{}) (interface{}, error) {
	e := p.(models.Product)
	q := fmt.Sprintf(`UPDATE %s
    SET
	    category_id = ?,
	    title = ?,
	    image_url = ?,
	    price = ?,
	    description = ?
	WHERE %s = ?;`, e.TableName(), e.PrimaryKey())

	_, err := r.DB().Exec(r.Ctx, q, e.CategoryID, e.Title, e.ImageURL, e.Price, e.Description, e.ID)
	if err != nil {
		return models.Product{}, err
	}

	return e, nil
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
func (r *productRepository) PartialUpdate(id int64, attrs map[string]interface{}) (int, error) {
	return r.Repository.PartialUpdate(r.Ctx, id, productUpdateSchema, attrs)
}
