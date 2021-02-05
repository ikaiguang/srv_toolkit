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

// BitCountParam .
type BitCountParam struct {
	Start, End int64
}

// BitCount .
// 统计指定位区间上值为1的个数
// BITCOUNT key [start end]
// 从左向右从0开始，从右向左从-1开始，注意start和end是字节
func BitCount(ctx context.Context, key string, bitCount *BitCountParam) (reply interface{}, err error) {
	args := []interface{}{key}
	if bitCount != nil {
		args = append(
			args,
			bitCount.Start,
			bitCount.End,
		)
	}
	reply, err = tkredis.Redis().Do(ctx, "bitcount", args...)
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
	reply, err = tkredis.Redis().Do(ctx, "bitop", args...)
	return
}

// BitOpOr .
// BITOP OR destkey key [key ...] ，对一个或多个key求逻辑或，并将结果保存到destkey
func BitOpOr(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("or", destKey, keys...)
	reply, err = tkredis.Redis().Do(ctx, "bitop", args...)
	return
}

// BitOpXor .
// BITOP XOR destkey key [key ...] ，对一个或多个key求逻辑异或，并将结果保存到destkey
func BitOpXor(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("xor", destKey, keys...)
	reply, err = tkredis.Redis().Do(ctx, "bitop", args...)
	return
}

// BitOpNot .
// BITOP NOT destkey key ，对给定key求逻辑非，并将结果保存到destkey
func BitOpNot(ctx context.Context, destKey string, keys ...string) (reply interface{}, err error) {
	args := bitOp("not", destKey, keys...)
	reply, err = tkredis.Redis().Do(ctx, "bitop", args...)
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
	reply, err = tkredis.Redis().Do(ctx, "bitpos", args...)
	return
}

// Decr .
// Redis Decr 命令将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func Decr(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "decr", key)
	return
}

// DecrBy .
// Redis Decrby 命令将 key 所储存的值减去指定的减量值。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func DecrBy(ctx context.Context, key string, decrement int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "decrby", key, decrement)
	return
}

// Get .
// Redis Get 命令用于获取指定 key 的值。如果 key 不存在，返回 nil 。如果key 储存的值不是字符串类型，返回一个错误。
func Get(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "get", key)
	return
}

// GetBit .
// Redis Getbit 命令用于对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
func GetBit(ctx context.Context, key string, offset int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "getbit", key, offset)
	return
}

// GetRange .
// Redis Getrange 命令用于获取存储在指定 key 中字符串的子字符串。字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
func GetRange(ctx context.Context, key string, start, end int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "getrange", key, start, end)
	return
}

// GetSet .
// Redis Getset 命令用于设置指定 key 的值，并返回 key 的旧值。
func GetSet(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "getset", key, value)
	return
}

// Incr .
// Redis Incr 命令将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func Incr(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "incr", key)
	return
}

// IncrBy .
// Redis Incrby 命令将 key 中储存的数字加上指定的增量值。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func IncrBy(ctx context.Context, key string, value int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "incrby", key, value)
	return
}

// IncrByFloat .
// Redis Incrbyfloat 命令为 key 中所储存的值加上指定的浮点数增量值。
// 如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
func IncrByFloat(ctx context.Context, key string, value float64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "incrbyfloat", key, value)
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
	reply, err = tkredis.Redis().Do(ctx, "mget", keys...)
	return
}

// MSet .
// Redis Mset 命令用于同时设置一个或多个 key-value 对。
func MSet(ctx context.Context, pairs ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, len(pairs))
	args = appendArgs(args, pairs)
	reply, err = tkredis.Redis().Do(ctx, "mset", args...)
	return
}

// MSetNX .
// Redis Msetnx 命令用于所有给定 key 都不存在时，同时设置一个或多个 key-value 对。
// 当所有 key 都成功设置，返回 1 。 如果所有给定 key 都设置失败(至少有一个 key 已经存在)，那么返回 0 。
func MSetNX(ctx context.Context, pairs ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, len(pairs))
	args = appendArgs(args, pairs)
	reply, err = tkredis.Redis().Do(ctx, "msetnx", args...)
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
	reply, err = tkredis.Redis().Do(ctx, "set", args...)
	return
}

// SetBit .
// Redis Setbit 命令用于对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
func SetBit(ctx context.Context, key string, offset int64, value int) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "setbit", key, offset, value)
	return
}

