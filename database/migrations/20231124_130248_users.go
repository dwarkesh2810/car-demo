package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Users_20231124_130248 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Users_20231124_130248{}
	m.Created = "20231124_130248"

	migration.Register("Users_20231124_130248", m)
}

// Run the migrations
func (m *Users_20231124_130248) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE users(id serial primary key,first_name TEXT NOT NULL,last_name TEXT NOT NULL,email TEXT NOT NULL,mobile TEXT NOT NULL, password TEXT NOT NULL,status integer DEFAULT NULL,role TEXT NOT NULL,otp TEXT NOT NULL)")
}

// Reverse the migrations
func (m *Users_20231124_130248) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE users")
}

