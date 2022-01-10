package repositories

import (
	"errors"
	"sync"

	"morshed/data/engine/sql"
	"morshed/data/models"
)

// Query represents the visitor and action queries.
type Query func(models.User) bool

// NewUserRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserRepository(source sql.Database) DataRepository {
	return &userMysqlRepository{source: source}
}

// userMysqlRepository is a "DataRepository"
// which manages the users using the memory data source (map).
type userMysqlRepository struct {
	source sql.Database
	mu     sync.RWMutex
}

const (
	// ReadOnlyMode will RLock(read) the data .
	ReadOnlyMode = iota
	// ReadWriteMode will Lock(read/write) the data.
	ReadWriteMode
)

func (r *userMysqlRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0

	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	for _, user := range r.source {
		ok = query(user)
		if ok {
			if action(user) {
				loops++
				if actionLimit >= loops {
					break // break
				}
			}
		}
	}

	return
}

// Select receives a query function
// which is fired for every single user model inside
// our imaginary data source.
// When that function returns true then it stops the iteration.
//
// It returns the query's return last known boolean value
// and the last known user model
// to help callers to reduce the LOC.
//
// It's actually a simple but very clever prototype function
// I'm using everywhere since I firstly think of it,
// hope you'll find it very useful as well.
func (r *userMysqlRepository) Select(query Query) (user models.User, found bool) {
	found = r.Exec(query, func(m models.User) bool {
		user = m
		return true
	}, 1, ReadOnlyMode)

	// set an empty models.User if not found at all.
	if !found {
		user = models.User{}
	}

	return
}

// SelectMany same as Select but returns one or more models.User as a slice.
// If limit <=0 then it returns everything.
func (r *userMysqlRepository) SelectMany(query Query, limit int) (results []models.User) {
	r.Exec(query, func(m models.User) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)

	return
}

// InsertOrUpdate adds or updates a user to the (memory) storage.
//
// Returns the new user and an error if any.
func (r *userMysqlRepository) InsertOrUpdate(user models.User) (models.User, error) {
	id := user.ID

	if id == 0 { // Create new action
		var lastID int64
		// find the biggest ID in order to not have duplications
		// in productions apps you can use a third-party
		// library to generate a UUID as string.
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()

		id = lastID + 1
		user.ID = id

		// map-specific thing
		r.mu.Lock()
		r.source[id] = user
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
	r.source[id] = current
	r.mu.Unlock()

	return user, nil
}

func (r *userMysqlRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(m models.User) bool {
		delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)
}
