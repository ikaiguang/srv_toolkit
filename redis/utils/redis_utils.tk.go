package tkredisutils

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/env"
	tkredis "github.com/ikaiguang/srv_toolkit/redis"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// key
var (
	_keyPrefix string
	_keyOnce   sync.Once
)

// usePrecise .
func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

// formatMs .
func formatMs(dur time.Duration) int64 {
	//if dur > 0 && dur < time.Millisecond {
	//	fmt.Printf(
	//		"specified duration is %s, but minimal supported value is %s",
	//		dur, time.Millisecond,
	//	)
	//}
	return int64(dur / time.Millisecond)
}

// formatSec .
func formatSec(dur time.Duration) int64 {
	//if dur > 0 && dur < time.Second {
	//	fmt.Printf(
	//		"specified duration is %s, but minimal supported value is %s",
	//		dur, time.Second,
	//	)
	//}
	return int64(dur / time.Second)
}

// KeyPrefix .
func KeyPrefix() string {
	_keyOnce.Do(func() {
		_keyPrefix = env.AppID
	})
	return _keyPrefix
}

// Key .
func Key(key string) string {
	return KeyPrefix() + ":" + key
}

//------------------------------------------------------------------------------

// Command .
// 命令用于返回所有的Redis命令的详细信息，以数组形式展示。
func Command(ctx context.Context) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "command")
	return
}

// ClientGetName .
// 命令用于返回 CLIENT SETNAME 命令为连接设置的名字。 因为新创建的连接默认是没有名字的， 对于没有名字的连接， CLIENT GETNAME 返回空白回复。
func ClientGetName(ctx context.Context) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "client", "getname")
	return
}

// Echo .
// Redis Echo 命令用于打印给定的字符串。
func Echo(ctx context.Context, message interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "echo", message)
	return
}

// Ping .
// Redis Ping 命令使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
func Ping(ctx context.Context) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "ping")
	return
}

// Quit .
// not implemented
func Quit(ctx context.Context) (reply interface{}, err error) {
	//reply, err = tkredis.Redis().Do(ctx, "quit")
	err = errors.New("not implemented")
	return
}

// Del .
// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
func Del(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "del", keys...)
	return
}

// Unlink .
// 该命令和DEL十分相似：删除指定的key(s),若key不存在则该key被跳过。
// 但是，相比DEL会产生阻塞，该命令会在另一个线程中回收内存，因此它是非阻塞的。
// 这也是该命令名字的由来：仅将keys从keyspace元数据中删除，真正的删除会在后续异步操作。
func Unlink(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "unlink", keys...)
	return
}

// Dump .
// Redis DUMP 命令用于序列化给定 key ，并返回被序列化的值。
func Dump(ctx context.Context, key interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "dump", key)
	return
}

// Exists .
// Redis EXISTS 命令用于检查给定 key 是否存在。
func Exists(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "exists", keys...)
	return
}

// Expire .
// Redis Expire 命令用于设置 key 的过期时间，key 过期后将不再可用。单位以秒计。
func Expire(ctx context.Context, key string, expiration time.Duration) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "expire", key, formatSec(expiration))
	return
}

// ExpireAt .
// Redis Expireat 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间。key 过期后将不再可用。
func ExpireAt(ctx context.Context, key string, tm time.Time) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "expireat", key, tm.Unix())
	return
}

// Keys .
// Redis Keys 命令用于查找所有符合给定模式 pattern 的 key 。。
func Keys(ctx context.Context, pattern string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "keys", pattern)
	return
}

// Migrate .
// 将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
// 这个命令是一个原子操作，它在执行的时候会阻塞进行迁移的两个实例，直到以下任意结果发生：迁移成功，迁移失败，等到超时。
func Migrate(ctx context.Context, host, port, key string, db int64, timeout time.Duration) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "migrate",
		host,
		port,
		key,
		db,
		formatMs(timeout),
	)
	return
}

// Move .
// Redis MOVE 命令用于将当前数据库的 key 移动到给定的数据库 db 当中。
func Move(ctx context.Context, key string, db int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "move", key, db)
	return
}

