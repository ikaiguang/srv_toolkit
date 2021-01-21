package testmodels

import models "github.com/ikaiguang/srv_toolkit/srv_hello/internal/model"

// Test .
type Test struct{}

// TableName : table name
func (m *Test) TableName() string {
	return models.TablePrefix + "test"
}

// test Test
type test struct {
	Test
}

// TestModel test
var TestModel = new(test)
