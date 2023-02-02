package tinycloud

const Name = "tiny-cloud"

type Ops struct {
	VmType string
}

type App interface {
	Run(ops Ops) error
	Destroy() error
}
