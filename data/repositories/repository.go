package repositories

import "morshed/data/models"

// DataRepository handles the basic operations of an entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type DataRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)

	Select(query Query) (user models.User, found bool)
	SelectMany(query Query, limit int) (results []models.User)

	InsertOrUpdate(user models.User) (updatedUser models.User, err error)
	Delete(query Query, limit int) (deleted bool)
}