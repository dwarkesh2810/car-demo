package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CarMaster_20231123_155827 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CarMaster_20231123_155827{}
	m.Created = "20231123_155827"

	migration.Register("CarMaster_20231123_155827", m)
}

// Run the migrations
func (m *CarMaster_20231123_155827) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	// m.SQL("CREATE TYPE car_types AS ENUM ('sedan', 'hatchback', 'SUV')")
	m.SQL("CREATE TABLE car_master(id serial primary key,user_id integer DEFAULT NULL,car_name TEXT NOT NULL,car_image TEXT NOT NULL,make TEXT NOT NULL,model TEXT NOT NULL,car_type car_types NOT NULL,created_at integer DEFAULT NULL, updated_at integer DEFAULT NULL)")
}

// Reverse the migrations
func (m *CarMaster_20231123_155827) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE car_master")
}
