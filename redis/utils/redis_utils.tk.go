package tkredisutils

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/pkg/errors"
	"io"
	"sync"
	"time"
)

// key
var (
	_client    *redis.Redis
	_keyPrefix string
	_keyOnce   sync.Once
)

// SetConn .
func SetConn(pool *redis.Redis) {
	_client = pool
}

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
	reply, err = _client.Do(ctx, "command")
	return
}

// ClientGetName .
// 命令用于返回 CLIENT SETNAME 命令为连接设置的名字。 因为新创建的连接默认是没有名字的， 对于没有名字的连接， CLIENT GETNAME 返回空白回复。
func ClientGetName(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "getname")
	return
}

// Echo .
// Redis Echo 命令用于打印给定的字符串。
func Echo(ctx context.Context, message interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "echo", message)
	return
}

// Ping .
// Redis Ping 命令使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
func Ping(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "ping")
	return
}

// Quit .
// not implemented
func Quit(ctx context.Context) (reply interface{}, err error) {
	//reply, err = _client.Do(ctx, "quit")
	err = errors.New("not implemented")
	return
}

// Del .
// Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
func Del(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "del", keys...)
	return
}

// Unlink .
// 该命令和DEL十分相似：删除指定的key(s),若key不存在则该key被跳过。
// 但是，相比DEL会产生阻塞，该命令会在另一个线程中回收内存，因此它是非阻塞的。
// 这也是该命令名字的由来：仅将keys从keyspace元数据中删除，真正的删除会在后续异步操作。
func Unlink(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "unlink", keys...)
	return
}

// Dump .
// Redis DUMP 命令用于序列化给定 key ，并返回被序列化的值。
func Dump(ctx context.Context, key interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "dump", key)
	return
}

// Exists .
// Redis EXISTS 命令用于检查给定 key 是否存在。
func Exists(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "exists", keys...)
	return
}

// Expire .
// Redis Expire 命令用于设置 key 的过期时间，key 过期后将不再可用。单位以秒计。
func Expire(ctx context.Context, key string, expiration time.Duration) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "expire", key, formatSec(expiration))
	return
}

// ExpireAt .
// Redis Expireat 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间。key 过期后将不再可用。
func ExpireAt(ctx context.Context, key string, tm time.Time) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "expireat", key, tm.Unix())
	return
}

// Keys .
// Redis Keys 命令用于查找所有符合给定模式 pattern 的 key 。。
func Keys(ctx context.Context, pattern string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "keys", pattern)
	return
}

// Migrate .
// 将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
// 这个命令是一个原子操作，它在执行的时候会阻塞进行迁移的两个实例，直到以下任意结果发生：迁移成功，迁移失败，等到超时。
func Migrate(ctx context.Context, host, port, key string, db int64, timeout time.Duration) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "migrate",
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
	reply, err = _client.Do(ctx, "move", key, db)
	return
}

// ObjectRefCount .
// OBJECT REFCOUNT 该命令主要用于调试(debugging)，它能够返回指定key所对应value被引用的次数.
// OBJECT REFCOUNT <key> 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
func ObjectRefCount(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "object", "refcount", key)
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
	reply, err = _client.Do(ctx, "object", "encoding", key)
	return
}

// ObjectIdleTime .
// OBJECT IDLETIME <key> 返回给定 key 自储存以来的空转时间(idle， 没有被读取也没有被写入)，以秒为单位。
func ObjectIdleTime(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "object", "idletime", key)
	return
}

// Persist .
// Redis PERSIST 命令用于移除给定 key 的过期时间，使得 key 永不过期。
func Persist(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "persist", key)
	return
}

// PExpire .
// Redis PEXPIRE 命令和 EXPIRE 命令的作用类似，但是它以毫秒为单位设置 key 的生存时间，而不像 EXPIRE 命令那样，以秒为单位。
func PExpire(ctx context.Context, key string, expiration time.Duration) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "pexpire", key, formatMs(expiration))
	return
}

// PExpireAt .
// Redis PEXPIREAT 命令用于设置 key 的过期时间，以毫秒计。key 过期后将不再可用。
func PExpireAt(ctx context.Context, key string, tm time.Time) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "pexpireat",
		key,
		tm.UnixNano()/int64(time.Millisecond),
	)
	return
}

// PTTL .
// Redis Pttl 命令以毫秒为单位返回 key 的剩余过期时间。
func PTTL(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "pttl", key)
	return
}

// RandomKey .
// Redis RANDOMKEY 命令从当前数据库中随机返回一个 key 。
func RandomKey(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "randomkey")
	return
}

// Rename .
// Redis Rename 命令用于修改 key 的名称 。
func Rename(ctx context.Context, key, newkey string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "rename", key, newkey)
	return
}

// RenameNX .
// Redis Renamenx 命令用于在新的 key 不存在时修改 key 的名称 。
func RenameNX(ctx context.Context, key, newkey string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "renamenx", key, newkey)
	return
}

// Restore .
// 反序列化给定的序列化值，并将它和给定的 key 关联。
// 那么使用反序列化得出的值来代替键 key 原有的值； 相反地， 如果键 key 已经存在， 但是没有给定 REPLACE 选项， 那么命令返回一个错误。
func Restore(ctx context.Context, key string, ttl time.Duration, value string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "restore", key, formatMs(ttl), value)
	return
}

// RestoreReplace .
// 反序列化给定的序列化值，并将它和给定的 key 关联。
// 那么使用反序列化得出的值来代替键 key 原有的值； 相反地， 如果键 key 已经存在， 但是没有给定 REPLACE 选项， 那么命令返回一个错误。
func RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "restore",
		key,
		formatMs(ttl),
		value,
		"replace",
	)
	return
}

// SortArgs .
type SortArgs struct {
	By            string
	Offset, Count int64
	Get           []string
	Order         string
	Alpha         bool
}

// args .
func (sort *SortArgs) args(key string) []interface{} {
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
func Sort(ctx context.Context, key string, sort *SortArgs) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sort", sort.args(key)...)
	return
}

// SortStore .
// redis支持对list，set，sorted set、hash元素（元素可以为数值与字符串）的排序。
// sort key [BY pattern] [LIMIT start count] [GET pattern] [ASC|DESC] [ALPHA] [STORE dstkey]
// 如果对集合经常按照固定的模式去排序，那么把排序结果缓存起来会减少不少cpu开销，使用store选项可以将排序内容保存到指定key中，保存的类型是list
func SortStore(ctx context.Context, key, store string, sort *SortArgs) (reply interface{}, err error) {
	args := sort.args(key)
	if store != "" {
		args = append(args, "store", store)
	}
	reply, err = _client.Do(ctx, "sort", args...)
	return
}

// SortInterfaces .
// redis支持对list，set，sorted set、hash元素（元素可以为数值与字符串）的排序。
// sort key [BY pattern] [LIMIT start count] [GET pattern] [ASC|DESC] [ALPHA] [STORE dstkey]
func SortInterfaces(ctx context.Context, key string, sort *SortArgs) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sort", sort.args(key)...)
	return
}

// Touch .
// 修改指定 key 的 最后访问时间。忽略不存在的 key。
func Touch(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "touch", keys...)
	return
}

// TTL .
// Redis TTL 命令以秒为单位返回 key 的剩余过期时间。
func TTL(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "ttl", key)
	return
}

// Type .
// Redis Type 命令用于返回 key 所储存的值的类型。
func Type(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "type", key)
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
	reply, err = _client.Do(ctx, "scan", args...)
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
	reply, err = _client.Do(ctx, "sscan", args...)
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
	reply, err = _client.Do(ctx, "hscan", args...)
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
	reply, err = _client.Do(ctx, "zscan", args...)
	return
}

//------------------------------------------------------------------------------

// Append .
// Redis Append 命令用于为指定的 key 追加值。
// 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func Append(ctx context.Context, key, value string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "append", key, value)
	return
}

// BitCountArgs .
type BitCountArgs struct {
	Start, End int64
}

// BitCount .
// 统计指定位区间上值为1的个数
// BITCOUNT key [start end]
// 从左向右从0开始，从右向左从-1开始，注意start和end是字节
func BitCount(ctx context.Context, key string, bitCount *BitCountArgs) (reply interface{}, err error) {
	args := []interface{}{key}
	if bitCount != nil {
		args = append(
			args,
			bitCount.Start,
			bitCount.End,
		)
	}
	reply, err = _client.Do(ctx, "bitcount", args...)
	return
}

