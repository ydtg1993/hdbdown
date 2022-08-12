package mongo

import (
	"context"
	"fmt"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// MongoDB 连接池
var DBClient *mongo.Database

// pool 连接池模式
func init() {

	host := os.Getenv("mongo_host")
	port := os.Getenv("mongo_port")
	user := os.Getenv("mongo_user")
	password := os.Getenv("mongo_password")
	dbName := os.Getenv("mongo_database")
	timeOut := 2
	maxNum := os.Getenv("mongo_database")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", user, password, host, port, dbName)
	//fmt.Println(uri)
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut))
	defer cancel()

	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)

	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	maxNums, _ := strconv.ParseUint(maxNum, 10, 64)
	o.SetMaxPoolSize(maxNums)

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
	DBClient = client.Database(dbName)

}
