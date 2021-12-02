package redisMng

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"time"
)

var pool redis.Pool

/**
 * @func：实例化一个redis连接池
 * @author Hank
 * @date   2019-02-24
 */
type RedisMng struct {
	//config configStruct.Redis
	Conn redis.Conn
}

func Init(redisC *configStruct.RedisConfig) (err error){
	pool = redis.Pool{
		MaxActive:   redisC.MaxActive,
		MaxIdle:     redisC.MaxIdle,
		IdleTimeout: time.Duration(redisC.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			redisURL := redisC.Host + ":" + redisC.Port
			log.Println("【redis-dsn】",redisURL)
			var conn redis.Conn
			conn, err = redis.Dial("tcp", redisURL)
			if err != nil {
				fmt.Println("【redis-dial-err】", err)
				return nil, err
			}

			if _, err := conn.Do("AUTH", redisC.Password); err != nil {
				fmt.Println("【redis-auth-err】", err)
				conn.Close()
				return nil, err
			}
			return conn, err
		},
	}
	return
}

/**
 * @func:   NewRedisMng 返回一个redis管理器实例
 * @author: Wiidz
 * @date:   2020-04-15
 */
func NewRedisMng() *RedisMng {
	redisMng := RedisMng{}
	return &redisMng
}

/**
 * @func:   Get 读取指定键的值
 * @author: Wiidz
 * @date:   2020-04-15
 * @Param: [key] 键名
 * @return: [str] 键值
 */
func (mng *RedisMng) GetString(key string) (string, error) {

	//【1】取出一条连接
	rc := pool.Get()
	defer rc.Close()

	//【2】读取值
	res, err := redis.String(rc.Do("GET", key))
	if err != nil && err.Error() == "redigo: nil returned" {
		return "",nil
	}

	//【3】返回
	return res, err
}

/**
 * @func:   Get 读取指定键的值
 * @author: Wiidz
 * @date:   2020-04-15
 * @Param: [key] 键名
 * @return: [str] 键值
 */
func (mng *RedisMng) Get(key string) (interface{}, error) {

	//【1】取出一条连接
	rc := pool.Get()
	defer rc.Close()

	//【2】读取值
	res, err := rc.Do("GET", key)
	if err != nil && err.Error() == "redigo: nil returned" {
		return "",nil
	}

	//【3】返回
	return res, err
}

/**
 * @func:   Set 设置键值
 * @author: Wiidz
 * @date:   2020-04-15
 * @Param: [key] 键名
 * 		   [value] 键值
 *         [expire] 过期时间（秒）
 */
func (mng *RedisMng) Set(key string, value interface{}, expire int) error {

	//【1】取出一条连接
	rc := pool.Get()
	defer rc.Close()

	//【2】获取redis连接
	_, err := rc.Do("SET", key, value)

	//【3】判断结果
	if err != nil {
		return err
	}

	//【4】设置过期时间
	if expire != 0 {
		_, err = rc.Do("EXPIRE", key, expire)
		if err != nil {
			return err
		}
	}

	//【5】返回
	return nil

}

// -------BEGIN------哈希相关的操作-----BEGIN--------
/**
* @func: HSGet 设置key
* @author Hank
* @date   2019-06-21
 */
func (mng *RedisMng) HGet(key_name, field string) (interface{}, error) {

	rc := pool.Get()
	defer rc.Close()

	result, err := redis.String(rc.Do("hget", key_name, field))
	return result, err
}

/**
* @func: HSET 设置key
* @author Hank
* @date   2019-06-21
 */
func (mng *RedisMng) HSet(key_name, field string, value interface{}) (interface{}, error) {

	rc := pool.Get()
	defer rc.Close()

	result, err := rc.Do("hset", key_name, field, value)
	return result, err
}

/**
* @func: HDEL 删除hash里的key
* @author Hank
* @date   2019-06-21
 */
func (mng *RedisMng) HDel(key_name string, fields ...string) (interface{}, error) {

	rc := pool.Get()
	defer rc.Close()

	result, err := rc.Do("hdel", key_name, fields)
	return result, err

}

/**
* @func: HEXISTS 判断hashkey中有没有这个字段
* @author Hank
* @date   2019-06-21
 */

