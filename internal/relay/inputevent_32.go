//go:build 386 || arm

package relay

type InputEvent struct {
	Time struct {
		Sec  uint32
		Usec uint32
	}
	Type  uint16
	Code  uint16
	Value int32
}
