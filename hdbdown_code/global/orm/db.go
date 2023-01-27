package orm

import (
	goredis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// mysql 连接池
var Eloquent *gorm.DB

// MongoDB 连接池
var DBClient *mongo.Database

// redis 连接池
var Client *goredis.Client