func (mng *RedisMng) HExists(key_name, field_name string) (interface{}, error) {
	rc := pool.Get()
	defer rc.Close()

	result, err := rc.Do("hexists", key_name, field_name)
	return result, err
}

//
///**
// * @func: HKEYS 获取hash中所有的field
// * @author Hank
// * @date   2019-06-21
// */
//
func (mng *RedisMng) HKeys(key_name string) (interface{}, error) {

	rc := pool.Get()
	defer rc.Close()

	result, err := redis.Strings(rc.Do("hkeys", key_name))
	return result, err

}

//
///**
// * @func: HINCRBY 增加hash中的字段的值  返回的是字段被修改过后的值
// * @author Hank
// * @date   2019-06-21
// */
//
func (mng *RedisMng) HIncrby(key_name, field_name, incr_by_number interface{}) (res interface{}, err error) {

	rc := pool.Get()
	defer rc.Close()

	result, err := rc.Do("hincrby", key_name, field_name, incr_by_number)

	return result, err

}

//
///**
// * @func: HLEN  hash 中 一个key下的数量
// * @author Hank
// * @date   2019-06-21
// */
//
func (mng *RedisMng) HLen(key_name string) (res interface{}, err error) {

	rc := pool.Get()
	defer rc.Close()

	result, err := redis.String(rc.Do("hlen", key_name))
	return result, err

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
//func GETSET(key_name string, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("set", key_name, value))
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
//func SETEX(key_name string, expire interface{}, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("setex", key_name, expire, value)
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
//func SETNX(key_name string, expire interface{}, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("setnx", key_name, value)
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
//func INCR(key_name string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("incr", key_name)
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
//func DECR(key_name string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("decr", key_name)
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
//func INCRBY(key_name string, increment interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("incrby", key_name, increment)
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
//func DECRBY(key_name string, increment interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("decrby", key_name, increment)
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
//func EXPIRE(key_name string, expire interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("expire", key_name, expire)
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
//func LINDEX(key_name string, position interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("lindex", key_name, position))
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
//func LLEN(key_name string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int64(client.Do("llen", key_name))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////从列表中左边删除第一个元素 ，返回这个元素的值
//func LPOP(key_name string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("lpop", key_name))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////从列表中左边删除第一个元素 ，返回这个元素的值
//func RPOP(key_name string) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("rpop", key_name))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////添加元素 如果失败 0  成功  返回列表的长度
//func RPUSH(key_name string, value interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("rpush", key_name, value))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result, err
//}
//
////列表的长度 0 -1 全部
//func LRANGE(key_name string, start, end interface{}) (res interface{}, err error) {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("lrange", key_name, start, end))
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
//func ZADD(key_name string, score interface{}, member interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := client.Do("zadd", key_name, score, member)
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
//func ZCARD(key_name string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zcard", key_name))
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
//func ZCOUNT(key_name string, min interface{}, end interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zcount", key_name, min, end))
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
//func ZINCRBY(key_name string, increment interface{}, field string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zincrby", key_name, increment, field))
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
//func ZRANGE(key_name string, start_index, end_index interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("zrevrange", key_name, start_index, end_index))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////ZREM 删除一个有序集合中的指定成员  被成功移除的成员的数量，不包括被忽略的成员。
//func ZREM(key_name string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("zrem", key_name))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////Zscore 命令返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。 成员的分数值，以字符串形式表示。
//func ZSCORE(key_name, member string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("zscore", key_name, member))
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
//func SADD(key_name string, value interface{}) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.String(client.Do("zscore", key_name, value))
//
//	defer pool.Close()
//
//	defer client.Close()
//
//	return result
//}
//
////Scard 命令返回集合中元素的数量。
//func SCARD(key_name string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("scard", key_name))
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
//func SMEMBERS(key_name string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Strings(client.Do("smembers", key_name))
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
//func SREM(key_name string) interface{} {
//
//	pool := newRedisClient()
//
//	client := pool.Get()
//
//	result, _ := redis.Int(client.Do("srem", key_name))
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
//func SRANDMEMBER(key_name string, count interface{}) interface{} {
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
//		result, _ := redis.String(client.Do("srandmember", key_name))
//
//		return result
//
//	} else {
//
//		result, _ := redis.Strings(client.Do("srandmember", key_name, count))
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
