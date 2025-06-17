package storage

type Storage interface {
	Get(key string) (data any, exists bool)
	GetAll() (data map[string]any)
	Keys() (keys []string)
	Set(key string, data any) (err error)
	Remove(key string) (err error)
	Save() (err error)
}
