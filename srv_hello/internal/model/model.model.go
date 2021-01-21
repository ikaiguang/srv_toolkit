package models

import tkdb "github.com/ikaiguang/srv_toolkit/db"

// db
var (
	TablePrefix string
)

// Setup .
func Setup() {
	TablePrefix = tkdb.TablePrefix()
}
