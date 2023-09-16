package dao

import "testing"

func setup() {
	Init()
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}
