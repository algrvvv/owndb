package exec

type Executor interface {
	Execute(command string) (res any, err error)
}