// SetNX .
// Redis Setnx（SET if Not eXists） 命令在指定的 key 不存在时，为 key 设置指定的值。
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (reply interface{}, err error) {
	if expiration == 0 {
		// Use old `SETNX` to support old Redis versions.
		reply, err = tkredis.Redis().Do(ctx, "setnx", key, value)
	} else {
		if usePrecise(expiration) {
			reply, err = tkredis.Redis().Do(ctx, "set", key, value, "px", formatMs(expiration), "nx")
		} else {
			reply, err = tkredis.Redis().Do(ctx, "set", key, value, "ex", formatMs(expiration), "nx")
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
		reply, err = tkredis.Redis().Do(ctx, "set", key, value, "xx")
	} else {
		if usePrecise(expiration) {
			reply, err = tkredis.Redis().Do(ctx, "set", key, value, "px", formatMs(expiration), "xx")
		} else {
			reply, err = tkredis.Redis().Do(ctx, "set", key, value, "ex", formatSec(expiration), "xx")
		}
	}
	return
}

// SetRange .
// Redis Setrange 命令用指定的字符串覆盖给定 key 所储存的字符串值，覆盖的位置从偏移量 offset 开始。
func SetRange(ctx context.Context, key string, offset int64, value string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "setrange", key, offset, value)
	return
}

// StrLen .
// Redis Strlen 命令用于获取指定 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误。
func StrLen(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "strlen", key)
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
	reply, err = tkredis.Redis().Do(ctx, "hdel", args...)
	return
}

// HExists .
// Redis Hexists 命令用于查看哈希表的指定字段是否存在。
func HExists(ctx context.Context, key string, field string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hexists", key, field)
	return
}

// HGet .
// Redis Hget 命令用于返回哈希表中指定字段的值。
func HGet(ctx context.Context, key string, field string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hget", key, field)
	return
}

// HGetAll .
// Redis Hgetall 命令用于返回哈希表中，所有的字段和值。
// 在返回值里，紧跟每个字段名(field name)之后是字段的值(value)，所以返回值的长度是哈希表大小的两倍
func HGetAll(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hgetall", key)
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
	reply, err = tkredis.Redis().Do(ctx, "hincrby", key, field, incr)
	return
}

// HIncrByFloat .
// Redis Hincrbyfloat 命令用于为哈希表中的字段值加上指定浮点数增量值。
// 如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0
func HIncrByFloat(ctx context.Context, key, field string, incr float64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hincrbyfloat", key, field, incr)
	return
}

// HKeys .
// Redis Hkeys 命令用于获取哈希表中的所有域（field）。
func HKeys(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hkeys", key)
	return
}

// HLen .
// Redis Hlen 命令用于获取哈希表中字段的数量。
func HLen(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hlen", key)
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
	reply, err = tkredis.Redis().Do(ctx, "hmget", args...)
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
	reply, err = tkredis.Redis().Do(ctx, "hmset", args...)
	return
}

// HSet .
// Redis Hset 命令用于为哈希表中的字段赋值 。
// 如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果字段已经存在于哈希表中，旧值将被覆盖。
func HSet(ctx context.Context, key, field string, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hset", key, field, value)
	return
}

// HSetNX .
// Redis Hsetnx 命令用于为哈希表中不存在的的字段赋值 。
// 如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果字段已经存在于哈希表中，操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func HSetNX(ctx context.Context, key, field string, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hsetnx", key, field, value)
	return
}

// HVals .
// Redis Hvals 命令返回哈希表所有域(field)的值。
func HVals(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "hvals", key)
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
	reply, err = tkredis.Redis().Do(ctx, "blpop", args...)
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
	reply, err = tkredis.Redis().Do(ctx, "brpop", args...)
	return
}

// BRPopLPush .
// Redis Brpoplpush 命令从列表中取出最后一个元素，并插入到另外一个列表的头部； 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
// 假如在指定时间内没有任何元素被弹出，则返回一个 nil 和等待时长。 反之，返回一个含有两个元素的列表，第一个元素是被弹出元素的值，第二个元素是等待时长。
func BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "brpoplpush",
		source,
		destination,
		formatSec(timeout),
	)
	return
}

// LIndex .
// Redis Lindex 命令用于通过索引获取列表中的元素。你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func LIndex(ctx context.Context, key string, index int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "lindex", key, index)
	return
}

// LInsert .
// Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
// 当列表不存在时，被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// LINSERT key BEFORE|AFTER pivot value
func LInsert(ctx context.Context, key, op string, pivot, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "linsert", key, op, pivot, value)
	return
}

// LInsertBefore .
// Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
// 当列表不存在时，被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// LINSERT key BEFORE|AFTER pivot value
func LInsertBefore(ctx context.Context, key, op string, pivot, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "linsert", key, "before", pivot, value)
	return
}

// LInsertAfter .
// Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
// 当列表不存在时，被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// LINSERT key BEFORE|AFTER pivot value
func LInsertAfter(ctx context.Context, key, op string, pivot, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "linsert", key, "after", pivot, value)
	return
}

// LLen .
// Redis Llen 命令用于返回列表的长度。 如果列表 key 不存在，则 key 被解释为一个空列表，返回 0 。 如果 key 不是列表类型，返回一个错误。
func LLen(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "llen", key)
	return
}

// LPop .
// Redis Lpop 命令用于移除并返回列表的第一个元素。
// 列表的第一个元素。 当列表 key 不存在时，返回 nil 。
func LPop(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "lpop", key)
	return
}

