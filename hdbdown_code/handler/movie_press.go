package handler

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models"
	"hdbdown/tools"
	"hdbdown/tools/config"
	"hdbdown/tools/download"
	"hdbdown/tools/rd"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
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
	if config.Spe.AppEnv == "dev" {
		wd, _ := os.Getwd()
		diskStatus := tools.DiskUsage(wd)
		use := float64(diskStatus.Free) / float64(diskStatus.All)
		if use < 0.2 {
			logs.Error("UAT 磁盘空间不足，暂停图片下载功能...")
			return false
		}
	}
	return true
}

var threadPool = make(chan bool, 10)

func PressDownloadMoviePicture() {
	if checkDisk() == false {
		return
	}

	total, err := rd.LLen(models.MoviePicturePress)
	if err != nil {
		logs.Error("PressDownloadMoviePicture 进程", "redis 队列读取错误:", err.Error())
		return
	}
	logs.Debug("开启图片资料下载队列处理....", total)

	if total < 1 {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(int(total))

	for i := 0; i < int(total); i++ {

		if rd.CheckLock() {
			return
		}
		// 根据队列读取
		number, err := rd.LPop(models.MoviePicturePress)
		if err != nil {
			logs.Error("PressDownloadMoviePicture 进程", "redis 队列读取错误:", err.Error())
			continue
		}
		var movieData models.Movie
		rest := orm.Eloquent.Where("number = ?", number).First(&movieData)
		if rest.Error != nil && rest.Error != gorm.ErrRecordNotFound {
			logs.Error("PressDownloadMoviePicture 进程", "movie 查询错误:", err, number)
			continue
		}

		var mvBase models.TemporaryMovie
		err = mvBase.GetDataByNumber(number)
		if err != nil {
			logs.Error("PressDownloadMoviePicture 进程", "temporary_movie 数据查询错误:", err, number)
			continue
		}

		if mvBase.Id == 0 {
			logs.Error("PressDownloadMoviePicture 进程", "异常数据:", number)
			continue
		}

		threadPool <- true
		go pressPicture(movieData, wg)

	}

	wg.Wait()

}

func pressPicture(movieData models.Movie, wg *sync.WaitGroup) {

	mvQueueForOther, err := makePictureQueue(movieData)
	if err != nil {
		logs.Error("PressDownloadMoviePicture 进程", err.Error())
		return
	}

	downloadPicture(mvQueueForOther, func(mv models.Movie) {
		if err := mv.AutoSuccess(mv.Id); err != nil {
			logs.Error("PressDownloadMoviePicture 进程", "状态修改失败:", err.Error())
		}
	}, wg)

}

func makePictureQueue(movieData models.Movie) (*MoviePictureQueue, error) {
	mvQueueForOther := new(MoviePictureQueue)
	mvQueueForOther.Id = strconv.Itoa(movieData.Id)
	mvQueueForOther.Mv = movieData
	mvQueueForOther.AppendList(movieData.SmallCover)
	mvQueueForOther.AppendList(movieData.BigCove)

	// UAT 不下载组图
	if config.Spe.MapDown != 1 {
		return mvQueueForOther, nil
	}

	mvQueueForOther.AppendList(movieData.Trailer)
	var pictureMap []string
	err, pictureMap := movieData.GetMapList()
	if err != nil {
		return nil, err
	}

	for _, v := range pictureMap {
		mvQueueForOther.AppendList(v)
	}

	return mvQueueForOther, nil
}

func downloadPicture(mvQueue *MoviePictureQueue, successFun func(mv models.Movie), wg *sync.WaitGroup) {
	defer wg.Done()

	if len(mvQueue.Lists) > 0 {
		//开启子线程下载
		//下载结果
		oDown := new(download.DoRun)
		done, msg := oDown.Work(fmt.Sprintf("%s:%s:%s", mvQueue.Mv.Number, "picture", mvQueue.Id), mvQueue.Lists)
		if done == false {
			movieError := new(models.MovieError)
			if err := orm.Eloquent.Where("mid =?", mvQueue.Mv.Id).First(&movieError).Error; err != nil && err != gorm.ErrRecordNotFound {
				logs.Error("PressDownloadMoviePicture 进程", "电影错误信息查找失败:", err.Error())
			}

			err := movieError.Add(mvQueue.Mv.Id, "图片下载错误", msg)
			if err != nil {
				logs.Error("PressDownloadMoviePicture 进程", "错误标记失败:", err.Error())
			}
			logs.Warning("影片下载记录失败,down error:", mvQueue.Id, msg, mvQueue.Mv.Number)
		} else {
			if successFun != nil {
				successFun(mvQueue.Mv)
			}
		}
	}

	<-threadPool
}
