package entity

type SchemaMigration struct {
	Version int64
	Dirty   bool
}

type SchemaMigrations []SchemaMigration
