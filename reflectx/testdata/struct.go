package testdata

// Ao Ao
type Ao struct {
	Base
	M1 M
	M2 M
}

// M M
type M struct {
	Base
}

// Base Base
type Base struct {
	db *DB
}

// DB DB
type DB struct {
	Name string
}
