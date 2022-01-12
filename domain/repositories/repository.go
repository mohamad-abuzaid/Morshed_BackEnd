package repositories

type DataRepository interface {
	Size(int64) (int64, error)
	Select(int64) (interface{}, error)
	SelectByAttrs(map[string]interface{}) (interface{}, error)
	SelectAll() ([]interface{}, error)
	Delete(int64) (int, error)
	Insert(interface{}) (interface{}, error)
	BatchInsert([]interface{}) (int, error)
	Update(interface{}) (interface{}, error)
	PartialUpdate(int64, map[string]interface{}) (int, error)
}
