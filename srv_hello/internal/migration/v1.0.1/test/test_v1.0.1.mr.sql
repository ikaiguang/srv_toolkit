-- 测试表
CREATE TABLE IF NOT EXISTS srv_test
(
    id            BIGINT AUTO_INCREMENT COMMENT 'id',
    column_string VARCHAR(255) NOT NULL DEFAULT '' COMMENT '字符串',
    column_int    INT          NOT NULL DEFAULT '0' COMMENT '数字',
    PRIMARY KEY (id)
) ENGINE InnoDB
  DEFAULT CHARSET utf8mb4
    COMMENT '测试表';