// LPush .
// Redis Lpush 命令将一个或多个值插入到列表头部。 如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。 当 key 存在但不是列表类型时，返回一个错误。
// 注意：在Redis 2.4版本以前的 LPUSH 命令，都只接受单个 value 值。
func LPush(ctx context.Context, key string, values ...interface{}) (reply interface{}, err error) {
	args := make([]interface{}, 1, 1+len(values))
	args[0] = key
	args = appendArgs(args, values)
	reply, err = tkredis.Redis().Do(ctx, "lpush", args...)
	return
}

// LPushX .
// Redis Lpushx 将一个值插入到已存在的列表头部，列表不存在时操作无效。
func LPushX(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "lpushx", key, value)
	return
}

// LRange .
// Redis Lrange 返回列表中指定区间内的元素，区间以偏移量 START 和 END 指定。
// 其中 0 表示列表的第一个元素， 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func LRange(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "lrange", key, start, stop)
	return
}

// LRem .
// Redis Lrem 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素。
// COUNT 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
// count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
// count = 0 : 移除表中所有与 VALUE 相等的值。
func LRem(ctx context.Context, key string, count int64, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "lrem", key, count, value)
	return
}

// LSet .
// Redis Lset 通过索引来设置元素的值。
// 当索引参数超出范围，或对一个空列表进行 LSET 时，返回一个错误。
// 关于列表下标的更多信息，请参考 LINDEX 命令。
func LSet(ctx context.Context, key string, index int64, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "lset", key, index, value)
	return
}

// LTrim .
// Redis Ltrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 下标 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func LTrim(ctx context.Context, key string, start, stop int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "ltrim", key, start, stop)
	return
}

// RPop .
// Redis Rpop 命令用于移除列表的最后一个元素，返回值为移除的元素。
func RPop(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "rpop", key)
	return
}

// RPopLPush .
// Redis Rpoplpush 命令用于移除列表的最后一个元素，并将该元素添加到另一个列表并返回。
func RPopLPush(ctx context.Context, source, destination string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "rpoplpush", source, destination)
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
	reply, err = tkredis.Redis().Do(ctx, "rpush", args...)
	return
}

// RPushX .
// Redis rpushx，命令用于将一个或多个值插入到已存在的列表尾部(最右边)，如果列表不存在，操作无效。
func RPushX(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "rpushx", key, value)
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
	reply, err = tkredis.Redis().Do(ctx, "sadd", args...)
	return
}

// SCard .
// Redis Scard 命令返回集合中元素的数量。
func SCard(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "scard", key)
	return
}

// SDiff .
// Redis Sdiff 命令返回第一个集合与其他集合之间的差异，也可以认为说第一个集合中独有的元素。不存在的集合 key 将视为空集。
// 差集的结果来自前面的 FIRST_KEY ,而不是后面的 OTHER_KEY1，也不是整个 FIRST_KEY OTHER_KEY1..OTHER_KEYN 的差集。
func SDiff(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "sdiff", keys...)
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
	reply, err = tkredis.Redis().Do(ctx, "sdiffstore", args...)
	return
}

// SInter .
// Redis Sinter 命令返回给定所有给定集合的交集。 不存在的集合 key 被视为空集。 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
func SInter(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "sinter", keys...)
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
	reply, err = tkredis.Redis().Do(ctx, "sinterstore", args...)
	return
}

// SIsMember .
// Redis Sismember 命令判断成员元素是否是集合的成员。
func SIsMember(ctx context.Context, key string, member interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "sismember", key, member)
	return
}

// SMembers .
// Redis Smembers 命令返回集合中的所有的成员。 不存在的集合 key 被视为空集合。
func SMembers(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "smembers", key)
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
	reply, err = tkredis.Redis().Do(ctx, "smove", source, destination, member)
	return
}

// SPop .
// Redis Spop 命令用于移除集合中的指定 key 的一个或多个随机元素，移除后会返回移除的元素。
// 该命令类似 Srandmember 命令，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SPop(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "spop", key)
	return
}

// SPopN .
// Redis Spop 命令用于移除集合中的指定 key 的一个或多个随机元素，移除后会返回移除的元素。
// 该命令类似 Srandmember 命令，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SPopN(ctx context.Context, key string, count int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "spop", key, count)
	return
}

// SRandMember .
// Redis Srandmember 命令用于返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， Srandmember 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SRandMember(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "srandmember", key)
	return
}

// SRandMemberN .
// Redis Srandmember 命令用于返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， Srandmember 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
func SRandMemberN(ctx context.Context, key string, count int64) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "srandmember", key, count)
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
	reply, err = tkredis.Redis().Do(ctx, "srem", args...)
	return
}

// SUnion .
// Redis Sunion 命令返回给定集合的并集。不存在的集合 key 被视为空集。
func SUnion(ctx context.Context, keys ...interface{}) (reply interface{}, err error) {
	reply, err = tkredis.Redis().Do(ctx, "sunion", keys...)
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
	reply, err = tkredis.Redis().Do(ctx, "sunionstore", args...)
	return
}

//------------------------------------------------------------------------------