// ObjectRefCount .
// OBJECT REFCOUNT 该命令主要用于调试(debugging)，它能够返回指定key所对应value被引用的次数.
// OBJECT REFCOUNT <key> 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
func ObjectRefCount(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "object", "refcount", key)
	return
}

// ObjectEncoding .
// OBJECT ENCODING <key> 返回给定 key 锁储存的值所使用的内部表示(representation)。
// 对象可以以多种方式编码：
// 字符串可以被编码为 raw (一般字符串)或 int (用字符串表示64位数字是为了节约空间)。
// 列表可以被编码为 ziplist 或 linkedlist 。 ziplist 是为节约大小较小的列表空间而作的特殊表示。
// 集合可以被编码为 intset 或者 hashtable 。 intset 是只储存数字的小集合的特殊表示。
// 哈希表可以编码为 zipmap 或者 hashtable 。 zipmap 是小哈希表的特殊表示。
// 有序集合可以被编码为 ziplist 或者 skiplist 格式。 ziplist 用于表示小的有序集合，而 skiplist 则用于表示任何大小的有序集合。
func ObjectEncoding(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "object", "encoding", key)
	return
}

// ObjectIdleTime .
// OBJECT IDLETIME <key> 返回给定 key 自储存以来的空转时间(idle， 没有被读取也没有被写入)，以秒为单位。
func ObjectIdleTime(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "object", "idletime", key)
	return
}

// Persist .
// Redis PERSIST 命令用于移除给定 key 的过期时间，使得 key 永不过期。
func Persist(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "persist", key)
	return
}

// PExpire .
// Redis PEXPIRE 命令和 EXPIRE 命令的作用类似，但是它以毫秒为单位设置 key 的生存时间，而不像 EXPIRE 命令那样，以秒为单位。
func PExpire(ctx context.Context, key string, expiration time.Duration) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "pexpire", key, formatMs(expiration))
	return
}

// PExpireAt .
// Redis PEXPIREAT 命令用于设置 key 的过期时间，以毫秒计。key 过期后将不再可用。
func PExpireAt(ctx context.Context, key string, tm time.Time) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "pexpireat",
		key,
		tm.UnixNano()/int64(time.Millisecond),
	)
	return
}

// PTTL .
// Redis Pttl 命令以毫秒为单位返回 key 的剩余过期时间。
func PTTL(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "pttl", key)
	return
}

// RandomKey .
// Redis RANDOMKEY 命令从当前数据库中随机返回一个 key 。
func RandomKey(ctx context.Context) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "randomkey")
	return
}

// Rename .
// Redis Rename 命令用于修改 key 的名称 。
func Rename(ctx context.Context, key, newkey string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "rename", key, newkey)
	return
}

// RenameNX .
// Redis Renamenx 命令用于在新的 key 不存在时修改 key 的名称 。
func RenameNX(ctx context.Context, key, newkey string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "renamenx", key, newkey)
	return
}

// Restore .
// 反序列化给定的序列化值，并将它和给定的 key 关联。
// 那么使用反序列化得出的值来代替键 key 原有的值； 相反地， 如果键 key 已经存在， 但是没有给定 REPLACE 选项， 那么命令返回一个错误。
func Restore(ctx context.Context, key string, ttl time.Duration, value string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "restore", key, formatMs(ttl), value)
	return
}

// RestoreReplace .
// 反序列化给定的序列化值，并将它和给定的 key 关联。
// 那么使用反序列化得出的值来代替键 key 原有的值； 相反地， 如果键 key 已经存在， 但是没有给定 REPLACE 选项， 那么命令返回一个错误。
func RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "restore",
		key,
		formatMs(ttl),
		value,
		"replace",
	)
	return
}

// SortParam .
type SortParam struct {
	By            string
	Offset, Count int64
	Get           []string
	Order         string
	Alpha         bool
}

