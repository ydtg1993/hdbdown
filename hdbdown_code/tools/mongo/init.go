package mongo

import (
	"context"
	"fmt"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"hdbdown/global/orm"
	"hdbdown/tools/config"
)

type Mange struct {
}

// pool 连接池模式
func (m Mange) SetUp() (err error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", config.Spe.MongoUser, config.Spe.MongoPassword, config.Spe.MongoHost, config.Spe.MongoPort, config.Spe.MongoDatabase)
	fmt.Println(uri)
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), config.Spe.MongoTimeout)
	defer cancel()

	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)

	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	o.SetMaxPoolSize(config.Spe.MongoMaxNum)

	// 发起链接
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		fmt.Println("ConnectToDB", err)
		log.Error("ConnectToDB", err)
		return
	}

	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		fmt.Println("ConnectToDB", err)
		return
	}

	// 返回 client
	orm.DBClient = client.Database(config.Spe.MongoDatabase)

	return
}
