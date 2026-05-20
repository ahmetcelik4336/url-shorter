package db

import "fmt"

type DialectHandler interface {
	FormatDate(column string) string
	GroupByDate(column string) string
}

type MySQLHandler struct{}

func (MySQLHandler) FormatDate(c string) string {
	return fmt.Sprintf("DATE_FORMAT(%s, '%%Y-%%m-%%d')", c)
}

func (MySQLHandler) GroupByDate(column string) string {
	return fmt.Sprintf("DATE(%s)", column)
}

type PostgresHandler struct{}

func (PostgresHandler) FormatDate(c string) string {
	return fmt.Sprintf("CAST(%s AS DATE)", c)
}

func (PostgresHandler) GroupByDate(column string) string {
	return fmt.Sprintf("%s::date", column)
}

type SQLiteHandler struct{}

func (SQLiteHandler) FormatDate(c string) string {
	return fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", c)
}
func (SQLiteHandler) GroupByDate(column string) string {
	return fmt.Sprintf("DATE(%s)", column)
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