// bitOp
func bitOp(op, destKey string, keys ...string) []interface{} {
	args := make([]interface{}, 2+len(keys))
	args[0] = op
	args[1] = destKey
	for i, key := range keys {
		args[2+i] = key
	}
	return args
}

// BitOpAnd .
// BITOP AND destkey key [key ...]  ，对一个或多个key求逻辑并，并将结果保存到destkey
func BitOpAnd(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("and", destKey, keys...)
	reply, err = _client.Do(ctx, "bitop", args...)
	return
}

// BitOpOr .
// BITOP OR destkey key [key ...] ，对一个或多个key求逻辑或，并将结果保存到destkey
func BitOpOr(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("or", destKey, keys...)
	reply, err = _client.Do(ctx, "bitop", args...)
	return
}

// BitOpXor .
// BITOP XOR destkey key [key ...] ，对一个或多个key求逻辑异或，并将结果保存到destkey
func BitOpXor(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("xor", destKey, keys...)
	reply, err = _client.Do(ctx, "bitop", args...)
	return
}

// BitOpNot .
// BITOP NOT destkey key ，对给定key求逻辑非，并将结果保存到destkey
func BitOpNot(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("not", destKey, keys...)
	reply, err = _client.Do(ctx, "bitop", args...)
	return
}

// BitPos .
// 返回字符串里面第一个被设置为1或者0的bit位。
func BitPos(ctx context.Context, key string, bit int64, pos ...int64) (reply interface{}, err error) {
	args := make([]interface{}, 2+len(pos))
	args[0] = key
	args[1] = bit
	switch len(pos) {
	case 0:
	case 1:
		args[2] = pos[0]
	case 2:
		args[2] = pos[0]
		args[3] = pos[1]
	default:
		panic("too many arguments")
	}
	reply, err = _client.Do(ctx, "bitpos", args...)
	return
}

// Decr .
// Redis Decr 命令将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func Decr(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "decr", key)
	return
}

// DecrBy .
// Redis Decrby 命令将 key 所储存的值减去指定的减量值。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func DecrBy(ctx context.Context, key string, decrement int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "decrby", key, decrement)
	return
}

// Get .
// Redis Get 命令用于获取指定 key 的值。如果 key 不存在，返回 nil 。如果key 储存的值不是字符串类型，返回一个错误。
func Get(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "get", key)
	return
}

// GetBit .
// Redis Getbit 命令用于对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
func GetBit(ctx context.Context, key string, offset int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "getbit", key, offset)
	return
}

// GetRange .
// Redis Getrange 命令用于获取存储在指定 key 中字符串的子字符串。字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
func GetRange(ctx context.Context, key string, start, end int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "getrange", key, start, end)
	return
}

// GetSet .
// Redis Getset 命令用于设置指定 key 的值，并返回 key 的旧值。
func GetSet(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "getset", key, value)
	return
}

// Incr .
// Redis Incr 命令将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func Incr(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "incr", key)
	return
}

// IncrBy .
// Redis Incrby 命令将 key 中储存的数字加上指定的增量值。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func IncrBy(ctx context.Context, key string, value int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "incrby", key, value)
	return
}

// IncrByFloat .
// Redis Incrbyfloat 命令为 key 中所储存的值加上指定的浮点数增量值。
// 如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
func IncrByFloat(ctx context.Context, key string, value float64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "incrbyfloat", key, value)
	return
}

// appendArgs .
func appendArgs(dst, src []interface{}) []interface{} {
	if len(src) == 1 {
		if ss, ok := src[0].([]string); ok {
			for _, s := range ss {
				dst = append(dst, s)
			}
			return dst
		}
	}

	for _, v := range src {
		dst = append(dst, v)
	}
	return dst
}

// MGet .
// Redis Mget 命令返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。
func MGet(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "mget", keys...)
	return
}

// MSet .
// Redis Mset 命令用于同时设置一个或多个 key-value 对。
func MSet(ctx context.Context, pairs ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, len(pairs))
	args = appendArgs(args, pairs)
	reply, err = _client.Do(ctx, "mset", args...)
	return
}

// MSetNX .
// Redis Msetnx 命令用于所有给定 key 都不存在时，同时设置一个或多个 key-value 对。
// 当所有 key 都成功设置，返回 1 。 如果所有给定 key 都设置失败(至少有一个 key 已经存在)，那么返回 0 。
func MSetNX(ctx context.Context, pairs ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, len(pairs))
	args = appendArgs(args, pairs)
	reply, err = _client.Do(ctx, "msetnx", args...)
	return
}

// Set .
// Redis SET 命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型。
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (reply interface{}, err error) {
	args := make([]interface{}, 2, 3)
	args[0] = key
	args[1] = value
	if expiration > 0 {
		if usePrecise(expiration) {
			args = append(args, "px", formatMs(expiration))
		} else {
			args = append(args, "ex", formatSec(expiration))
		}
	}
	reply, err = _client.Do(ctx, "set", args...)
	return
}

// SetBit .
// Redis Setbit 命令用于对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
func SetBit(ctx context.Context, key string, offset int64, value int) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "setbit", key, offset, value)
	return
}

// SetNX .
// Redis Setnx（SET if Not eXists） 命令在指定的 key 不存在时，为 key 设置指定的值。
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (reply interface{}, err error) {
	if expiration == 0 {
		// Use old `SETNX` to support old Redis versions.
		reply, err = _client.Do(ctx, "setnx", key, value)
	} else {
		if usePrecise(expiration) {
			reply, err = _client.Do(ctx, "set", key, value, "px", formatMs(expiration), "nx")
		} else {
			reply, err = _client.Do(ctx, "set", key, value, "ex", formatMs(expiration), "nx")
		}
	}
	return
}

// SetXX .
// Redis `SET key value [expiration] XX` command.
// Zero expiration means the key has no expiration time.
// NX ： 只在键不存在时， 才对键进行设置操作。 执行 SET key value NX 的效果等同于执行 SETNX key value 。
// XX ： 只在键已经存在时， 才对键进行设置操作。
func SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (reply interface{}, err error) {
	if expiration == 0 {
		reply, err = _client.Do(ctx, "set", key, value, "xx")
	} else {
		if usePrecise(expiration) {
			reply, err = _client.Do(ctx, "set", key, value, "px", formatMs(expiration), "xx")
		} else {
			reply, err = _client.Do(ctx, "set", key, value, "ex", formatSec(expiration), "xx")
		}
	}
	return
}

// SetRange .
// Redis Setrange 命令用指定的字符串覆盖给定 key 所储存的字符串值，覆盖的位置从偏移量 offset 开始。
func SetRange(ctx context.Context, key string, offset int64, value string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "setrange", key, offset, value)
	return
}

// StrLen .
// Redis Strlen 命令用于获取指定 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误。
func StrLen(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "strlen", key)
	return
}

//------------------------------------------------------------------------------

// HDel .
// Redis Hdel 命令用于删除哈希表 key 中的一个或多个指定字段，不存在的字段将被忽略。
func HDel(ctx context.Context, key string, fields ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(fields))
	args[0] = key
	for i, field := range fields {
		args[1+i] = field
	}
	reply, err = _client.Do(ctx, "hdel", args...)
	return
}

// HExists .
// Redis Hexists 命令用于查看哈希表的指定字段是否存在。
func HExists(ctx context.Context, key string, field string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hexists", key, field)
	return
}

// HGet .
// Redis Hget 命令用于返回哈希表中指定字段的值。
func HGet(ctx context.Context, key string, field string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hget", key, field)
	return
}

// HGetAll .
// Redis Hgetall 命令用于返回哈希表中，所有的字段和值。
// 在返回值里，紧跟每个字段名(field name)之后是字段的值(value)，所以返回值的长度是哈希表大小的两倍
func HGetAll(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hgetall", key)
	return
}

// HIncrBy .
// Redis Hincrby 命令用于为哈希表中的字段值加上指定增量值。
// 增量也可以为负数，相当于对指定字段进行减法操作。
// 如果哈希表的 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0 。
// 对一个储存字符串值的字段执行 HINCRBY 命令将造成一个错误。
// 本操作的值被限制在 64 位(bit)有符号数字表示之内。
func HIncrBy(ctx context.Context, key, field string, incr int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hincrby", key, field, incr)
	return
}

