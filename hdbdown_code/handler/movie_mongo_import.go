package handler

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/golang-module/carbon"
	"hdbdown/global/orm"
	"hdbdown/models"
	"hdbdown/tools/config"
	"hdbdown/tools/mongo"
	"hdbdown/tools/rd"
	"math"
	"time"
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

//  从 mongo 读取数据 ，写入 movie
func importJavData(dbName string) {

	//进程开始时间
	startTime := time.Now()
	//写入计数
	total := 0     // 总处理
	totalSkip := 0 // 跳过
	totalNew := 0  // 新增

	// 分页处理，每页条数
	var pageNum float64 = 500

	tmpMovie := new(models.TemporaryMovie)
	if err := orm.Eloquent.Where("db_name = ?", dbName).Order("ctime desc").Limit(1).Find(&tmpMovie).Error; err != nil {
		logs.Error("ImportMongoDb 进程", "temp..._movie 数据查询错误数据", err.Error())
		return
	}

	lastTime := tmpMovie.Ctime
	//lastTime = "2022-08-22 00:00:01"
	if lastTime == "" {
		// 项目开始时间，爬取的最早时间点
		lastTime = carbon.CreateFromDateTime(2021, 1, 1, 0, 0, 0).ToDateTimeString()
	} else {
		t := carbon.Parse(lastTime).Timestamp()
		// !!! 向前推 7 天更新数据  !!!
		//t = t - 604800 //(3600 * 24 * 7)
		t = t - 3600*24*config.Spe.DaysInAdvance //(3600 * 24 * 7)
		lastTime = carbon.CreateFromTimestamp(t).ToDateTimeString()
	}

	err, maxNum := mongo.Count(dbName, lastTime)
	logs.Debug(dbName, "数据导入开始, 总数据:", maxNum, "条")
	if err != nil {
		logs.Error("ImportMongoDb 进程", "mongo 数据统计查询错误数据", err.Error())
		return
	}

	pageAll := math.Ceil(float64(maxNum) / pageNum)

	for i := 0; i < int(pageAll); i++ {
		logs.Debug("ImportMongoDb 进程, 数据库", dbName, fmt.Sprintf(" 第%d页, 总页数:%d", i, int(pageAll)))
		if rd.CheckLock() {
			return
		}
		err, movieBaseList := mongo.Find(dbName, lastTime, int64(i+1), int64(pageNum))
		if err != nil {
			logs.Error("ImportMongoDb 进程", "mongo 数据查询错误数据", err.Error())
			return
		}

		for _, v := range movieBaseList {

			if v.Uid == "" {
				obj, _ := json.Marshal(v)
				logs.Error("ImportMongoDb 进程", "mongo 数据异常", string(obj))
				continue
			}

			tmpModel := new(models.TemporaryMovie)
			if err := orm.Eloquent.Where("number = ? and db_name = ?", v.Uid, dbName).Select("id", "ctime", "utime").Find(&tmpModel).Error; err != nil {
				logs.Error("ImportMongoDb 进程", "mongo 数据查询错误数据", err.Error())
				continue
			}

			if tmpModel.Id > 0 {
				vuTime := carbon.Parse(v.Utime).Timestamp()
				tuTime := carbon.Parse(tmpModel.Utime).Timestamp()

				if vuTime > tuTime {
					data := models.TemporaryMovie{
						IsUpdate: models.NeedUpdate, // 需要更新
						Utime:    v.Utime,
					}

					if v.BigCover == "" && v.SmallCover == "" {
						data.Status = models.StatusNotUnusual // 数据不完整，不做处理
					}

					if err := orm.Eloquent.Model(&tmpModel).Updates(data).Error; err != nil {
						logs.Error("ImportMongoDb 进程", "temp__movie 更新失败", v.Uid, err.Error())
						continue
					}
				}

				totalSkip++
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
					logs.Error("ImportMongoDb 进程", "temp__movie 创建失败", v.Uid, err)
					continue
				}
				totalNew++
			}
			total++
		}
	}
	logs.Debug(dbName, "总数据:", maxNum, " 成功处理:", total, "条", "新增:", totalNew, "跳过:", totalSkip, "总耗时：", time.Now().Unix()-startTime.Unix())
}
