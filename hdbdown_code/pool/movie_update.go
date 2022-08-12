package pool

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"hdbdown/models"
	"hdbdown/mongo"
	"hdbdown/rd"
	"math"
	"time"
)

func UpdateMovie()  {
	logs.Debug("开始更新任务....")
	//  根据 temproty_movie 表中 status 状态判断是否需要更新
	var pageNum = 500
	var page = 0
	var tmModel models.TemporaryMovie
	err , totalNum := tmModel.ListOfNeedUpdateCount()
	if err != nil {
		logs.Error("UpdateMovie 进程","temporary_movie 统计查询错误:", err.Error())
		return
	}

	pageAll := math.Ceil(float64(totalNum) / float64(pageNum))

	for page = 0; page < int(pageAll) ; page ++   {
		if rd.CheckLock() {
			return
		}
		var lists []*models.TemporaryMovie
		var updateData models.Movie
		var lastId int
		err, lastId ,lists = tmModel.ListOfNeedUpdate(lastId, pageNum)
		if err != nil {
			logs.Error("UpdateMovie 进程","temporary_movie 数据查询错误:", err.Error())
			continue
		}

		for _, v := range lists {
			err, mongoData := mongo.GetMovieDataByNumber(v.DbName, v.Number)
			if err != nil {
				logs.Error("UpdateMovie 进程","mongo 查询错误:", err.Error())
				continue
			}

			var movie models.Movie
			err = movie.FindByNumber(mongoData.Uid)
			if err != nil {
				logs.Error("UpdateMovie 进程","movie 查询错误:", v.Number, err.Error())
				continue
			}

			if movie.Id == 0 {
				if err := models.GetGormDb().Model(&v).Updates(models.TemporaryMovie{
					IsUpdate:  models.NoNeedUpdate,
				}).Error; err != nil{
					logs.Error("UpdateMovie 进程","temporary_movie 数据存储错误:", err.Error())
				}
				continue
			}

			updateData.ReleaseTime = mongoData.ReleaseTime
			updateData.Name = mongoData.VideoTitle
			updateData.Sell = mongoData.Sell

			/**
			磁链更新
			 */
			var magnets [] mongo.MagnetMode
			if movie.FluxLinkage != "" {
				err = json.Unmarshal([]byte(movie.FluxLinkage), &magnets)
				if err != nil {
					logs.Error("UpdateMovie 进程","movie磁链解析错误:",v.Number, movie.FluxLinkage, err.Error())
					continue
				}
			}

			// 需要更新的数据
			cloneMagents := magnets
			for _, val := range mongoData.Magnet {
				var update = false
				for _, value := range magnets {
					if val.Url == value.Url {
						update = true
					}
				}

				if update == false {
					// 不存在
					cloneMagents = append(cloneMagents, val)
				}
			}

			if len(cloneMagents) > len(magnets) {
				// 有更新
				jsonStr, err := json.Marshal(cloneMagents)
				if err != nil {
					logs.Error("UpdateMovie 进程","movie 磁链数据 json 压缩错误:", v.Number, err.Error())
					continue
				}
				updateData.FluxLinkageNum =  len(cloneMagents)
				updateData.FluxLinkage =  string(jsonStr)
				updateData.FluxLinkageTime =  time.Now().Format("2006-01-02 15:04:05")
			}

			/*
				判断图片信息是否有变更，如果有变更加入到图片下载队列
			*/
			var pictureMap []map[string]string
			for k, v := range mongoData.PreviewImg {
				var picture = make(map[string]string)
				picture["img"] = v
				if len(mongoData.PreviewBigImg) > 0 {
					picture["big_img"] = getValue(k,  mongoData.PreviewBigImg)
				} else {
					picture["big_img"] = v
				}

				pictureMap = append(pictureMap, picture)
			}
			mspPictureJson, err := json.Marshal(pictureMap)
			mspPicture := string(mspPictureJson)
			var isUpdatePicture = false
			if movie.SmallCover != mongoData.SmallCover  {
				updateData.SmallCover =  mongoData.SmallCover
				isUpdatePicture = true
			}

			if  movie.BigCove != mongoData.BigCover  {
				updateData.BigCove =  mongoData.BigCover
				isUpdatePicture = true
			}

			if  movie.Trailer != mongoData.Trailer  {
				updateData.Trailer =  mongoData.Trailer
				isUpdatePicture = true
			}

			if movie.Map != mspPicture {
				updateData.Map = mspPicture
				isUpdatePicture = true
			}

			if isUpdatePicture == true {
				logs.Debug(movie.Number, "图片需要更新，加入图片下载队列...")
				err = rd.RPush(models.MoviePicturePress, v.Number)
				if err != nil {
					return
				}
			}

			var reShip RelationshipManager

			if err, _ = reShip.pressMongoMovie(mongoData) ; err != nil {
				logs.Error("UpdateMovie 进程","movie 关联数据处理失败:", v.Number, err.Error())
				continue
			}

			if err = reShip.Update(movie.Id, movie.Cid) ; err != nil{
				logs.Error("UpdateMovie 进程","movie 关联数据处理失败:", v.Number, err.Error())
				continue
			}


			if err := models.GetGormDb().Model(&movie).Where("id = ?", movie.Id).Updates(updateData).Error; err != nil {
				logs.Error("UpdateMovie 进程","movie 数据更新失败:", v.Number, err.Error())
				continue
			}

			if err := models.GetGormDb().Model(&v).Updates(models.TemporaryMovie{
				IsUpdate:  models.NoNeedUpdate,
			}).Error; err != nil{
				logs.Error("UpdateMovie 进程","temporary_movie 状态更新错误错误:", v.Number, v.Id, err.Error())
			}
		}


	}

}