// HIncrByFloat .
// Redis Hincrbyfloat 命令用于为哈希表中的字段值加上指定浮点数增量值。
// 如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0
func HIncrByFloat(ctx context.Context, key, field string, incr float64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hincrbyfloat", key, field, incr)
	return
}

// HKeys .
// Redis Hkeys 命令用于获取哈希表中的所有域（field）。
func HKeys(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hkeys", key)
	return
}

// HLen .
// Redis Hlen 命令用于获取哈希表中字段的数量。
func HLen(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hlen", key)
	return
}

// HMGet .
// Redis Hmget 命令用于返回哈希表中，一个或多个给定字段的值。
// 如果指定的字段不存在于哈希表，那么返回一个 nil 值。
func HMGet(ctx context.Context, key string, fields ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(fields))
	args[0] = key
	for i, field := range fields {
		args[1+i] = field
	}
	reply, err = _client.Do(ctx, "hmget", args...)
	return
}

// HMSet .
// Redis Hmset 命令用于同时将多个 field-value (字段-值)对设置到哈希表中。
// 此命令会覆盖哈希表中已存在的字段。
// 如果哈希表不存在，会创建一个空哈希表，并执行 HMSET 操作。
func HMSet(ctx context.Context, key string, fields map[string]interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(fields)*2)
	args[0] = key
	i := 1
	for k, v := range fields {
		args[i] = k
		args[i+1] = v
		i += 2
	}
	reply, err = _client.Do(ctx, "hmset", args...)
	return
}

// HSet .
// Redis Hset 命令用于为哈希表中的字段赋值 。
// 如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果字段已经存在于哈希表中，旧值将被覆盖。
func HSet(ctx context.Context, key, field string, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hset", key, field, value)
	return
}

// HSetNX .
// Redis Hsetnx 命令用于为哈希表中不存在的的字段赋值 。
// 如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果字段已经存在于哈希表中，操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func HSetNX(ctx context.Context, key, field string, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hsetnx", key, field, value)
	return
}

// HVals .
// Redis Hvals 命令返回哈希表所有域(field)的值。
func HVals(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "hvals", key)
	return
}

//------------------------------------------------------------------------------

// BLPop .
// Redis Blpop 命令移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func BLPop(ctx context.Context, timeout time.Duration, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[i] = key
	}
	args[len(args)-1] = formatSec(timeout)
	//ctx, cancel := context.WithTimeout(ctx, timeout)
	//defer cancel()
	reply, err = _client.Do(ctx, "blpop", args...)
	return
}

// BRPop .
// Redis Brpop 命令移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func BRPop(ctx context.Context, timeout time.Duration, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[i] = key
	}
	args[len(args)-1] = formatSec(timeout)
	//ctx, cancel := context.WithTimeout(ctx, timeout)
	//defer cancel()
	reply, err = _client.Do(ctx, "brpop", args...)
	return
}

// BRPopLPush .
// Redis Brpoplpush 命令从列表中取出最后一个元素，并插入到另外一个列表的头部； 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
// 假如在指定时间内没有任何元素被弹出，则返回一个 nil 和等待时长。 反之，返回一个含有两个元素的列表，第一个元素是被弹出元素的值，第二个元素是等待时长。
func BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "brpoplpush",
		source,
		destination,
		formatSec(timeout),
	)
	return
}

// LIndex .
// Redis Lindex 命令用于通过索引获取列表中的元素。你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func LIndex(ctx context.Context, key string, index int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lindex", key, index)
	return
}

// LInsert .
// Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
// 当列表不存在时，被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// LINSERT key BEFORE|AFTER pivot value
func LInsert(ctx context.Context, key, op string, pivot, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "linsert", key, op, pivot, value)
	return
}

// LInsertBefore .
// Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
// 当列表不存在时，被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// LINSERT key BEFORE|AFTER pivot value
func LInsertBefore(ctx context.Context, key, op string, pivot, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "linsert", key, "before", pivot, value)
	return
}

// LInsertAfter .
// Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
// 当列表不存在时，被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// LINSERT key BEFORE|AFTER pivot value
func LInsertAfter(ctx context.Context, key, op string, pivot, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "linsert", key, "after", pivot, value)
	return
}

// LLen .
// Redis Llen 命令用于返回列表的长度。 如果列表 key 不存在，则 key 被解释为一个空列表，返回 0 。 如果 key 不是列表类型，返回一个错误。
func LLen(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "llen", key)
	return
}

// LPop .
// Redis Lpop 命令用于移除并返回列表的第一个元素。
// 列表的第一个元素。 当列表 key 不存在时，返回 nil 。
func LPop(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lpop", key)
	return
}

// LPush .
// Redis Lpush 命令将一个或多个值插入到列表头部。 如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。 当 key 存在但不是列表类型时，返回一个错误。
// 注意：在Redis 2.4版本以前的 LPUSH 命令，都只接受单个 value 值。
func LPush(ctx context.Context, key string, values ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = appendArgs(args, values)
	reply, err = _client.Do(ctx, "lpush", args...)
	return
}

// LPushX .
// Redis Lpushx 将一个值插入到已存在的列表头部，列表不存在时操作无效。
func LPushX(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lpushx", key, value)
	return
}

// LRange .
// Redis Lrange 返回列表中指定区间内的元素，区间以偏移量 START 和 END 指定。
// 其中 0 表示列表的第一个元素， 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func LRange(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lrange", key, start, stop)
	return
}

// LRem .
// Redis Lrem 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素。
// COUNT 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
// count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
// count = 0 : 移除表中所有与 VALUE 相等的值。
func LRem(ctx context.Context, key string, count int64, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lrem", key, count, value)
	return
}

// LSet .
// Redis Lset 通过索引来设置元素的值。
// 当索引参数超出范围，或对一个空列表进行 LSET 时，返回一个错误。
// 关于列表下标的更多信息，请参考 LINDEX 命令。
func LSet(ctx context.Context, key string, index int64, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lset", key, index, value)
	return
}

// LTrim .
// Redis Ltrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 下标 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func LTrim(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "ltrim", key, start, stop)
	return
}

// RPop .
// Redis Rpop 命令用于移除列表的最后一个元素，返回值为移除的元素。
func RPop(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "rpop", key)
	return
}

// RPopLPush .
// Redis Rpoplpush 命令用于移除列表的最后一个元素，并将该元素添加到另一个列表并返回。
func RPopLPush(ctx context.Context, source, destination string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "rpoplpush", source, destination)
	return
}

// RPush .
// Redis Rpush 命令用于将一个或多个值插入到列表的尾部(最右边)。
// 如果列表不存在，一个空列表会被创建并执行 RPUSH 操作。 当列表存在但不是列表类型时，返回一个错误。
// 注意：在 Redis 2.4 版本以前的 RPUSH 命令，都只接受单个 value 值。
func RPush(ctx context.Context, key string, values ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = appendArgs(args, values)
	reply, err = _client.Do(ctx, "rpush", args...)
	return
}

// RPushX .
// Redis rpushx，命令用于将一个或多个值插入到已存在的列表尾部(最右边)，如果列表不存在，操作无效。
func RPushX(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "rpushx", key, value)
	return
}

//------------------------------------------------------------------------------

// SAdd .
// Redis Sadd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
// 假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
// 当集合 key 不是集合类型时，返回一个错误。
// 注意：在 Redis2.4 版本以前， SADD 只接受单个成员值。
func SAdd(ctx context.Context, key string, members ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(members))
	args[0] = key
	args = appendArgs(args, members)
	reply, err = _client.Do(ctx, "sadd", args...)
	return
}

// SCard .
// Redis Scard 命令返回集合中元素的数量。
func SCard(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "scard", key)
	return
}

// SDiff .
// Redis Sdiff 命令返回第一个集合与其他集合之间的差异，也可以认为说第一个集合中独有的元素。不存在的集合 key 将视为空集。
// 差集的结果来自前面的 FIRST_KEY ,而不是后面的 OTHER_KEY1，也不是整个 FIRST_KEY OTHER_KEY1..OTHER_KEYN 的差集。
func SDiff(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sdiff", keys...)
	return
}

