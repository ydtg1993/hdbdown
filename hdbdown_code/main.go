package main

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	_ "hdbdown/log"
	"hdbdown/mongo"
	"hdbdown/pool"
	"hdbdown/rd"
	_ "hdbdown/rd"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func init() {
	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for s := range c {
			switch s {
			case os.Kill: // kill -9 pid，下面的fmt无效
				fmt.Println("强制退出", s)
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				exitFunc()
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
	go pool.TaskImportMongo()
	//根据 redis 队列更新演员信息
	go pool.TaskUpdateActorData()

	// update movie 更新队列
	go pool.TaskUpdateMovieData()

	// 开启图片资料下载队列处理
	go pool.DownloadMoviePicture()

	pool.TaskImportDB()

}

func exitFunc(){
	rd.LockSystem()
	err := mongo.DBClient.Client().Disconnect(context.TODO())
	if err != nil {
		logs.Error("mongo 链接关闭错误:", err.Error())
	}
	fmt.Println("安全退出中,等待 60s...")
}
