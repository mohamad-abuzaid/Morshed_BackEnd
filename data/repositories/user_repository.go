package repositories

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"morshed/data/engine/sql"
	"morshed/data/models"
	"morshed/domain/repositories"

	"github.com/kataras/iris/v12"
)

type userRepository struct {
	Ctx iris.Context
	*sql.Repository
	rec sql.Record
	mu  sync.RWMutex
}

// NewUserRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserRepository(db sql.Database) repositories.DataRepository {
	return &userRepository{Repository: sql.NewRepository(db, new(models.User))}
}

const (
	// ReadOnlyMode will RLock(read) the data .
	ReadOnlyMode = iota
	// ReadWriteMode will Lock(read/write) the data.
	ReadWriteMode
)

func (r *userRepository) Size(id int64) (int64, error) {
	total, err := r.Count(r.Ctx)
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (r *userRepository) Select(id int64) (user interface{}, err error) {
	err = r.GetByID(r.Ctx, user, id)
	return
}

func (r *userRepository) SelectByAttrs(attrs map[string]interface{}) (user interface{}, err error) {
	err = r.GetByAttrs(r.Ctx, user, attrs)
	return
}

func (r *userRepository) SelectAll() (users []interface{}, err error) {
	err = r.GetAll(r.Ctx, users)
	return
}

func (r *userRepository) Delete(id int64) (int, error) {
	rows, err := r.DeleteByID(r.Ctx, id)
	return rows, err
}

func (r *userRepository) Insert(u interface{}) (interface{}, error) {
	e := u.(models.User)
	if !e.ValidateInsert() {
		return models.User{}, sql.ErrUnprocessable
	}

	q := fmt.Sprintf(`INSERT INTO %s (firstname, username, hashpass)
	VALUES (?,?,?);`, e.TableName())

	_, err := r.DB().Exec(r.Ctx, q, e.Firstname, e.Username, e.HashedPassword)
	if err != nil {
		return models.User{}, err
	}

	return e, nil
}

// BatchInsert inserts one or more users at once and returns the total length created.
func (r *userRepository) BatchInsert(users []interface{}) (int, error) {
	if len(users) == 0 {
		return 0, nil
	}

	var (
		valuesLines []string
		args        []interface{}
	)

	for _, v := range users {
		u := v.(models.User)
		if !u.ValidateInsert() {
			// all users should be "valid", we don't skip, we cancel.
			return 0, sql.ErrUnprocessable
		}

		valuesLines = append(valuesLines, "(?,?,?)")
		args = append(args, []interface{}{u.Firstname, u.Username, u.HashedPassword}...)
	}

	q := fmt.Sprintf("INSERT INTO %s (firstname, username, hashpass) VALUES %s;",
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
func (r *userRepository) Update(u interface{}) (interface{}, error) {
	e := u.(models.User)
	q := fmt.Sprintf(`UPDATE %s
    SET
		firstname = ?,
	    username = ?,
	    hashpass = ?
	WHERE %s = ?;`, e.TableName(), e.PrimaryKey())

	_, err := r.DB().Exec(r.Ctx, q, e.Firstname, e.Username, e.HashedPassword, e.ID)
	if err != nil {
		return models.User{}, err
	}

	return e, nil
}

var userUpdateSchema = map[string]reflect.Kind{
	"firstname": reflect.Int,
	"username":  reflect.String,
	"hashpass":  reflect.String,
}

// PartialUpdate accepts a key-value map to
// update the record based on the given "id".
func (r *userRepository) PartialUpdate(id int64, attrs map[string]interface{}) (int, error) {
	return r.Repository.PartialUpdate(r.Ctx, id, userUpdateSchema, attrs)
}
