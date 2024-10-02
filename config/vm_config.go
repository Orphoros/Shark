package config

type VmConf struct {
	StackSize   int `json:"stackSize"`
	GlobalsSize int `json:"globalsSize"`
	MaxFrames   int `json:"maxFrames"`
}

func NewDefaultVmConf() VmConf {
	return VmConf{
		StackSize:   2048,
		GlobalsSize: 65536,
		MaxFrames:   1024,
	}
}