// SDiffStore .
// Redis Sdiffstore 命令将给定集合之间的差集存储在指定的集合中。如果指定的集合 key 已存在，则会被覆盖。
func SDiffStore(ctx context.Context, destination string, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(keys))
	args[0] = destination
	for i, key := range keys {
		args[1+i] = key
	}
	reply, err = _client.Do(ctx, "sdiffstore", args...)
	return
}

// SInter .
// Redis Sinter 命令返回给定所有给定集合的交集。 不存在的集合 key 被视为空集。 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
func SInter(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sinter", keys...)
	return
}

// SInterStore .
// Redis Sinterstore 命令将给定集合之间的交集存储在指定的集合中。如果指定的集合已经存在，则将其覆盖。
func SInterStore(ctx context.Context, destination string, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(keys))
	args[0] = destination
	for i, key := range keys {
		args[1+i] = key
	}
	reply, err = _client.Do(ctx, "sinterstore", args...)
	return
}

// SIsMember .
// Redis Sismember 命令判断成员元素是否是集合的成员。
func SIsMember(ctx context.Context, key string, member interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sismember", key, member)
	return
}

// SMembers .
// Redis Smembers 命令返回集合中的所有的成员。 不存在的集合 key 被视为空集合。
func SMembers(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "smembers", key)
	return
}

// SMove .
// Redis Smove 命令将指定成员 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。
// 否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
func SMove(ctx context.Context, source, destination string, member interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "smove", source, destination, member)
	return
}

// SPop .
// Redis Spop 命令用于移除集合中的指定 key 的一个或多个随机元素，移除后会返回移除的元素。
// 该命令类似 Srandmember 命令，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SPop(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "spop", key)
	return
}

// SPopN .
// Redis Spop 命令用于移除集合中的指定 key 的一个或多个随机元素，移除后会返回移除的元素。
// 该命令类似 Srandmember 命令，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SPopN(ctx context.Context, key string, count int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "spop", key, count)
	return
}

// SRandMember .
// Redis Srandmember 命令用于返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， Srandmember 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SRandMember(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "srandmember", key)
	return
}

// SRandMemberN .
// Redis Srandmember 命令用于返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， Srandmember 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SRandMemberN(ctx context.Context, key string, count int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "srandmember", key, count)
	return
}

// SRem .
// Redis Srem 命令用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略。
// 当 key 不是集合类型，返回一个错误。
// 在 Redis 2.4 版本以前， SREM 只接受单个成员值。
func SRem(ctx context.Context, key string, members ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(members))
	args[0] = key
	args = appendArgs(args, members)
	reply, err = _client.Do(ctx, "srem", args...)
	return
}

// SUnion .
// Redis Sunion 命令返回给定集合的并集。不存在的集合 key 被视为空集。
func SUnion(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sunion", keys...)
	return
}

// SUnionStore .
// Redis Sunionstore 命令将给定集合的并集存储在指定的集合 destination 中。如果 destination 已经存在，则将其覆盖。
func SUnionStore(ctx context.Context, destination string, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(keys))
	args[0] = destination
	for i, key := range keys {
		args[1+i] = key
	}
	reply, err = _client.Do(ctx, "sunionstore", args...)
	return
}

//------------------------------------------------------------------------------

// XAddArgs .
type XAddArgs struct {
	Stream       string
	MaxLen       int64 // MAXLEN N
	MaxLenApprox int64 // MAXLEN ~ N
	ID           string
	Values       map[string]interface{}
}

// XAdd .
// https://www.runoob.com/redis/redis-stream.html
// XADD - 添加消息到末尾
func XAdd(ctx context.Context, a *XAddArgs) (reply interface{}, err error) {
	args := make([]interface{}, 0, 5+len(a.Values)*2)
	args = append(args, a.Stream)
	if a.MaxLen > 0 {
		args = append(args, "maxlen", a.MaxLen)
	} else if a.MaxLenApprox > 0 {
		args = append(args, "maxlen", "~", a.MaxLenApprox)
	}
	if a.ID != "" {
		args = append(args, a.ID)
	} else {
		args = append(args, "*")
	}
	for k, v := range a.Values {
		args = append(args, k)
		args = append(args, v)
	}

	reply, err = _client.Do(ctx, "xadd", args...)
	return
}

// XDel .
// XDel - 删除消息
func XDel(ctx context.Context, stream string, ids ...string) (reply interface{}, err error) {
	args := []interface{}{stream}
	for _, id := range ids {
		args = append(args, id)
	}
	reply, err = _client.Do(ctx, "xdel", args...)
	return
}

// XLen .
// XLEN - 获取流包含的元素数量，即消息长度
func XLen(ctx context.Context, stream string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xlen", stream)
	return
}

// XRange .
// XRANGE - 获取消息列表，会自动过滤已经删除的消息
func XRange(ctx context.Context, stream, start, stop string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xrange", stream, start, stop)
	return
}

// XRangeN .
// XRANGE - 获取消息列表，会自动过滤已经删除的消息
func XRangeN(ctx context.Context, stream, start, stop string, count int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xrange", stream, start, stop, "count", count)
	return
}

// XRevRange .
// XREVRANGE - 反向获取消息列表，ID 从大到小
func XRevRange(ctx context.Context, stream, start, stop string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xrevrange", stream, start, stop)
	return
}

// XRevRangeN .
// XREVRANGE - 反向获取消息列表，ID 从大到小
func XRevRangeN(ctx context.Context, stream, start, stop string, count int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xrevrange", stream, start, stop, "count", count)
	return
}

// XReadArgs .
type XReadArgs struct {
	Streams []string
	Count   int64
	Block   time.Duration
}

// XRead .
// XREAD - 以阻塞或非阻塞方式获取消息列表
func XRead(ctx context.Context, a *XReadArgs) (reply interface{}, err error) {
	args := make([]interface{}, 0, 4+len(a.Streams))
	if a.Count > 0 {
		args = append(args, "count")
		args = append(args, a.Count)
	}
	if a.Block >= 0 {
		args = append(args, "block")
		args = append(args, int64(a.Block/time.Millisecond))
	}
	args = append(args, "streams")
	for _, s := range a.Streams {
		args = append(args, s)
	}
	if a.Block >= 0 {
		//ctx, cancel := context.WithTimeout(ctx, a.Block)
		//defer cancel()
	}
	reply, err = _client.Do(ctx, "xread", args...)
	return
}

// XReadStreams .
// XREAD - 以阻塞或非阻塞方式获取消息列表
func XReadStreams(ctx context.Context, streams ...string) (reply interface{}, err error) {
	reply, err = XRead(ctx, &XReadArgs{
		Streams: streams,
		Block:   -1,
	})
	return
}

// XGroupCreate .
// XGROUP CREATE - 创建消费者组
func XGroupCreate(ctx context.Context, stream, group, start string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xgroup", "create", stream, group, start)
	return
}

// XGroupCreateMkStream .
// XGROUP CREATE - 创建消费者组
func XGroupCreateMkStream(ctx context.Context, stream, group, start string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xgroup", "create", stream, group, start, "mkstream")
	return
}

// XGroupSetID .
// XGROUP SETID - 为消费者组设置新的最后递送消息ID
func XGroupSetID(ctx context.Context, stream, group, start string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xgroup", "setid", stream, group, start)
	return
}

// XGroupDestroy .
// XGROUP DESTROY - 删除消费者组
func XGroupDestroy(ctx context.Context, stream, group string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xgroup", "destroy", stream, group)
	return
}

// XGroupDelConsumer .
// XGROUP DELCONSUMER - 删除消费者
func XGroupDelConsumer(ctx context.Context, stream, group, consumer string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xgroup", "delconsumer", stream, group, consumer)
	return
}

// XReadGroupArgs .
type XReadGroupArgs struct {
	Group    string
	Consumer string
	// List of streams and ids.
	Streams []string
	Count   int64
	Block   time.Duration
	NoAck   bool
}

// XReadGroup .
// XREADGROUP GROUP - 读取消费者组中的消息
func XReadGroup(ctx context.Context, a *XReadGroupArgs) (reply interface{}, err error) {
	args := make([]interface{}, 0, 7+len(a.Streams))
	args = append(args, "group", a.Group, a.Consumer)
	if a.Count > 0 {
		args = append(args, "count", a.Count)
	}
	if a.Block >= 0 {
		args = append(args, "block", int64(a.Block/time.Millisecond))
	}
	if a.NoAck {
		args = append(args, "noack")
	}
	args = append(args, "streams")
	for _, s := range a.Streams {
		args = append(args, s)
	}
	if a.Block >= 0 {
		//ctx, cancel := context.WithTimeout(ctx, a.Block)
		//defer cancel()
	}
	reply, err = _client.Do(ctx, "xreadgroup", args...)
	return
}

