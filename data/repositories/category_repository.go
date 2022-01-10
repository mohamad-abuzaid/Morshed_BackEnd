package repositories

import (
	"context"
	"fmt"
	"reflect"

	"morshed/data/engine/sql"
	"morshed/data/models"
)

// CategoryRepository represents the category models service.
// Note that the given models (request) should be already validated
// before service's calls.
type CategoryRepository struct {
	*sql.Repository
}

// NewCategoryRepository returns a new category service to communicate with the database.
func NewCategoryRepository(db sql.Database) *CategoryRepository {
	return &CategoryRepository{Repository: sql.NewRepository(db, new(models.Category))}
}

// Insert stores a category to the database and returns its ID.
func (r *CategoryRepository) Insert(ctx context.Context, e models.Category) (int64, error) {
	if e.Title == "" || e.ImageURL == "" {
		return 0, sql.ErrUnprocessable
	}

	q := fmt.Sprintf(`INSERT INTO %s (title, position, image_url)
	VALUES (?,?,?);`, e.TableName())

	res, err := r.DB().Exec(ctx, q, e.Title, e.Position, e.ImageURL)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Update updates a category based on its `ID`.
func (r *CategoryRepository) Update(ctx context.Context, e models.Category) (int, error) {
	if e.ID == 0 || e.Title == "" || e.ImageURL == "" {
		return 0, sql.ErrUnprocessable
	}

	q := fmt.Sprintf(`UPDATE %s
    SET
	    title = ?,
	    position = ?,
	    image_url = ?
	WHERE %s = ?;`, e.TableName(), e.PrimaryKey())

	res, err := r.DB().Exec(ctx, q, e.Title, e.Position, e.ImageURL, e.ID)
	if err != nil {
		return 0, err
	}

	n := sql.GetAffectedRows(res)
	return n, nil
}

// The updatable fields, separately from that we create for any possible future necessities.
var categoryUpdateSchema = map[string]reflect.Kind{
	"title":     reflect.String,
	"image_url": reflect.String,
	"position":  reflect.Int,
}

// PartialUpdate accepts a key-value map to
// update the record based on the given "id".
func (r *CategoryRepository) PartialUpdate(ctx context.Context, id int64, attrs map[string]interface{}) (int, error) {
	return r.Repository.PartialUpdate(ctx, id, categoryUpdateSchema, attrs)
}
