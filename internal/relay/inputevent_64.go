//go:build amd64 || arm64

package relay

type InputEvent struct {
	Time struct {
		Sec  uint64
		Usec uint64
	}
	Type  uint16
	Code  uint16
	Value int32
}
