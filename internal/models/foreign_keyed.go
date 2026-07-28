package models

type ForeignKeyed interface {
	ForeignKeys() []string
}
