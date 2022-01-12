package services

import (
	"morshed/data/models"
	repo "morshed/domain/repositories"

	"github.com/kataras/iris/v12"
)

// ProductService handles CRUID operations of a product datamodel,
// it depends on a product repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type ProductService interface {
	Count(int64) (int64, error)
	GetByID(int64) (models.Product, error)
	GetByAttrs(map[string]interface{}) (models.Product, error)
	GetAll() ([]models.Product, error)
	DeleteByID(int64) (int, error)
	Create(models.Product) (models.Product, error)
	InsertAll([]interface{}) (int, error)
	Update(models.Product) (models.Product, error)
	PatchUpdate(int64, map[string]interface{}) (int, error)
}

// NewProductService returns the default product service.
func NewProductService(repo repo.DataRepository) ProductService {
	return &productService{repo: repo}
}

type productService struct {
	Ctx  iris.Context
	repo repo.DataRepository
}

func (s *productService) Count(id int64) (int64, error) {
	total, err := s.repo.Size(id)
	return total, err
}

func (s *productService) GetByID(id int64) (models.Product, error) {
	prod, err := s.repo.Select(id)
	return prod.(models.Product), err
}

func (s *productService) GetByAttrs(attrs map[string]interface{}) (models.Product, error) {
	prod, err := s.repo.SelectByAttrs(attrs)
	return prod.(models.Product), err
}

func (s *productService) GetAll() ([]models.Product, error) {
	ps, err := s.repo.SelectAll()
	var prods []models.Product
	for _, v := range ps {
		prods = append(prods, v.(models.Product))
	}
	return prods, err
}

func (s *productService) DeleteByID(id int64) (int, error) {
	row, err := s.repo.Delete(id)
	return row, err
}

func (s *productService) Create(product models.Product) (models.Product, error) {
	prod, err := s.repo.Insert(product)
	return prod.(models.Product), err
}

func (s *productService) InsertAll(products []interface{}) (int, error) {
	len, err := s.repo.BatchInsert(products)
	return len, err
}

func (s *productService) Update(product models.Product) (models.Product, error) {
	prod, err := s.repo.Update(product)
	return prod.(models.Product), err
}

func (s *productService) PatchUpdate(id int64, attr map[string]interface{}) (int, error) {
	row, err := s.repo.PartialUpdate(id, attr)
	return row, err
}
