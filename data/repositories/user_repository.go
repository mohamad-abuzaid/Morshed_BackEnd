package repositories

import (
	"errors"
	"sync"

	"morshed/data/engine/sql"
	"morshed/data/models"

	"github.com/kataras/iris/v12"
)

type UserRepository struct {
	Ctx iris.Context
	*sql.Repository
	rec sql.Record
	mu     sync.RWMutex
}

// NewUserRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserRepository(db sql.Database) *UserRepository {
	return &UserRepository{Repository: sql.NewRepository(db, new(models.User))}
}

const (
	// ReadOnlyMode will RLock(read) the data .
	ReadOnlyMode = iota
	// ReadWriteMode will Lock(read/write) the data.
	ReadWriteMode
)

func (r *UserRepository) Size(id int64) (int64, error) {
	total, err := r.Count(r.Ctx)
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (r *UserRepository) Select(id int64) (user models.User, err error) {
	err = r.GetByID(r.Ctx, user, id)
	return
}

// SelectMany same as Select but returns one or more models.User as a slice.
// If limit <=0 then it returns everything.
func (r *UserRepository) SelectMany(query Query, limit int) (results []models.User) {
	r.Exec(query, func(m models.User) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)

	return
}

// InsertOrUpdate adds or updates a user to the (memory) storage.
//
// Returns the new user and an error if any.
func (r *UserRepository) InsertOrUpdate(user models.User) (models.User, error) {
	id := user.ID

	if id == 0 { // Create new action
		var lastID int64
		// find the biggest ID in order to not have duplications
		// in productions apps you can use a third-party
		// library to generate a UUID as string.
		r.mu.RLock()
		// for _, item := range r.source {
		// 	if item.ID > lastID {
		// 		lastID = item.ID
		// 	}
		// }
		r.mu.RUnlock()

		id = lastID + 1
		user.ID = id

		// map-specific thing
		r.mu.Lock()
		// r.source[id] = user
		r.mu.Unlock()

		return user, nil
	}

	// Update action based on the user.ID,
	// here we will allow updating the poster and genre if not empty.
	// Alternatively we could do pure replace instead:
	// r.source[id] = user
	// and comment the code below;
	current, exists := r.Select(func(m models.User) bool {
		return m.ID == id
	})

	if !exists { // ID is not a real one, return an error.
		return models.User{}, errors.New("failed to update a nonexistent user")
	}

	// or comment these and r.source[id] = user for pure replace
	if user.Username != "" {
		current.Username = user.Username
	}

	if user.Firstname != "" {
		current.Firstname = user.Firstname
	}

	// map-specific thing
	r.mu.Lock()
	// r.source[id] = current
	r.mu.Unlock()

	return user, nil
}

func (r *UserRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(m models.User) bool {
		// delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)
}
