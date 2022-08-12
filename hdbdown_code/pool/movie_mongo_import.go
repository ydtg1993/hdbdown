package pool

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/golang-module/carbon"
	"hdbdown/models"
	"hdbdown/mongo"
	"hdbdown/rd"
	"math"
	"os"
	"strconv"

	"time"

	_ "github.com/joho/godotenv/autoload"
)

/**
从 mongo 拉去数据 ，导入 temporary_movie
*/
func ImportMongoDb() {
	javlist := []string{"javlibrary", "javbus", "javdb"}

	for _, v := range javlist {
		importJavData(v)
	}
}

var DaysInAdvance int64 = 1

func init() {

	if days := os.Getenv("DaysInAdvance");days != ""{
		var err error
		DaysInAdvance , err = strconv.ParseInt(days, 64, 10)
		if err != nil {
			DaysInAdvance = 1
		}
	}

}

//  从 mongo 读取数据 ，写入 movie
func importJavData(dbName string){

	//进程开始时间
	startTime := time.Now()
	//写入计数
	total := 0
	// 分页处理，每页条数
	var pageNum float64 = 500

	tmpMovie := new(models.TemporaryMovie)
	if err := models.GetGormDb().Where("db_name = ?", dbName).Order("ctime desc").Limit(1).Find(&tmpMovie).Error ; err != nil{
		logs.Error("ImportMongoDb 进程", "temp..._movie 数据查询错误数据", err.Error())
		return
	}


	lastTime := tmpMovie.Ctime
	if lastTime == "" {
		// 项目开始时间，爬取的最早时间点
		lastTime = carbon.CreateFromDateTime(2021, 1, 1, 0, 0, 0).ToDateTimeString()
	}else{
		t := carbon.Parse(lastTime).Timestamp()
		// !!! 向前推 7 天更新数据  !!!
		//t = t - 604800 //(3600 * 24 * 7)
		t = t - 3600 * 24 * DaysInAdvance //(3600 * 24 * 7)
		lastTime = carbon.CreateFromTimestamp(t).ToDateTimeString()
	}

	err, maxNum := mongo.Count(dbName, lastTime)
	logs.Debug(dbName, "数据导入开始", maxNum, "条")
	if err != nil {
		logs.Error("ImportMongoDb 进程", "mongo 数据统计查询错误数据", err.Error())
		return
	}

	pageAll := math.Ceil(float64(maxNum) / pageNum)

	for i := 0; i < int(pageAll); i++ {
		logs.Debug("ImportMongoDb 进程, 数据库", dbName, fmt.Sprintf(" 第%d页, 总页数:%f", i, pageAll))
		if rd.CheckLock() {
			return
		}
		err, movieBaseList := mongo.Find(dbName, lastTime, int64(i + 1), int64(pageNum))
		if err != nil {
			logs.Error("ImportMongoDb 进程", "mongo 数据查询错误数据",err.Error())
			return
		}

		for _, v := range movieBaseList {
			tmpModel := new(models.TemporaryMovie)
			if err := models.GetGormDb().Where("number = ? and db_name = ?", v.Uid, dbName).Select("id", "ctime", "utime").Find(&tmpModel).Error ; err != nil{
				logs.Error("ImportMongoDb 进程", "mongo 数据查询错误数据", err.Error())
				continue
			}

			if tmpModel.Id > 0 {
				if carbon.Parse(v.Utime).Timestamp()  > carbon.Parse(tmpModel.Utime).Timestamp() {
					tmpModel.Utime = v.Utime              // 更新时间节点
					tmpModel.IsUpdate = models.NeedUpdate // 状态修改为可更新
					if err := models.GetGormDb().Model(&tmpModel).Updates(models.TemporaryMovie{
						IsUpdate:  models.NeedUpdate,
						Utime:     v.Utime,
					}).Error ; err != nil{
						logs.Error("ImportMongoDb 进程", "temp__movie 更新失败", v.Uid,err.Error())
						continue
					}
				}
			} else {
				tmpModel.Status = models.StatusNotProcessed
				// 异常数据的判读, 封面不存在
				if v.BigCover == "" && v.SmallCover == "" {
					tmpModel.Status = models.StatusNotUnusual // 数据不完整，不做处理
				}

				tmpModel.IsUpdate = models.NoNeedUpdate
				tmpModel.Ctime = v.Ctime
				tmpModel.Number = v.Uid
				tmpModel.DbName = dbName

				err := tmpModel.Create()
				if err != nil {
					logs.Error("ImportMongoDb 进程", "temp__movie 创建失败" , v.Uid, err)
				}
			}
			total++
		}
	}
	logs.Debug(dbName, "数据导入成功", total, "条,耗时：", time.Now().Unix()-startTime.Unix())
}
