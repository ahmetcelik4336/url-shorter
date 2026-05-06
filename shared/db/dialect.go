package db

import "fmt"

type DialectHandler interface {
	FormatDate(column string) string
}

type MySQLHandler struct{}

func (MySQLHandler) FormatDate(c string) string {
	return fmt.Sprintf("DATE_FORMAT(%s, '%%Y-%%m-%%d')", c)
}

type PostgresHandler struct{}

func (PostgresHandler) FormatDate(c string) string {
	return fmt.Sprintf("TO_CHAR(%s, 'YYYY-MM-DD')", c)
}

type SQLiteHandler struct{}

func (SQLiteHandler) FormatDate(c string) string {
	return fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", c)
}

func GetDialect(driverName string) DialectHandler {
	switch driverName {
	case "postgres":
		return PostgresHandler{}
	case "mysql":
		return MySQLHandler{}
	case "sqlite3":
		return SQLiteHandler{}
	default:
		return MySQLHandler{}
	}
}
