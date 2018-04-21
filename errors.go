package blockchaindb

type BlockError struct {
	err string
}

func NewBlockError(message string) BlockError {
	return BlockError{err: message}
}

func (e BlockError) Error() string {
	return e.err
}
