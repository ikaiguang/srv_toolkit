package testmodels

import (
	"fmt"
	models "github.com/ikaiguang/srv_toolkit/srv_hello/internal/model"
)

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

// TableSQL : table SQL
func (m *test) TableSQL() string {
	execSQL := `
CREATE TABLE %s
(
    id            BIGINT AUTO_INCREMENT COMMENT 'id',
    column_string VARCHAR(255) NOT NULL DEFAULT '' COMMENT '字符串',
    column_int    INT          NOT NULL DEFAULT '0' COMMENT '数字',
    PRIMARY KEY (id)
) ENGINE InnoDB
  DEFAULT CHARSET utf8mb4
    COMMENT '测试表';
`
	return fmt.Sprintf(execSQL, m.TableName())
}
