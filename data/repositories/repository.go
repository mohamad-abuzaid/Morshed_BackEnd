package repositories

// DataRepository handles the basic operations of an entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type DataRepository interface {
	// Select() (user models.User, found bool)
	// SelectMany() (results []models.User)

	// InsertOrUpdate() (updatedUser models.User, err error)
	// Delete() (deleted bool)
}