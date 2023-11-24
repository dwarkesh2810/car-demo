package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Users_20231124_145928 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Users_20231124_145928{}
	m.Created = "20231124_145928"

	migration.Register("Users_20231124_145928", m)
}

// Run the migrations
func (m *Users_20231124_145928) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE users()")
}

// Reverse the migrations
func (m *Users_20231124_145928) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE users")
}
