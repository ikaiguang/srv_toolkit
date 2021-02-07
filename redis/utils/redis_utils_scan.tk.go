package tkredisutils

import "github.com/gomodule/redigo/redis"

// ScanRedis .
func ScanRedis(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	return redis.Scan(src, dest...)
}

// ScanStruct .
//  var p1, p2 struct {
//		Title  string `redis:"title"`
//		Author string `redis:"author"`
//		Body   string `redis:"body"`
//	}
//  p1.Title = "Example"
//	p1.Author = "Gary"
//	p1.Body = "Hello"
//  c.Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...)
//  m := map[string]string{
//		"title":  "Example2",
//		"author": "Steve",
//		"body":   "Map",
//	}
//  c.Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...)
//  redis.Values(c.Do("HGETALL", id))
//  redis.ScanStruct(v, &p2)
func ScanStruct(src []interface{}, dest interface{}) error {
	return redis.ScanStruct(src, dest)
}

// ScanSlice .
//  c.Send("HMSET", "album:1", "title", "Red", "rating", 5)
//	c.Send("HMSET", "album:2", "title", "Earthbound", "rating", 1)
//	c.Send("HMSET", "album:3", "title", "Beat", "rating", 4)
//	c.Send("LPUSH", "albums", "1")
//	c.Send("LPUSH", "albums", "2")
//	c.Send("LPUSH", "albums", "3")
//	values, err := redis.Values(c.Do("SORT", "albums",
//		"BY", "album:*->rating",
//		"GET", "album:*->title",
//		"GET", "album:*->rating"))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	var albums []struct {
//		Title  string
//		Rating int
//	}
func ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error {
	return redis.ScanSlice(src, dest, fieldNames...)
}
