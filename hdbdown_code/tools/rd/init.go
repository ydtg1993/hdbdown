package rd

import (
	"context"
	"fmt"
	goredis "github.com/go-redis/redis/v8"
	"hdbdown/global/orm"
	"hdbdown/tools/config"
	"time"
)

type RedisManage struct {
}

func (r RedisManage) SetUp() (err error) {
	orm.Client, err = open(fmt.Sprintf("%s:%d", config.Spe.RedisHost, config.Spe.RedisPort), config.Spe.RedisPass, config.Spe.RedisDb)
	return
}

// redis连接池
func open(server string, password string, db int) (*goredis.Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Network:  "tcp",
		Addr:     server,
		Password: password,
		DB:       db,
		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池数量
		MinIdleConns: 10, //好比最小连接数
		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	})

	_, err := rdb.Ping(context.Background()).Result()
	return rdb, err
}
