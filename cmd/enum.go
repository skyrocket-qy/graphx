package cmd

import errors "github.com/rotisserie/eris"

type DatabaseEnum string

const (
	databaseEnumPg     DatabaseEnum = "pg"
	databaseEnumSqlite DatabaseEnum = "sqlite"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *DatabaseEnum) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *DatabaseEnum) Set(v string) error {
	switch v {
	case string(databaseEnumPg), string(databaseEnumSqlite):
		*e = DatabaseEnum(v)
		return nil
	default:
		return errors.New(`must be one of "pg", "sqlite"`)
	}
}

// Type is only used in help text
func (e *DatabaseEnum) Type() string {
	return "DatabaseEnum"
}
