package tinycloud

type Ops struct {
	VmType string
}

type App interface {
	Run(ops Ops) error
	Destroy() error
}