// XAck .
// XACK - 将消息标记为"已处理"
func XAck(ctx context.Context, stream, group string, ids ...string) (reply interface{}, err error) {
	args := []interface{}{stream, group}
	for _, id := range ids {
		args = append(args, id)
	}
	reply, err = _client.Do(ctx, "xack", args...)
	return
}

// XPending .
// XPENDING - 显示待处理消息的相关信息
func XPending(ctx context.Context, stream, group string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xpending", stream, group)
	return
}

// XPendingExtArgs .
type XPendingExtArgs struct {
	Stream   string
	Group    string
	Start    string
	End      string
	Count    int64
	Consumer string
}

// XPendingExt .
// XPENDING - 显示待处理消息的相关信息
func XPendingExt(ctx context.Context, a *XPendingExtArgs) (reply interface{}, err error) {
	args := make([]interface{}, 0, 6)
	args = append(args, a.Stream, a.Group, a.Start, a.End, a.Count)
	if a.Consumer != "" {
		args = append(args, a.Consumer)
	}
	reply, err = _client.Do(ctx, "xpending", args...)
	return
}

// XClaimArgs .
type XClaimArgs struct {
	Stream   string
	Group    string
	Consumer string
	MinIdle  time.Duration
	Messages []string
}

// xClaimArgs .
func xClaimArgs(a *XClaimArgs) []interface{} {
	args := make([]interface{}, 0, 3+len(a.Messages))
	args = append(args, a.Stream, a.Group, a.Consumer, int64(a.MinIdle/time.Millisecond))
	for _, id := range a.Messages {
		args = append(args, id)
	}
	return args
}

// XClaim .
// XCLAIM - 转移消息的归属权
func XClaim(ctx context.Context, a *XClaimArgs) (reply interface{}, err error) {
	args := xClaimArgs(a)
	reply, err = _client.Do(ctx, "xclaim", args...)
	return
}

// XClaimJustID .
// XCLAIM - 转移消息的归属权
func XClaimJustID(ctx context.Context, a *XClaimArgs) (reply interface{}, err error) {
	args := xClaimArgs(a)
	args = append(args, "justid")
	reply, err = _client.Do(ctx, "xclaim", args...)
	return
}

// XTrim .
// XTRIM - 对流进行修剪，限制长度
func XTrim(ctx context.Context, key string, maxLen int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xtrim", key, "maxlen", maxLen)
	return
}

// XTrimApprox .
// XTRIM - 对流进行修剪，限制长度
func XTrimApprox(ctx context.Context, key string, maxLen int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "xtrim", key, "maxlen", "~", maxLen)
	return
}

//------------------------------------------------------------------------------

// BZPopMax .
// BZPOPMAX 是有序集合命令 ZPOPMAX带有阻塞功能的版本。
// 参数 timeout 可以理解为客户端被阻塞的最大秒数值，0 表示永久阻塞。
// 当有序集合空且执行超时时返回 nil
// 返回三元素结果，第一元素 key 名称，第二元素成员名称，第三元素分数。
func BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[1+i] = key
	}
	args[len(args)-1] = formatSec(timeout)
	reply, err = _client.Do(ctx, "bzpopmax", args...)
	return
}

// BZPopMin .
// BZPOPMIN 是有序集合命令 ZPOPMIN带有阻塞功能的版本。
// 参数 timeout 可以理解为客户端被阻塞的最大秒数值，0 表示永久阻塞。
// 当有序集合空且执行超时时返回 nil
// 返回三元素结果，第一元素 key 名称，第二元素成员名称，第三元素分数。
func BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[1+i] = key
	}
	args[len(args)-1] = formatSec(timeout)
	reply, err = _client.Do(ctx, "bzpopmin", args...)
	return
}

// Z represents sorted set member.
type Z struct {
	Score  float64
	Member interface{}
}

// ZWithKey represents sorted set member including the name of the key where it was popped.
type ZWithKey struct {
	Z
	Key string
}

// ZStore is used as an arg to ZInterStore and ZUnionStore.
type ZStore struct {
	Weights []float64
	// Can be SUM, MIN or MAX.
	Aggregate string
}

// zAdd .
func zAdd(a []interface{}, n int, members ...Z) []interface{} {
	for i, m := range members {
		a[n+2*i] = m.Score
		a[n+2*i+1] = m.Member
	}
	return a
}

