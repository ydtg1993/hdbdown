package main

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"hdbdown/global/orm"
	"hdbdown/handler"
	"hdbdown/tools/config"
	"hdbdown/tools/database"
	"hdbdown/tools/download"
	"hdbdown/tools/log"
	"hdbdown/tools/mongo"
	"hdbdown/tools/rd"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func Setup() {
	err := config.Spe.SetUp()
	if err != nil {
		panic(err)
	}

	mylog := new(log.LogsManage)
	err = mylog.SetUp()
	if err != nil {
		panic(err)
	}

	db := new(database.MysqlManage)
	err = db.Setup()
	if err != nil {
		panic(err)
	}

	mongoDb := new(mongo.Mange)
	err = mongoDb.SetUp()
	if err != nil {
		panic(err)
	}

	redisManage := new(rd.RedisManage)
	err = redisManage.SetUp()
	if err != nil {
		panic(err)
	}

	// 开始前的线程数
	logs.Debug("线程数量 starting: %d\n", runtime.NumGoroutine())

}

func init() {
	Setup()
	download.Setup()

	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for s := range c {
			switch s {
			case os.Kill: // kill -9 pid，下面的fmt无效
				fmt.Println("强制退出", s)
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				func() {
					rd.LockSystem()
					err := orm.DBClient.Client().Disconnect(context.TODO())
					if err != nil {
						logs.Error("mongo 链接关闭错误:", err.Error())
					}
					fmt.Println("安全退出中,等待 60s...")
				}()
				time.Sleep(time.Second * 60)
				os.Exit(0)
			default:
			}
		}
	}()

}

func main() {
	rd.UnLockSystem()
	//下载 mongo 数据 temporary_movie
	go handler.TaskImportMongo()
	//根据 redis 队列更新演员信息
	go handler.TaskUpdateActorData()

	// update movie 更新队列
	go handler.TaskUpdateMovieData()

	// 开启图片资料下载队列处理
	go handler.DownloadMoviePicture()

	handler.TaskImportDB()

}