// args .
func (sort *SortParam) args(key string) []interface{} {
	//args := []interface{}{"sort", key}
	args := []interface{}{key}
	if sort.By != "" {
		args = append(args, "by", sort.By)
	}
	if sort.Offset != 0 || sort.Count != 0 {
		args = append(args, "limit", sort.Offset, sort.Count)
	}
	for _, get := range sort.Get {
		args = append(args, "get", get)
	}
	if sort.Order != "" {
		args = append(args, sort.Order)
	}
	if sort.Alpha {
		args = append(args, "alpha")
	}
	return args
}

// Sort .
// redis支持对list，set，sorted set、hash元素（元素可以为数值与字符串）的排序。
// sort key [BY pattern] [LIMIT start count] [GET pattern] [ASC|DESC] [ALPHA] [STORE dstkey]
func Sort(ctx context.Context, key string, sort *SortParam) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "sort", sort.args(key)...)
	return
}

// SortStore .
// redis支持对list，set，sorted set、hash元素（元素可以为数值与字符串）的排序。
// sort key [BY pattern] [LIMIT start count] [GET pattern] [ASC|DESC] [ALPHA] [STORE dstkey]
// 如果对集合经常按照固定的模式去排序，那么把排序结果缓存起来会减少不少cpu开销，使用store选项可以将排序内容保存到指定key中，保存的类型是list
func SortStore(ctx context.Context, key, store string, sort *SortParam) (reply interface{}, err error) {
	args := sort.args(key)
	if store != "" {
		args = append(args, "store", store)
	}
	reply, err = tkredis.Redis().Do(ctx, "sort", args...)
	return
}

// SortInterfaces .
// redis支持对list，set，sorted set、hash元素（元素可以为数值与字符串）的排序。
// sort key [BY pattern] [LIMIT start count] [GET pattern] [ASC|DESC] [ALPHA] [STORE dstkey]
func SortInterfaces(ctx context.Context, key string, sort *SortParam) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "sort", sort.args(key)...)
	return
}

// Touch .
// 修改指定 key 的 最后访问时间。忽略不存在的 key。
func Touch(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "touch", keys...)
	return
}

// TTL .
// Redis TTL 命令以秒为单位返回 key 的剩余过期时间。
func TTL(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "ttl", key)
	return
}

// Type .
// Redis Type 命令用于返回 key 所储存的值的类型。
func Type(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "type", key)
	return
}

// Scan .
// SCAN 命令是一个基于游标的迭代器，每次被调用之后， 都会向用户返回一个新的游标， 用户在下次迭代时需要使用这个新游标作为 SCAN 命令的游标参数， 以此来延续之前的迭代过程。
// SCAN 返回一个包含两个元素的数组， 第一个元素是用于进行下一次迭代的新游标， 而第二个元素则是一个数组， 这个数组中包含了所有被迭代的元素。如果新游标返回 0 表示迭代已结束。
func Scan(ctx context.Context, cursor uint64, match string, count int64) (reply interface{}, err error) {
	args := []interface{}{cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	reply, err = tkredis.Redis().Do(ctx, "scan", args...)
	return
}

// SScan .
// Redis Sscan 命令用于迭代集合中键的元素，Sscan 继承自 Scan。
func SScan(ctx context.Context, key string, cursor uint64, match string, count int64) (reply interface{}, err error) {
	args := []interface{}{key, cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	reply, err = tkredis.Redis().Do(ctx, "sscan", args...)
	return
}

// HScan .
// Redis HSCAN 命令用于迭代哈希表中的键值对。
func HScan(ctx context.Context, key string, cursor uint64, match string, count int64) (reply interface{}, err error) {
	args := []interface{}{key, cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	reply, err = tkredis.Redis().Do(ctx, "hscan", args...)
	return
}

// ZScan .
// Redis Zscan 命令用于迭代有序集合中的元素（包括元素成员和元素分值）
func ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) (reply interface{}, err error) {
	args := []interface{}{key, cursor}
	if match != "" {
		args = append(args, "match", match)
	}
	if count > 0 {
		args = append(args, "count", count)
	}
	reply, err = tkredis.Redis().Do(ctx, "zscan", args...)
	return
}

//------------------------------------------------------------------------------

// Append .
// Redis Append 命令用于为指定的 key 追加值。
// 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func Append(ctx context.Context, key, value string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "append", key, value)
	return
}
