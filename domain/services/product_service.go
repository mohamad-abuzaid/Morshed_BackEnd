package services

import (
	"errors"

	"morshed/data/models"
	"morshed/data/repositories"
)

// ProductService handles CRUID operations of a product datamodel,
// it depends on a product repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type ProductService interface {
	GetAll() []models.Product
	GetByID(id int64) (models.Product, bool)
	DeleteByID(id int64) bool
	Update(id int64, product models.Product) (models.Product, error)
	Create(product models.Product) (models.Product, error)
}

// NewProductService returns the default product service.
func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

type productService struct {
	repo repositories.ProductRepository
}

// GetAll returns all users.
func (s *productService) GetAll() []models.Product {
	return s.repo.SelectMany(func(_ models.Product) bool {
		return true
	}, -1)
}

// GetByID returns a product based on its id.
func (s *productService) GetByID(id int64) (models.Product, bool) {
	return s.repo.Select(func(m models.Product) bool {
		return m.ID == id
	})
}

// GetByUsernameAndPassword returns a product based on its username and passowrd,
// used for authentication.
func (s *productService) GetByUsernameAndPassword(username, userPassword string) (models.Product, bool) {
	if username == "" || userPassword == "" {
		return models.Product{}, false
	}

	return s.repo.Select(func(m models.Product) bool {
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
func (s *productService) Update(id int64, product models.Product) (models.Product, error) {
	product.ID = id
	return s.repo.InsertOrUpdate(product)
}

// UpdatePassword updates a product's password.
func (s *productService) UpdatePassword(id int64, newPassword string) (models.Product, error) {
	// update the product and return it.
	hashed, err := models.GeneratePassword(newPassword)
	if err != nil {
		return models.Product{}, err
	}

	return s.Update(id, models.Product{
		HashedPassword: hashed,
	})
}

// UpdateUsername updates a product's username.
func (s *productService) UpdateUsername(id int64, newUsername string) (models.Product, error) {
	return s.Update(id, models.Product{
		Username: newUsername,
	})
}

// Create inserts a new User,
// the userPassword is the client-typed password
// it will be hashed before the insertion to our repository.
func (s *productService) Create(userPassword string, product models.Product) (models.Product, error) {
	if product.ID > 0 || userPassword == "" || product.Firstname == "" || product.Username == "" {
		return models.Product{}, errors.New("unable to create this product")
	}

	hashed, err := models.GeneratePassword(userPassword)
	if err != nil {
		return models.Product{}, err
	}
	product.HashedPassword = hashed

	return s.repo.InsertOrUpdate(product)
}

// DeleteByID deletes a product by its id.
//
// Returns true if deleted otherwise false.
func (s *productService) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m models.Product) bool {
		return m.ID == id
	}, 1)
}
