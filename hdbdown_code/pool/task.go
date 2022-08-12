package pool

import (
	"time"
)

// 下载 mongo 数据 temporary_movie
func TaskImportMongo()  {
	t := time.NewTicker(time.Second * time.Duration(600))

	defer t.Stop()

	ImportMongoDb()

	for {
		<-t.C
		ImportMongoDb()
	}
}


// 根据 temporary_movie 的数据查询 mongo ，将数据写入到 movie * 细列表中
func TaskImportDB()  {
	t := time.NewTicker(time.Second * time.Duration(300))
	defer t.Stop()

	ProcessMovieData()

	for {
		<-t.C
		ProcessMovieData()
	}
}


/**
根据 redis 队列更新演员信息
 */
func TaskUpdateActorData()  {
	t := time.NewTicker(time.Second * time.Duration(600))

	defer t.Stop()

	ImportUpdateActorData()

	for {
		<-t.C
		ImportUpdateActorData()
	}
}

/**
资源更新
暂时只支持磁链更新
 */
func TaskUpdateMovieData()  {
	t := time.NewTicker(time.Second * time.Duration(600))

	defer t.Stop()

	UpdateMovie()

	for {
		<-t.C
		UpdateMovie()
	}
}



/**
下载影片图片资源
 */
func DownloadMoviePicture()  {
	t := time.NewTicker(time.Second * time.Duration(600))

	defer t.Stop()

	PressDownloadMoviePicture()

	for {
		<-t.C
		PressDownloadMoviePicture()
	}
}