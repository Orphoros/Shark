package bytecode

type Type byte

var minType Type = 0
var maxType Type = 1

const (
	BcTypeNormal Type = iota
	BcTypeCompressedBrotli
)

func (c *Type) String() string {
	switch *c {
	case BcTypeNormal:
		return "none"
	case BcTypeCompressedBrotli:
		return "brotli"
	default:
		return "unknown"
	}
}
