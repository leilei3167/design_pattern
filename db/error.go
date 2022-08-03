package db

type ErrDB int8

func (e ErrDB) Error() string {
	return errDBMsg[e]
}

const (
	ErrPrimaryKeyConflict ErrDB = 100 + iota
	ErrRecordNotFound
	ErrRecordTypeInvalid
	ErrTableNotExist
	ErrTableAlreadyExist
	ErrTransactionNotBegin
	ErrDslInvalidGrammar
)

var errDBMsg = map[ErrDB]string{
	ErrPrimaryKeyConflict:  "primary key conflict",
	ErrRecordNotFound:      "model not found",
	ErrRecordTypeInvalid:   "model type invalid",
	ErrTableNotExist:       "table not exist",
	ErrTableAlreadyExist:   "table already exist",
	ErrTransactionNotBegin: "transaction not begin",
	ErrDslInvalidGrammar:   "dsl expression invalid grammar",
}
