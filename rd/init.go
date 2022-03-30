package rd

import (
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

var RdClient *redis.Pool
var RdDB string

func init() {
	redishost, _ := beego.AppConfig.String("redishost")
	redisport, _ := beego.AppConfig.String("redisport")
	redispass, _ := beego.AppConfig.String("redispass")
	RdDB, _ = beego.AppConfig.String("redisdb")

	//使用redis连接池
	RdClient = PoolInitRedis(redishost+":"+redisport, redispass)
	RdClient.Stats()
}

// redis连接池
func PoolInitRedis(server string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 60 * time.Second,
		MaxActive:   10000, //最大数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if RdDB != "" {
				if _, err := c.Do("SELECT", RdDB); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
