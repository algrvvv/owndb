package snapshot

import (
	"io"
	"os"

	"github.com/algrvvv/owndb/internal/storage"
)

type SnapshotManager struct {
	marshaller storage.Marshaller
	src        *os.File
}

func NewSnapshotManager(marshaller storage.Marshaller, src *os.File) *SnapshotManager {
	return &SnapshotManager{
		marshaller: marshaller,
		src:        src,
	}
}

func (s *SnapshotManager) Write(m map[string]any) (int, error) {
	data, err := s.marshaller.Marshal(m)
	if err != nil {
		panic(err)
	}

	// чистим старые данные
	err = s.src.Truncate(0)
	if err != nil {
		return 0, err
	}

	_, err = s.src.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	return s.src.Write(data)
}

func (s *SnapshotManager) Read() (map[string]any, error) {
	data, err := io.ReadAll(s.src)
	if err != nil {
		return nil, err
	}

	return s.marshaller.Unmarshal(data)
}
