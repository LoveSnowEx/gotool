package errors

type Error uint64

const (
	ErrInvalidField Error = 1 << iota
	ErrUnknownDatabaseDriver
)

var errorStrings = map[Error]string{
	ErrInvalidField:          "invalid field",
	ErrUnknownDatabaseDriver: "unknown database driver",
}

func (e Error) Error() string {
	if s, ok := errorStrings[e]; ok {
		return s
	}
	return "unknown error"
}
