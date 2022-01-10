package services

import (
	"errors"

	"morshed/data/models"
	"morshed/data/repositories"
)

// CategoryService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type CategoryService interface {
	GetAll() []models.User
	GetByID(id int64) (models.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (models.User, bool)
	DeleteByID(id int64) bool

	Update(id int64, user models.User) (models.User, error)
	UpdatePassword(id int64, newPassword string) (models.User, error)
	UpdateUsername(id int64, newUsername string) (models.User, error)

	Create(userPassword string, user models.User) (models.User, error)
}

// NewCategoryService returns the default user service.
func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

type categoryService struct {
	repo repositories.ProductRepository
}

// GetAll returns all users.
func (s *categoryService) GetAll() []models.User {
	return s.repo.SelectMany(func(_ models.User) bool {
		return true
	}, -1)
}

// GetByID returns a user based on its id.
func (s *categoryService) GetByID(id int64) (models.User, bool) {
	return s.repo.Select(func(m models.User) bool {
		return m.ID == id
	})
}

// GetByUsernameAndPassword returns a user based on its username and passowrd,
// used for authentication.
func (s *categoryService) GetByUsernameAndPassword(username, userPassword string) (models.User, bool) {
	if username == "" || userPassword == "" {
		return models.User{}, false
	}

	return s.repo.Select(func(m models.User) bool {
		if m.Username == username {
			hashed := m.HashedPassword
			if ok, _ := models.ValidatePassword(userPassword, hashed); ok {
				return true
			}
		}
		return false
	})
}

// Update updates every field from an existing User,
// it's not safe to be used via public API,
// however we will use it on the web/controllers/user_controller.go#PutBy
// in order to show you how it works.
func (s *categoryService) Update(id int64, user models.User) (models.User, error) {
	user.ID = id
	return s.repo.InsertOrUpdate(user)
}

// UpdatePassword updates a user's password.
func (s *categoryService) UpdatePassword(id int64, newPassword string) (models.User, error) {
	// update the user and return it.
	hashed, err := models.GeneratePassword(newPassword)
	if err != nil {
		return models.User{}, err
	}

	return s.Update(id, models.User{
		HashedPassword: hashed,
	})
}

// UpdateUsername updates a user's username.
func (s *categoryService) UpdateUsername(id int64, newUsername string) (models.User, error) {
	return s.Update(id, models.User{
		Username: newUsername,
	})
}

// Create inserts a new User,
// the userPassword is the client-typed password
// it will be hashed before the insertion to our repository.
func (s *categoryService) Create(userPassword string, user models.User) (models.User, error) {
	if user.ID > 0 || userPassword == "" || user.Firstname == "" || user.Username == "" {
		return models.User{}, errors.New("unable to create this user")
	}

	hashed, err := models.GeneratePassword(userPassword)
	if err != nil {
		return models.User{}, err
	}
	user.HashedPassword = hashed

	return s.repo.InsertOrUpdate(user)
}

// DeleteByID deletes a user by its id.
//
// Returns true if deleted otherwise false.
func (s *categoryService) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m models.User) bool {
		return m.ID == id
	}, 1)
}
