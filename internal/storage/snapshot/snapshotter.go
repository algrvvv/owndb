package snapshot

type Snapshotter interface {
	Write(m map[string]any) (int, error)
	Read() (map[string]any, error)
}