// ZAdd .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
func ZAdd(ctx context.Context, key string, members ...Z) (reply interface{}, err error) {
	const n = 1
	a := make([]interface{}, n+2*len(members))
	a[0] = key
	args := zAdd(a, n, members...)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZAddNX .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
func ZAddNX(ctx context.Context, key string, members ...Z) (reply interface{}, err error) {
	const n = 2
	a := make([]interface{}, n+2*len(members))
	a[0], a[1] = key, "nx"
	args := zAdd(a, n, members...)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZAddXX .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
func ZAddXX(ctx context.Context, key string, members ...Z) (reply interface{}, err error) {
	const n = 2
	a := make([]interface{}, n+2*len(members))
	a[0], a[1] = key, "xx"
	args := zAdd(a, n, members...)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZAddCh .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
// CH: 修改返回值为发生变化的成员总数，原始是返回新添加成员的总数 (CH 是 changed 的意思)。
// 更改的元素是新添加的成员，已经存在的成员更新分数。 所以在命令中指定的成员有相同的分数将不被计算在内。
// 注：在通常情况下，ZADD返回值只计算新添加成员的数量。
func ZAddCh(ctx context.Context, key string, members ...Z) (reply interface{}, err error) {
	const n = 2
	a := make([]interface{}, n+2*len(members))
	a[0], a[1] = key, "ch"
	args := zAdd(a, n, members...)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZAddNXCh .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
// CH: 修改返回值为发生变化的成员总数，原始是返回新添加成员的总数 (CH 是 changed 的意思)。
// 更改的元素是新添加的成员，已经存在的成员更新分数。 所以在命令中指定的成员有相同的分数将不被计算在内。
// 注：在通常情况下，ZADD返回值只计算新添加成员的数量。
func ZAddNXCh(ctx context.Context, key string, members ...Z) (reply interface{}, err error) {
	const n = 3
	a := make([]interface{}, n+2*len(members))
	a[0], a[1], a[2] = key, "nx", "ch"
	args := zAdd(a, n, members...)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZAddXXCh .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
// CH: 修改返回值为发生变化的成员总数，原始是返回新添加成员的总数 (CH 是 changed 的意思)。
// 更改的元素是新添加的成员，已经存在的成员更新分数。 所以在命令中指定的成员有相同的分数将不被计算在内。
// 注：在通常情况下，ZADD返回值只计算新添加成员的数量。
func ZAddXXCh(ctx context.Context, key string, members ...Z) (reply interface{}, err error) {
	const n = 3
	a := make([]interface{}, n+2*len(members))
	a[0], a[1], a[2] = key, "xx", "ch"
	args := zAdd(a, n, members...)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// zIncr .
func zIncr(a []interface{}, n int, members ...Z) []interface{} {
	for i, m := range members {
		a[n+2*i] = m.Score
		a[n+2*i+1] = m.Member
	}
	return a
}

// ZIncr .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
// INCR: 当ZADD指定这个选项时，成员的操作就等同ZINCRBY命令，对成员的分数进行递增操作。
func ZIncr(ctx context.Context, key string, member Z) (reply interface{}, err error) {
	const n = 2
	a := make([]interface{}, n+2)
	a[0], a[1] = key, "incr"
	args := zIncr(a, n, member)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZIncrNX .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
// INCR: 当ZADD指定这个选项时，成员的操作就等同ZINCRBY命令，对成员的分数进行递增操作。
func ZIncrNX(ctx context.Context, key string, member Z) (reply interface{}, err error) {
	const n = 3
	a := make([]interface{}, n+2)
	a[0], a[1], a[2] = key, "incr", "nx"
	args := zIncr(a, n, member)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZIncrXX .
// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上。
// 分数值可以是整数值或双精度浮点数。
// 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZADD 每次只能添加一个元素。
// XX: 仅仅更新存在的成员，不添加新成员。
// NX: 不更新存在的成员。只添加新成员。
// INCR: 当ZADD指定这个选项时，成员的操作就等同ZINCRBY命令，对成员的分数进行递增操作。
func ZIncrXX(ctx context.Context, key string, member Z) (reply interface{}, err error) {
	const n = 3
	a := make([]interface{}, n+2)
	a[0], a[1], a[2] = key, "incr", "xx"
	args := zIncr(a, n, member)
	reply, err = _client.Do(ctx, "zadd", args...)
	return
}

// ZCard .
// Redis ZCARD 命令用于返回有序集的成员个数。
// 整数: 返回有序集的成员个数，当 key 不存在时，返回 0 。
func ZCard(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zcard", key)
	return
}

// ZCount .
// Redis Zcount 命令用于计算有序集合中指定分数区间的成员数量。
// 分数值在 min 和 max 之间的成员的数量。
func ZCount(ctx context.Context, key, min, max string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zcount", key, min, max)
	return
}

// ZLexCount .
// Redis Zlexcount 命令在计算有序集合中指定字典区间内成员数量。
func ZLexCount(ctx context.Context, key, min, max string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zlexcount", key, min, max)
	return
}

// ZIncrBy .
// Redis Zincrby 命令对有序集合中指定成员的分数加上增量 increment
// 可以通过传递一个负数值 increment ，让分数减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或分数不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
// 当 key 不是有序集类型时，返回一个错误。
// 分数值可以是整数值或双精度浮点数。
func ZIncrBy(ctx context.Context, key string, increment float64, member string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zincrby", key, increment, member)
	return
}

// ZInterStore .
// Redis Zinterstore 命令计算给定的一个或多个有序集的交集，
// 其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的分数值是所有给定集下该成员分数值之和。
// [WEIGHTS weight [weight …]] ～ 指定参与交集运算各集合score的权重参数
// [AGGREGATE SUM|MIN|MAX] ～ 指定交集中元素的score的取值方式，例如：sum 等于各集合中该元素的score乘以权重求和
// 结果集中元素score为 取值方式计算各集合中该元素score乘以权重结果
// -----
// 例子 ：取zkey1 和zkey2 2个有序集合的交集 保存至有序集合zkey3,权重配置为 10 和 1（zkey3若不存在则新建，若存在则覆盖）
// zinterstore zkey3 2 zkey1 zkey2 weights 10 1
// zkey3 score = (zkey1 score)*10 + (zkey2 score)*1
// 例子 ：取zkey1 和zkey2 2个有序集合的交集 保存至有序集合zkey3,取值方式为min（zkey3若不存在则新建，若存在则覆盖）
// zkey3 score = min((zkey1 score)*1, (zkey2 score)*1)
func ZInterStore(ctx context.Context, destination string, store ZStore, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 2+len(keys))
	args[0] = destination
	args[1] = len(keys)
	for i, key := range keys {
		args[3+i] = key
	}
	if len(store.Weights) > 0 {
		args = append(args, "weights")
		for _, weight := range store.Weights {
			args = append(args, weight)
		}
	}
	if store.Aggregate != "" {
		args = append(args, "aggregate", store.Aggregate)
	}
	reply, err = _client.Do(ctx, "zinterstore", args...)
	return
}

// ZPopMax .
// ZPOPMAX 命令用于移除并弹出有序集合中分值最大的 count 个元素：
func ZPopMax(ctx context.Context, key string, count ...int64) (reply interface{}, err error) {
	args := []interface{}{key}

	switch len(count) {
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	}
	reply, err = _client.Do(ctx, "zpopmax", args...)
	return
}

// ZPopMin .
// ZPOPMIN 命令则用于移除并弹出有序集合中分值最小的 count 个元素：
func ZPopMin(ctx context.Context, key string, count ...int64) (reply interface{}, err error) {
	args := []interface{}{key}

	switch len(count) {
	case 0:
		break
	case 1:
		args = append(args, count[0])
	default:
		panic("too many arguments")
	}
	reply, err = _client.Do(ctx, "zpopmin", args...)
	return
}

// zRange .
func zRange(key string, start, stop int64, withScores bool) []interface{} {
	args := []interface{}{
		key,
		start,
		stop,
	}
	if withScores {
		args = append(args, "withscores")
	}
	return args
}

// ZRange .
// Redis Zrange 返回有序集中，指定区间内的成员。
// 其中成员的位置按分数值递增(从小到大)来排序。
// 具有相同分数值的成员按字典序(lexicographical order )来排列。
// 如果你需要成员按
// 值递减(从大到小)来排列，请使用 ZREVRANGE 命令。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// 指定区间内，带有分数值(可选)的有序集成员的列表。
func ZRange(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	args := zRange(key, start, start, false)
	reply, err = _client.Do(ctx, "zrange", args...)
	return
}

// ZRangeWithScores .
// Redis Zrange 返回有序集中，指定区间内的成员。
// 其中成员的位置按分数值递增(从小到大)来排序。
// 具有相同分数值的成员按字典序(lexicographical order )来排列。
// 如果你需要成员按
// 值递减(从大到小)来排列，请使用 ZREVRANGE 命令。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// 指定区间内，带有分数值(可选)的有序集成员的列表。
func ZRangeWithScores(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zrange", key, start, stop, "withscores")
	return
}

// ZRangeBy .
type ZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

// zRangeBy .
func zRangeBy(key string, opt ZRangeBy, withScores bool) []interface{} {
	args := []interface{}{key, opt.Min, opt.Max}
	if withScores {
		args = append(args, "withscores")
	}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	return args
}

// ZRangeByScore .
// Redis Zrangebyscore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列(该属性是有序集提供的，不需要额外的计算)。
// 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。
func ZRangeByScore(ctx context.Context, key string, opt ZRangeBy) (reply interface{}, err error) {
	args := zRangeBy(key, opt, false)
	reply, err = _client.Do(ctx, "zrangebyscore", args...)
	return
}

// ZRangeByLex .
// Redis Zrangebylex 通过字典区间返回有序集合的成员。
func ZRangeByLex(ctx context.Context, key string, opt ZRangeBy) (reply interface{}, err error) {
	args := zRangeBy(key, opt, false)
	reply, err = _client.Do(ctx, "zrangebylex", args...)
	return
}

// ZRangeByScoreWithScores .
// Redis Zrangebylex 通过字典区间返回有序集合的成员。
func ZRangeByScoreWithScores(ctx context.Context, key string, opt ZRangeBy) (reply interface{}, err error) {
	args := []interface{}{key, opt.Min, opt.Max, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	reply, err = _client.Do(ctx, "zrangebyscore", args...)
	return
}

// ZRank .
// Redis Zrank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。
func ZRank(ctx context.Context, key, member string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zrank", key, member)
	return
}

// ZRem .
// Redis Zrem 命令用于移除有序集中的一个或多个成员，不存在的成员将被忽略。
// 当 key 存在但不是有序集类型时，返回一个错误。
// 注意： 在 Redis 2.4 版本以前， ZREM 每次只能删除一个元素。
func ZRem(ctx context.Context, key string, members ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(members))
	args[0] = key
	args = appendArgs(args, members)
	reply, err = _client.Do(ctx, "zrem", args...)
	return
}

// ZRemRangeByRank .
// Redis Zremrangebyrank 命令用于移除有序集中，指定排名(rank)区间内的所有成员。
func ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zremrangebyrank", key, start, stop)
	return
}

// ZRemRangeByScore .
// Redis Zremrangebyscore 命令用于移除有序集中，指定分数（score）区间内的所有成员。
func ZRemRangeByScore(ctx context.Context, key, min, max string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zremrangebyscore", key, min, max)
	return
}

// ZRemRangeByLex .
// Redis Zremrangebylex 命令用于移除有序集合中给定的字典区间的所有成员。
func ZRemRangeByLex(ctx context.Context, key, min, max string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zremrangebylex", key, min, max)
	return
}

// ZRevRange .
// Redis Zrevrange 命令返回有序集中，指定区间内的成员。
// 其中成员的位置按分数值递减(从大到小)来排列。
// 具有相同分数值的成员按字典序的逆序(reverse lexicographical order)排列。
// 除了成员按分数值递减的次序排列这一点外， ZREVRANGE 命令的其他方面和 ZRANGE 命令一样。
func ZRevRange(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zrevrange", key, start, stop)
	return
}

// ZRevRangeWithScores .
// Redis Zrevrange 命令返回有序集中，指定区间内的成员。
// 其中成员的位置按分数值递减(从大到小)来排列。
// 具有相同分数值的成员按字典序的逆序(reverse lexicographical order)排列。
// 除了成员按分数值递减的次序排列这一点外， ZREVRANGE 命令的其他方面和 ZRANGE 命令一样。
func ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zrevrange", key, start, stop, "withscores")
	return
}

// zRevRangeBy .
func zRevRangeBy(key string, opt ZRangeBy) []interface{} {
	args := []interface{}{key, opt.Max, opt.Min}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	return args
}

// ZRevRangeByScore .
// Redis Zrevrangebyscore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序的逆序(reverse lexicographical order )排列。
// 除了成员按分数值递减的次序排列这一点外， ZREVRANGEBYSCORE 命令的其他方面和 ZRANGEBYSCORE 命令一样。
func ZRevRangeByScore(ctx context.Context, key string, opt ZRangeBy) (reply interface{}, err error) {
	args := zRevRangeBy(key, opt)
	reply, err = _client.Do(ctx, "zrevrangebyscore", args...)
	return
}

// ZRevRangeByLex .
// ZREVRANGEBYLEX 返回指定成员区间内的成员，按成员字典倒序排序, 分数必须相同。
// 在某些业务场景中,需要对一个字符串数组按名称的字典顺序进行倒序排列时,可以使用Redis中SortSet这种数据结构来处理。
func ZRevRangeByLex(ctx context.Context, key string, opt ZRangeBy) (reply interface{}, err error) {
	args := zRevRangeBy(key, opt)
	reply, err = _client.Do(ctx, "zrevrangebylex", args...)
	return
}

// ZRevRangeByScoreWithScores .
// Redis Zrevrangebyscore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序的逆序(reverse lexicographical order )排列。
// 除了成员按分数值递减的次序排列这一点外， ZREVRANGEBYSCORE 命令的其他方面和 ZRANGEBYSCORE 命令一样。
func ZRevRangeByScoreWithScores(ctx context.Context, key string, opt ZRangeBy) (reply interface{}, err error) {
	args := []interface{}{key, opt.Max, opt.Min, "withscores"}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"limit",
			opt.Offset,
			opt.Count,
		)
	}
	reply, err = _client.Do(ctx, "zrevrangebyscore", args...)
	return
}

// ZRevRank .
// Redis Zrevrank 命令返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。
// 排名以 0 为底，也就是说， 分数值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按分数值递增(从小到大)排列的排名。
func ZRevRank(ctx context.Context, key, member string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zrevrank", key, member)
	return
}

// ZScore .
// Redis Zscore 命令返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func ZScore(ctx context.Context, key, member string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "zscore", key, member)
	return
}

// ZUnionStore .
// Redis Zunionstore 命令计算给定的一个或多个有序集的并集，
// 其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的分数值是所有给定集下该成员分数值之和 。
// -----
// 例子 ：取zkey1 和zkey2 2个有序集合的并集 保存至有序集合zkey3,权重配置为 10 和 1（zkey3若不存在则新建，若存在则覆盖）
// zinterstore zkey3 2 zkey1 zkey2 weights 10 1
// zkey3 score = (zkey1 score)*10 + (zkey2 score)*1
// 例子 ：取zkey1 和zkey2 2个有序集合的并集 保存至有序集合zkey3,取值方式为min（zkey3若不存在则新建，若存在则覆盖）
// zkey3 score = min((zkey1 score)*1, (zkey2 score)*1)
func ZUnionStore(ctx context.Context, dest string, store ZStore, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 2+len(keys))
	args[0] = dest
	args[1] = len(keys)
	for i, key := range keys {
		args[2+i] = key
	}
	if len(store.Weights) > 0 {
		args = append(args, "weights")
		for _, weight := range store.Weights {
			args = append(args, weight)
		}
	}
	if store.Aggregate != "" {
		args = append(args, "aggregate", store.Aggregate)
	}
	reply, err = _client.Do(ctx, "zunionstore", args...)
	return
}

//------------------------------------------------------------------------------

// PFAdd .
// Redis Pfadd 命令将所有元素参数添加到 HyperLogLog 数据结构中。
func PFAdd(ctx context.Context, key string, els ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(els))
	args[0] = key
	args = appendArgs(args, els)
	reply, err = _client.Do(ctx, "pfadd", args...)
	return
}

