package meta

type Env interface {
	Apply(alias Alias) error
}

func New(global bool) Env {
	return instance(global)
}
