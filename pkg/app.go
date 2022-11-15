package tinycloud

const Name = "tiny-cloud"

type Ops struct {
	Image  string
	VmType string
}

type App interface {
	Init()
	Run(ops Ops)
	Destroy()
}
