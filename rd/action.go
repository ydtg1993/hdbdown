package rd

import (
	"github.com/gomodule/redigo/redis"
)

/**
* 获取list的数量
* param		key 	主键
 */
func ListCount(key string) int {
	RD := RdClient.Get()
	defer RD.Close()

	res, _ := redis.Int(RD.Do("LLEN", key))
	return res
}

/**
* 从list最后移除并读取一条数据
 */
func ListGetOneByLast(key string) string {
	RD := RdClient.Get()
	defer RD.Close()

	res, _ := redis.String(RD.Do("RPOP", key))
	return res
}

/**
* 写入字符串
* param 	key 	主键
* param		val 	字符串
 */
func StringSet(key, val string) {
	RD := RdClient.Get()
	defer RD.Close()

	//写入redis
	RD.Send("set", key, val)

}

/**
* 读取字符串
* param 	key 	主键
 */
func StringGet(key string) string {
	RD := RdClient.Get()
	defer RD.Close()

	res, _ := redis.String(RD.Do("get", key))
	return res
}

/**
* 删除一个key
 */
func DelKey(key string) {
	RD := RdClient.Get()
	defer RD.Close()

	RD.Do("del", key)
}