// PFCount .
// Redis Pfcount 命令返回给定 HyperLogLog 的基数估算值。
// 整数，返回给定 HyperLogLog 的基数值，如果多个 HyperLogLog 则返回基数估值之和。
func PFCount(ctx context.Context, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	reply, err = _client.Do(ctx, "pfcount", args...)
	return
}

// PFMerge .
// Redis PFMERGE 命令将多个 HyperLogLog 合并为一个 HyperLogLog ，
// 合并后的 HyperLogLog 的基数估算值是通过对所有 给定 HyperLogLog 进行并集计算得出的。
func PFMerge(ctx context.Context, dest string, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(keys))
	args[0] = dest
	for i, key := range keys {
		args[1+i] = key
	}
	reply, err = _client.Do(ctx, "pfmerge", args...)
	return
}

//------------------------------------------------------------------------------

// BgRewriteAOF .
// Redis Bgrewriteaof 命令用于异步执行一个 AOF（AppendOnly File） 文件重写操作。重写会创建一个当前 AOF 文件的体积优化版本。
// 即使 Bgrewriteaof 执行失败，也不会有任何数据丢失，因为旧的 AOF 文件在 Bgrewriteaof 成功之前不会被修改。
// 注意：从 Redis 2.4 开始， AOF 重写由 Redis 自行触发， BGREWRITEAOF 仅仅用于手动触发重写操作。
func BgRewriteAOF(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "bgrewriteaof")
	return
}

// BgSave .
// Redis Bgsave 命令用于在后台异步保存当前数据库的数据到磁盘。
// BGSAVE 命令执行之后立即返回 OK ，
// 然后 Redis fork 出一个新子进程，原来的 Redis 进程(父进程)继续处理客户端请求，而子进程则负责将数据保存到磁盘，然后退出。
func BgSave(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "bgsave")
	return
}

// ClientKill .
// Redis Client Kill 命令用于关闭客户端连接。
func ClientKill(ctx context.Context, ipPort string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "kill", ipPort)
	return
}

// ClientKillByFilter is new style synx, while the ClientKill is old
// CLIENT KILL <option> [value] ... <option> [value]
func ClientKillByFilter(ctx context.Context, keys ...string) (reply interface{}, err error) {
	args := make([]interface{}, 1+len(keys))
	args[0] = "kill"
	for i, key := range keys {
		args[1+i] = key
	}
	reply, err = _client.Do(ctx, "client", args...)
	return
}

// ClientList .
// Redis Client List 命令用于返回所有连接到服务器的客户端信息和统计数据。
func ClientList(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "list")
	return
}

// ClientPause .
// Redis Client Pause 命令用于阻塞客户端命令一段时间（以毫秒计）。
func ClientPause(ctx context.Context, dur time.Duration) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "pause", formatMs(dur))
	return
}

// ClientID .
// (error) ERR Syntax error, try CLIENT (LIST | KILL | GETNAME | SETNAME | PAUSE | REPLY)
func ClientID(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "id")
	return
}

// ClientUnblock .
// (error) ERR Syntax error, try CLIENT (LIST | KILL | GETNAME | SETNAME | PAUSE | REPLY)
func ClientUnblock(ctx context.Context, id int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "unblock", id)
	return
}

// ClientUnblockWithError .
// (error) ERR Syntax error, try CLIENT (LIST | KILL | GETNAME | SETNAME | PAUSE | REPLY)
func ClientUnblockWithError(ctx context.Context, id int64) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "unblock", id, "error")
	return
}

// ClientSetName .
// Redis Client Setname 命令用于指定当前连接的名称。
// 这个名字会显示在 CLIENT LIST 命令的结果中， 用于识别当前正在与服务器进行连接的客户端。
func ClientSetName(ctx context.Context, name string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "client", "setname", name)
	return
}

