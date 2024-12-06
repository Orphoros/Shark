package bytecode

type Version byte

const (
	BcVersionOnos1 Version = iota
)

func (v *Version) String() string {
	switch *v {
	case BcVersionOnos1:
		return "onos1"
	default:
		return "unknown"
	}
}
