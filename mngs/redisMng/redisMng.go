package redisMng

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"time"
)

var ctx = context.Background()

//var pool redis.PoolStats

var client *redis.Client

type RedisMng struct {
	//config configStruct.Redis
	Conn redis.Conn
}

// Init 初始化
func Init(redisC *configStruct.RedisConfig) (err error) {

	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	redisURL := redisC.Host + ":" + redisC.Port
	log.Println("【redis-dsn】", redisURL)
	client = redis.NewClient(&redis.Options{
		Addr:         redisURL,
		Username:     redisC.Username,
		Password:     redisC.Password,
		DB:           0,
		DialTimeout:  time.Duration(redisC.IdleTimeout),
		MinIdleConns: redisC.MaxIdle,
		PoolSize:     redisC.MaxActive,
	})
	log.Println("redis Init", client.Context().Err())
	log.Println("redis PoolStats", client.PoolStats())
	ping := client.Ping(ctx)
	log.Println("ping", ping)
	log.Println("\tping.Err()", ping.Err())
	err = ping.Err()
	return
}

// NewRedisMng 返回一个redis管理器实例
func NewRedisMng() *RedisMng {
	redisMng := RedisMng{}
	return &redisMng
}

// GetString 读取指定键的字符串值
func (mng *RedisMng) GetString(key string) (val string, err error) {
	return client.Get(ctx, key).Result()
	//  err == redis.Nil 不存在
}

// Set 设置键值
func (mng *RedisMng) Set(key string, value interface{}, expire time.Duration) (err error) {
	return client.Set(ctx, key, value, expire).Err()
}

// -------BEGIN------哈希相关的操作-----BEGIN--------