// ConfigGet .
// Redis Config Get 命令用于获取 redis 服务的配置参数。
// 在 Redis 2.4 版本中， 有部分参数没有办法用 CONFIG GET 访问，但是在最新的 Redis 2.6 版本中，所有配置参数都已经可以用 CONFIG GET 访问了。
func ConfigGet(ctx context.Context, parameter string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "config", "get", parameter)
	return
}

// ConfigResetStat .
// Redis Config Resetstat 命令用于重置 INFO 命令中的某些统计数据，包括：
// Keyspace hits (键空间命中次数)
// Keyspace misses (键空间不命中次数)
// Number of commands processed (执行命令的次数)
// Number of connections received (连接服务器的次数)
// Number of expired keys (过期key的数量)
// Number of rejected connections (被拒绝的连接数量)
// Latest fork(2) time(最后执行 fork(2) 的时间)
// The aof_delayed_fsync counter(aof_delayed_fsync 计数器的值)
func ConfigResetStat(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "config", "resetstat")
	return
}

// ConfigSet .
// Redis Config Set 命令可以动态地调整 Redis 服务器的配置(configuration)而无须重启。
// 你可以使用它修改配置参数，或者改变 Redis 的持久化(Persistence)方式。
func ConfigSet(ctx context.Context, parameter, value string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "config", "set", parameter, value)
	return
}

// ConfigRewrite .
// Redis Config rewrite 命令对启动 Redis 服务器时所指定的 redis.conf 配置文件进行改写。
// CONFIG SET 命令可以对服务器的当前配置进行修改，
// 而修改后的配置可能和 redis.conf 文件中所描述的配置不一样，
// CONFIG REWRITE 的作用就是通过尽可能少的修改，
// 将服务器当前所使用的配置记录到 redis.conf 文件中。
func ConfigRewrite(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "config", "rewrite")
	return
}

// DbSize .
// Redis Dbsize 命令用于返回当前数据库的 key 的数量。
func DbSize(ctx context.Context) (reply interface{}, err error) {
	reply, err = DBSize(ctx)
	return
}

// DBSize .
// Redis Dbsize 命令用于返回当前数据库的 key 的数量。
func DBSize(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "dbsize")
	return
}

// FlushAll .
// Redis Flushall 命令用于清空整个 Redis 服务器的数据(删除所有数据库的所有 key )。
// Redis 4.0 同样给这两个指令也带来了异步化，在指令后面增加 async 参数扔给后台线程销毁，不会阻塞当前线程。
func FlushAll(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "flushall")
	return
}

// FlushAllAsync .
// Redis Flushall 命令用于清空整个 Redis 服务器的数据(删除所有数据库的所有 key )。
// Redis 4.0 同样给这两个指令也带来了异步化，在指令后面增加 async 参数扔给后台线程销毁，不会阻塞当前线程。
func FlushAllAsync(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "flushall", "async")
	return
}

// FlushDb .
// Redis Flushdb 命令用于清空当前数据库中的所有 key。
// Redis 4.0 同样给这两个指令也带来了异步化，在指令后面增加 async 参数扔给后台线程销毁，不会阻塞当前线程。
func FlushDb(ctx context.Context) (reply interface{}, err error) {
	reply, err = FlushDB(ctx)
	return
}

// FlushDB .
// Redis Flushdb 命令用于清空当前数据库中的所有 key。
// Redis 4.0 同样给这两个指令也带来了异步化，在指令后面增加 async 参数扔给后台线程销毁，不会阻塞当前线程。
func FlushDB(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "flushdb")
	return
}

// FlushDBAsync .
// Redis Flushdb 命令用于清空当前数据库中的所有 key。
// Redis 4.0 同样给这两个指令也带来了异步化，在指令后面增加 async 参数扔给后台线程销毁，不会阻塞当前线程。
func FlushDBAsync(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "flushdb", "async")
	return
}

// Info .
// Redis Info 命令以一种易于理解和阅读的格式，返回关于 Redis 服务器的各种信息和统计数值。
// 通过给定可选的参数 section ，可以让命令只返回某一部分的信息：
// server : 一般 Redis 服务器信息
// clients : 已连接客户端信息
// memory : 内存信息
// persistence : RDB 和 AOF 的相关信息
// stats : 一般统计信息
// replication : 主/从复制信息
// cpu : CPU 计算量统计信息
// commandstats : Redis 命令统计信息
// cluster : Redis 集群信息
// keyspace : 数据库相关的统计信息
func Info(ctx context.Context, section ...string) (reply interface{}, err error) {
	var args []interface{}
	if len(section) > 0 {
		args = append(args, section[0])
	}
	reply, err = _client.Do(ctx, "info", args...)
	return
}

// LastSave .
// Redis Lastsave 命令返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示。
func LastSave(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "lastsave")
	return
}

// Save .
// Redis Save 命令执行一个同步保存操作，将当前 Redis 实例的所有数据快照(snapshot)以 RDB 文件的形式保存到硬盘。
func Save(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "save")
	return
}

// shutdown .
func shutdown(ctx context.Context, modifier string) (reply interface{}, err error) {
	if modifier == "" {
		reply, err = _client.Do(ctx, "shutdown")
	} else {
		reply, err = _client.Do(ctx, "shutdown", modifier)
	}
	if err != nil {
		if err == io.EOF {
			// Server quit as expected.
			err = nil
		}
	} else {
		// Server did not quit. String reply contains the reason.
		err = errors.New(fmt.Sprint(reply))
	}
	return
}

// Shutdown .
// Redis Shutdown 命令执行以下操作：
// 停止所有客户端
// 如果有至少一个保存点在等待，执行 SAVE 命令
// 如果 AOF 选项被打开，更新 AOF 文件
// 关闭 redis 服务器(server)
// 执行失败时返回错误。 执行成功时不返回任何信息，服务器和客户端的连接断开，客户端自动退出。
func Shutdown(ctx context.Context) (reply interface{}, err error) {
	reply, err = shutdown(ctx, "")
	return
}

// ShutdownSave .
// Redis Shutdown 命令执行以下操作：
// 停止所有客户端
// 如果有至少一个保存点在等待，执行 SAVE 命令
// 如果 AOF 选项被打开，更新 AOF 文件
// 关闭 redis 服务器(server)
// 执行失败时返回错误。 执行成功时不返回任何信息，服务器和客户端的连接断开，客户端自动退出。
func ShutdownSave(ctx context.Context) (reply interface{}, err error) {
	reply, err = shutdown(ctx, "save")
	return
}

// ShutdownNoSave .
// Redis Shutdown 命令执行以下操作：
// 停止所有客户端
// 如果有至少一个保存点在等待，执行 SAVE 命令
// 如果 AOF 选项被打开，更新 AOF 文件
// 关闭 redis 服务器(server)
// 执行失败时返回错误。 执行成功时不返回任何信息，服务器和客户端的连接断开，客户端自动退出。
func ShutdownNoSave(ctx context.Context) (reply interface{}, err error) {
	reply, err = shutdown(ctx, "nosave")
	return
}

// SlaveOf .
// Redis Slaveof 命令可以将当前服务器转变为指定服务器的从属服务器(slave server)。
// 如果当前服务器已经是某个主服务器(master server)的从属服务器，
// 那么执行 SLAVEOF host port 将使当前服务器停止对旧主服务器的同步，丢弃旧数据集，转而开始对新主服务器进行同步。
// 另外，对一个从属服务器执行命令 SLAVEOF NO ONE 将使得这个从属服务器关闭复制功能，并从从属服务器转变回主服务器，原来同步所得的数据集不会被丢弃。
// 利用『 SLAVEOF NO ONE 不会丢弃同步所得数据集』这个特性，可以在主服务器失败的时候，将从属服务器用作新的主服务器，从而实现无间断运行。
func SlaveOf(ctx context.Context, host, port string) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "slaveof", host, port)
	return
}

// Time .
// Redis Time 命令用于返回当前服务器时间。
func Time(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "time")
	return
}

// SlowLog .
func SlowLog() {
	panic("not implemented")
}

// Sync .
// Redis Sync 命令用于同步主从服务器。
func Sync(ctx context.Context) (reply interface{}, err error) {
	reply, err = _client.Do(ctx, "sync")
	return
}

//------------------------------------------------------------------------------
