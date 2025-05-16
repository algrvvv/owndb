package storage

type Marshaller interface {
	Marshal(map[string]any) ([]byte, error)
	Unmarshal(data []byte) (map[string]any, error)
}