// HGetString 集合获取
func (mng *RedisMng) HGetString(key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

// HSet 设置Hash
func (mng *RedisMng) HSet(key string, value interface{}) (int64, error) {

	return client.HSet(ctx, key, value).Result()
}

// HSetNX 设置Hash一个file
func (mng *RedisMng) HSetNX(key, field string, value interface{}) (bool, error) {

	return client.HSetNX(ctx, key, field, value).Result()
}

// HDel 删除hash里的key
func (mng *RedisMng) HDel(key string, fields []string) (int64, error) {

	return client.HDel(ctx, key, fields...).Result()

}

// HExists 判断hash中有没有这个field
func (mng *RedisMng) HExists(key, field string) (bool, error) {
	return client.HExists(ctx, key, field).Result()

}

// HKeys 获取hash中所有的field
func (mng *RedisMng) HKeys(key string) ([]string, error) {
	return client.HKeys(ctx, key).Result()
}

// HIncrBy 增加hash中的字段的值  返回的是字段被修改过后的值
func (mng *RedisMng) HIncrBy(key, field string, increase int64) (int64, error) {
	return client.HIncrBy(ctx, key, field, increase).Result()
}

// HLen  hash 中 一个key下的数量
func (mng *RedisMng) HLen(key string) (int64, error) {

	return client.HLen(ctx, key).Result()

}

//
//// -------END------哈希相关的操作----END---------
//
////-------BEGIN-----字符串相关操作-------BEGIN-------

//func GetFromDB(key string, dbStruct int) string {
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	res, err := client.Do("SELECT", dbStruct)
//
//	if err != nil {
//
//		return ""
//	}
//
//	res, _ = client.Do("GET", key)
//
//	if res == nil {
//
//		return ""
//	}
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return string(res.([]byte))
//
//}
//

//func SetToDB(key string, value interface{}, expire interface{}, dbStruct int) error {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	_, err := client.Do("SELECT", dbStruct)
//
//	if err != nil {
//
//		fmt.Println("失败", err)
//
//		return errors.New("切换数据库失败")
//
//	}
//
//	client.Do("SET", key, value)
//
//	if expire != nil {
//
//		client.Do("EXPIRE", key, expire)
//
//		return nil
//	}
//
//	return nil
//}
//
////设置新的值 返回久的值
//func GETSET(keyName string, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("set", keyName, value))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////设置新的值 返回久的值
//func SETEX(keyName string, expire interface{}, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("setex", keyName, expire, value)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////设置新的值 返回久的值  返回 1 成功 0  设置失败
//func SETNX(keyName string, expire interface{}, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("setnx", keyName, value)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////将 key 中储存的数字值增一
//func INCR(keyName string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("incr", keyName)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////将 key 中储存的数字值减一
//func DECR(keyName string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("decr", keyName)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////将 key 中储存的数字值增量
//func INCRBY(keyName string, increment interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("incrby", keyName, increment)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////将 key 中储存的数字值减去一定的量
//func DECRBY(keyName string, increment interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("decrby", keyName, increment)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////将 key 设置过期时间
//func EXPIRE(keyName string, expire interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("expire", keyName, expire)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////--------END-----字符串相关操作-----END---------
//
////--------BEGIN-----列表相关操作-----BEGIN---------
//
////Redis Lindex 命令用于通过索引获取列表中的元素。你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
//func LINDEX(keyName string, position interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("lindex", keyName, position))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//
//}
//
////一个列表的长度
//func LLEN(keyName string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int64(client.Do("llen", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////从列表中左边删除第一个元素 ，返回这个元素的值
//func LPOP(keyName string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("lpop", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////从列表中左边删除第一个元素 ，返回这个元素的值
//func RPOP(keyName string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("rpop", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////添加元素 如果失败 0  成功  返回列表的长度
//func RPUSH(keyName string, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("rpush", keyName, value))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////列表的长度 0 -1 全部
//func LRANGE(keyName string, start, end interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("lrange", keyName, start, end))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////--------END-----列表相关操作-----END---------
//
////--------BEGIN-----有序集合相关操作-----BEGIN---------
//
////向有序集合添加成员 被成功添加的新成员的数量，不包括那些被更新的、已经存在的成员。
//func ZADD(keyName string, score interface{}, member interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("zadd", keyName, score, member)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////Redis Zcard 命令用于计算集合中元素的数量。 返回这个集合的个数
//func ZCARD(keyName string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zcard", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////Zcount 命令用于计算有序集合中指定分数区间的成员数量。
//func ZCOUNT(keyName string, min interface{}, end interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zcount", keyName, min, end))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////ZINCRBY 命令对有序集合中指定成员的分数加上增量 increment  member 成员的新分数值，以字符串形式表示。
//func ZINCRBY(keyName string, increment interface{}, field string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zincrby", keyName, increment, field))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////Zrange 返回有序集中，指定区间内的成员。
//func ZRANGE(keyName string, start_index, end_index interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("zrevrange", keyName, start_index, end_index))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////ZREM 删除一个有序集合中的指定成员  被成功移除的成员的数量，不包括被忽略的成员。
//func ZREM(keyName string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zrem", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////Zscore 命令返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。 成员的分数值，以字符串形式表示。
//func ZSCORE(keyName, member string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("zscore", keyName, member))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////--------END-----有序集合相关操作-----END---------
//
////--------BEGIN-----无序集合相关操作-----BEGIN---------
//
////Sadd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
//func SADD(keyName string, value interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("zscore", keyName, value))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////Scard 命令返回集合中元素的数量。
//func SCARD(keyName string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("scard", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////Sdiff 命令返回给定集合之间的差集。不存在的集合 key 将视为空集。 包含差集成员的列表
//func SDIFF(key1, key2 string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("sdiff", key1, key2))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
//// Sinter 命令返回给定所有给定集合的交集
//func SINTER(key1, key2 string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("sinter", key1, key2))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////Smembers 命令返回集合中的所有的成员。 不存在的集合 key 被视为空集合。
//func SMEMBERS(keyName string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("smembers", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////Smove 命令将指定成员 member 元素从 source 集合移动到 destination 集合。
////如果成员元素被成功移除，返回 1 。 如果成员元素不是 source 集合的成员，并且没有任何操作对 destination 集合执行，那么返回 0 。
//func SMOVE(set1, set2 string, member interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("smove", set1, set2, member))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////Srem 命令用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略。
//func SREM(keyName string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("srem", keyName))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////
//
////Srandmember 命令用于返回集合中的一个随机元素。
//// 只提供集合 key 参数时，返回一个元素；如果集合为空，返回 nil 。 如果提供了 count 参数，那么返回一个数组；如果集合为空，返回空数组。
//func SRANDMEMBER(keyName string, count interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	if count == nil {
//
//		result, _ := redis.String(client.Do("srandmember", keyName))
//
//		return result
//
//	} else {
//
//		result, _ := redis.Strings(client.Do("srandmember", keyName, count))
//
//		return result
//	}
//
//}
//
////--------END-----无序集合相关操作-----END---------
//
////给一个集合中的field增加权重  然后根据获取这个集合
//func InsertSortedSet(setName string, score int, data interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("zincrby", setName, score, data)
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
//
////查询sorted set 有序集合
//func ZrevangeSortedSet(setName string, start_index, end_index int) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("zrevrange", setName, start_index, end_index))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//
//}
