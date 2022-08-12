package pool

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
	"hdbdown/common"
	"hdbdown/models"
	"hdbdown/rd"
	"hdbdown/runtimes"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type MoviePictureQueue struct {
	Id    string
	Mv    models.Movie
	Lists []string
}

func (m *MoviePictureQueue) AppendList(uri string) {
	if uri != "" {
		_, uri1 := formatUrl(uri)
		URL, _ := url.Parse(uri1)
		url2 := strings.ReplaceAll(uri1, URL.RawQuery, "")
		url3 := strings.ReplaceAll(url2, "?", "")
		m.Lists = append(m.Lists, url3)
	}
}

func checkDisk() bool {
	//APP_ENV=dev
	AppEnv := os.Getenv("APP_ENV")

	if AppEnv == "dev" {
		wd, _ := os.Getwd()
		diskStatus := common.DiskUsage(wd)
		use := float64(diskStatus.Free) / float64(diskStatus.All)
		if use < 0.2 {
			logs.Error("UAT 磁盘空间不足，暂停图片下载功能...")
			return false
		}
	}
	return true
}

var threadPool = make(chan interface{}, 50)

func PressDownloadMoviePicture() {
	if checkDisk() == false {
		return
	}
	total, err := rd.LLen(models.MoviePicturePress)
	if err != nil {
		logs.Error("PressDownloadMoviePicture 进程","redis 队列读取错误:", err.Error())
		return
	}


	for i := 0; i < int(total); i++ {
		if rd.CheckLock() {
			return
		}
		// 根据队列读取
		number, err := rd.LPop(models.MoviePicturePress)
		if err != nil {
			logs.Error("PressDownloadMoviePicture 进程","redis 队列读取错误:", err.Error())
			continue
		}
		var movieData models.Movie
		rest := models.GetGormDb().Where("number = ?", number).First(&movieData)
		if rest.Error != nil && rest.Error != gorm.ErrRecordNotFound {
			logs.Error("PressDownloadMoviePicture 进程","movie 查询错误:", err, number)
			continue
		}

		var mvBase models.TemporaryMovie
		rests := models.GetGormDb().Where("number = ?", number).First(&mvBase)
		if rests.Error != nil && rests.Error != gorm.ErrRecordNotFound {
			logs.Error("PressDownloadMoviePicture 进程","temporary_movie 数据查询错误:", err, number)
			continue
		}

		if rests == nil {
			continue
		}

		mvQueueForCover := new(MoviePictureQueue)
		mvQueueForCover.Id = strconv.Itoa(movieData.Id)
		mvQueueForCover.Mv = movieData

		mvQueueForCover.AppendList(movieData.SmallCover)
		mvQueueForCover.AppendList(movieData.BigCove)


		threadPool <- mvQueueForCover
		go func(vv *MoviePictureQueue) {
			if len(mvQueueForCover.Lists) > 0 {
				//开启子线程下载
				oDown := new(runtimes.DoRun)
				//下载结果
				done, msg := oDown.Work(fmt.Sprintf("%s:%s:%s", vv.Mv.Number, "cove", mvQueueForCover.Id), mvQueueForCover.Lists)
				if done == false {
					var resultDownload []string
					resultDownload = append(resultDownload, msg)
					movieError := new(models.MovieError)
					if err := models.GetGormDb().Where("mid =?", vv.Mv.Id).First(&movieError).Error; err!=nil && err!=gorm.ErrRecordNotFound{
						logs.Error("PressDownloadMoviePicture 进程","电影错误信息查找失败:", err.Error())
					}

					msg , err := json.Marshal(resultDownload)
					if err != nil {
						logs.Error("PressDownloadMoviePicture 进程","压缩失败:", err.Error())
					}

					err = movieError.Add(vv.Mv.Id, "图片下载错误", string(msg))
					if err != nil {
						logs.Error("PressDownloadMoviePicture 进程","错误标记失败:", err.Error())
					}

					logs.Warning("影片下载记录失败,down error:", vv.Id, msg, vv.Mv.Number)

				}else {
					// 下载成功封面就标记上架自动过审等
					if err := vv.Mv.AutoSuccess(vv.Mv.Id); err != nil {
						logs.Error("PressDownloadMoviePicture 进程", "状态修改失败:", err.Error())
					}
				}
			}
			<-threadPool
		}(mvQueueForCover)



		MapDown := os.Getenv("MAP_DOWN")
		if MapDown != "1" {
			return
		}

		mvQueueForOther := new(MoviePictureQueue)
		mvQueueForOther.Id = strconv.Itoa(movieData.Id)
		mvQueueForOther.Mv = movieData
		mvQueueForOther.Lists = nil
		mvQueueForOther.AppendList(movieData.Trailer)
		var pictureMap []map[string]string
		if movieData.Map != "" {
			if err := json.Unmarshal([]byte(movieData.Map), &pictureMap); err != nil {
				logs.Error("PressDownloadMoviePicture 进程","图片数据解析错误:", err.Error())
				continue
			}
			for _, v := range pictureMap {
				if val, ok := v["img"] ; ok == true{
					mvQueueForOther.AppendList(val)
				}

				if val, ok := v["big_img"] ; ok == true{
					mvQueueForOther.AppendList(val)
				}
			}
		}

		err = download(mvQueueForOther)
		if err != nil {
			continue
		}
	}

}

func download(mvQueue *MoviePictureQueue) error {
	threadPool <- mvQueue
	go func(vv *MoviePictureQueue) {
		if len(mvQueue.Lists) > 0 {
			//开启子线程下载
			oDown := new(runtimes.DoRun)
			//下载结果
			done, msg := oDown.Work(fmt.Sprintf("%s:%s:%s", vv.Mv.Number, "picture", mvQueue.Id), mvQueue.Lists)
			if done == false {
				var resultDownload []string
				resultDownload = append(resultDownload, msg)
				movieError := new(models.MovieError)
				if err := models.GetGormDb().Where("mid =?", vv.Mv.Id).First(&movieError).Error; err!=nil && err!=gorm.ErrRecordNotFound{
					logs.Error("PressDownloadMoviePicture 进程","电影错误信息查找失败:", err.Error())
				}

				msg , err := json.Marshal(resultDownload)
				if err != nil {
					logs.Error("PressDownloadMoviePicture 进程","压缩失败:", err.Error())
				}

				err = movieError.Add(vv.Mv.Id, "图片下载错误", string(msg))
				if err != nil {
					logs.Error("PressDownloadMoviePicture 进程","错误标记失败:", err.Error())
				}

				logs.Warning("影片下载记录失败,down error:", vv.Id, msg, vv.Mv.Number)

			}
		}
		<-threadPool
	}(mvQueue)

	return nil
}
