package services

import (
	"errors"
	"morshed/data/models"
	repo "morshed/domain/repositories"

	"github.com/kataras/iris/v12"
)

// UserService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserService interface {
	Count(int64) (int64, error)
	GetByID(int64) (models.User, error)
	GetByAttrs(map[string]interface{}) (models.User, error)
	GetAll() ([]models.User, error)
	DeleteByID(int64) (int, error)
	Create(models.User) (models.User, error)
	InsertAll([]interface{}) (int, error)
	Update(models.User) (models.User, error)
	PatchUpdate(int64, map[string]interface{}) (int, error)
	CreateUser(string, models.User) (models.User, error)
}

// NewUserService returns the default user service.
func NewUserService(repo repo.DataRepository) UserService {
	return &userService{repo: repo}
}

type userService struct {
	Ctx  iris.Context
	repo repo.DataRepository
}

func (s *userService) Count(id int64) (int64, error) {
	total, err := s.repo.Size(id)
	return total, err
}

func (s *userService) GetByID(id int64) (models.User, error) {
	user, err := s.repo.Select(id)
	return user.(models.User), err
}

func (s *userService) GetByAttrs(attrs map[string]interface{}) (models.User, error) {
	user, err := s.repo.SelectByAttrs(attrs)
	return user.(models.User), err
}

func (s *userService) GetAll() ([]models.User, error) {
	us, err := s.repo.SelectAll()
	var prods []models.User
	for _, v := range us {
		prods = append(prods, v.(models.User))
	}
	return prods, err
}

func (s *userService) DeleteByID(id int64) (int, error) {
	row, err := s.repo.Delete(id)
	return row, err
}

func (s *userService) Create(user models.User) (models.User, error) {
	us, err := s.repo.Insert(user)
	return us.(models.User), err
}

func (s *userService) InsertAll(users []interface{}) (int, error) {
	len, err := s.repo.BatchInsert(users)
	return len, err
}

func (s *userService) Update(user models.User) (models.User, error) {
	us, err := s.repo.Update(user)
	return us.(models.User), err
}

func (s *userService) PatchUpdate(id int64, attr map[string]interface{}) (int, error) {
	row, err := s.repo.PartialUpdate(id, attr)
	return row, err
}

func (s *userService) CreateUser(userPassword string, user models.User) (models.User, error) {
	if user.ID > 0 || userPassword == "" || user.Firstname == "" || user.Username == "" {
		return models.User{}, errors.New("unable to create this user")
	}

	hashed, err := models.GeneratePassword(userPassword)
	if err != nil {
		return models.User{}, err
	}
	user.HashedPassword = hashed

	us, err := s.repo.Insert(user)
	return us.(models.User), err
}
