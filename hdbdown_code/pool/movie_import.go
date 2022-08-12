package pool

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"hdbdown/models"
	"hdbdown/mongo"
	"hdbdown/rd"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RelationshipManager struct {
	ActorID     []int //演员数据
	LabelID     []int //标记
	SeriesID    []int //系列
	CompaniesID []int //片商
	CategoryID  []int //分类
	DirectorID  []int //导演
}

func (reShip *RelationshipManager) Update(mid int, cid int) (err error) {

	if mid == 0 || cid == 0{
		return errors.New("mid or cid不可以为空或者0")
	}

	// 建立演员和影片之间的关系
	for _, aid := range reShip.ActorID {
		maa := new(models.MovieActorAssociate)
		maa.Aid = aid
		maa.Mid = mid
		maa.Status = 1
		maa.AssociateTime = time.Now().Format("2006-01-02 15:04:05")
		rest := models.GetGormDb().Where("mid = ? and aid = ?", mid, aid).FirstOrCreate(&maa)
		if rest.Error != nil {
			return
		}

		// 更新演员和分类的关系
		var maca = new(models.MovieActorCategoryAssociate)
		maca.Cid = cid
		maca.Aid = aid
		maca.Status = 1
		maca.AssociateTime = time.Now().Format("2006-01-02 15:04:05")
		if err = models.GetGormDb().Where("aid = ? and cid = ?", aid, cid).FirstOrCreate(&maca).Error; err != nil {
			return
		}
	}


	// 建立与片商之间的关系
	for _, id := range reShip.CompaniesID {
		//fmt.Println(id)
		mcp := new(models.MovieFilmCompaniesAssociate)
		mcp.Mid = mid
		mcp.FilmCompaniesId = id
		mcp.Status = 1
		mcp.AssociateTime = time.Now().Format("2006-01-02 15:04:05")
		if err = models.GetGormDb().Where("mid = ? and film_companies_id = ?", mid, id).FirstOrCreate(&mcp).Error; err != nil {
			return
		}

		//  更新电影公司和分类的关系
		var mfca = new(models.MovieFilmCompaniesCategoryAssociate)
		mfca.Cid = cid
		mfca.FilmCompaniesId = id
		if err = models.GetGormDb().Where("cid = ? and film_companies_id = ?", cid, id).FirstOrCreate(&mfca).Error; err != nil {
			return
		}
	}


	// 建立与系类之间的关系
	for _, id := range reShip.SeriesID {
		//fmt.Println(id)
		ms := new(models.MovieSeriesAssociate)
		ms.AssociateTime = time.Now().Format("2006-01-02 15:04:05")
		ms.Status = 1
		ms.Mid = mid
		ms.SeriesId = id
		if err = models.GetGormDb().Where("mid = ? and series_id = ?", mid, id).FirstOrCreate(&ms).Error; err != nil {
			return
		}

		//  更新系列和分类的关系
		var msca = new(models.MovieSeriesCategoryAssociate)
		msca.Cid = cid
		msca.SeriesId = id
		if err = models.GetGormDb().Where("cid = ? and series_id = ?", cid, id).FirstOrCreate(&msca).Error; err != nil {
			return
		}
	}

	// 建立与标签之间的关系
	for _, id := range reShip.LabelID {
		//fmt.Println(id)
		ml := new(models.MovieLabelAssociate)
		ml.Mid = mid
		ml.Cid = id
		ml.Status = 1
		ml.AssociateTime = time.Now().Format("2006-01-02 15:04:05")

		if err = models.GetGormDb().Where("mid = ? and cid = ?", mid, id).FirstOrCreate(&ml).Error; err != nil {
			return
		}

		//  更新标签和分类的关系
		var mlca = new(models.MovieLabelCategoryAssociate)
		mlca.Cid = cid
		mlca.Lid = id
		if err = models.GetGormDb().Where("cid = ? and lid = ?", cid, id).FirstOrCreate(&mlca).Error; err != nil {
			return
		}

	}


	// 建立导演和影片之间的关系
	for _, id := range reShip.DirectorID {
		ml := new(models.MovieDirectorAssociate)
		ml.Mid = mid
		ml.Did = id
		ml.Status = 1
		ml.AssociateTime = time.Now().Format("2006-01-02 15:04:05")

		if err = models.GetGormDb().Where("mid = ? and did = ?", mid, id).FirstOrCreate(&ml).Error; err != nil {
			return
		}
	}

	return
}

func (reShip *RelationshipManager) pressMongoMovie(mongoData *mongo.MovieBase) (err error, downloadError []models.DownloadError) {

	// 更新演员信息
	err, reShip.ActorID = getActorId(mongoData)
	if err != nil {
		return
	}

	// 更新导演信息
	err, reShip.DirectorID = getDirectorId(mongoData)
	if err != nil {
		return
	}

	// 更新标签信息
	err , reShip.LabelID = getLabelId(mongoData)
	if err != nil {
		return
	}

	// 更新系列
	err, reShip.SeriesID = getSeriesId(mongoData)
	if err != nil {
		return
	}

	// 更新片商
	err, reShip.CompaniesID = getFilmCompaniesId(mongoData)
	if err != nil {
		return
	}

	if len(reShip.ActorID) == 0 {
		downloadError = append(downloadError, models.DownloadError{
			Title:   "演员信息",
			Message: "缺少数据",
		})
	}

	if len(reShip.LabelID) == 0 {
		downloadError = append(downloadError, models.DownloadError{
			Title:   "标签信息",
			Message: "缺少数据",
		})
	}

	if len(reShip.SeriesID) == 0 {
		downloadError = append(downloadError, models.DownloadError{
			Title:   "系列信息",
			Message: "缺少数据",
		})
	}

	if len(reShip.CompaniesID) == 0 {
		downloadError = append(downloadError, models.DownloadError{
			Title:   "片商信息",
			Message: "缺少数据",
		})
	}

	if len(reShip.DirectorID) == 0 {
		downloadError = append(downloadError, models.DownloadError{
			Title:   "导演信息",
			Message: "缺少数据",
		})
	}

	return
}

/**
查询 temporary_movie 的数据 ，将新番写入 movie
// 如果新增了演员，将演员数据添加到 redis 队列中 ，等待下载处理
// 如果新增了电影，将电影数据添加到 redis 队列中 ，等待下载处理
*/
func ProcessMovieData() {
	logs.Debug("ProcessMovieData 进程开启...")
	var pageNum = 500
	var lastId = 0
	var tmModel = new(models.TemporaryMovie)

	err, totalNum := tmModel.ListOfNotProcessedCount()
	if err != nil {
		logs.Error("temporary_movie 查询统计查询错误:", err.Error())
		return
	}

	pageAll := math.Ceil(float64(totalNum) / float64(pageNum))
	logs.Debug("ProcessMovieData 进程,", "总页数:", pageAll)
	for i := 0; i < int(pageAll); i++ {
		if rd.CheckLock() {
			return
		}
		logs.Debug("ProcessMovieData 进程,", "总页数:", pageAll, "当前页:", i)

		var res []*models.TemporaryMovie

		err, lastId, res = tmModel.ListOfNotProcessed(lastId, pageNum)
		if err != nil {
			logs.Error("temporary_movie 查询错误:", err.Error())
			continue
		}

		for _, v := range res {
			err = createMovieWithRelationship(v)
			if err != nil {
				continue
			}
		}

	}

}

func createMovieWithRelationship(temporaryData *models.TemporaryMovie) (err error) {
	var reShip RelationshipManager
	var downloadError []models.DownloadError

	//  检查影片是否被录入到 movie 表
	movie := new(models.Movie)
	err = movie.FindByNumber(temporaryData.Number)
	if err != nil {
		logs.Error(temporaryData.Number, "movie 影片查询错误:", err.Error())
		return
	}

	// 影片如果存在则跳过
	if movie.Id > 0 {
		if err = temporaryData.PressSuccess(); err != nil {
			logs.Error(temporaryData.Id, "temporary_movie 状态修改错误:", temporaryData.Id, err.Error())
		}
		return
	}

	// 从 mongo 读取电影数据
	err, mongoData := mongo.GetMovieDataByNumber(temporaryData.DbName, temporaryData.Number)
	if err != nil {
		logs.Error(temporaryData.Number, "movie 影片mongo查询错误:", err.Error())
		return
	}

	// 获得影片的关联信息，如果关联信息不存在则创建，最后获得关联信息相关的主键 ID
	err, downloadError = reShip.pressMongoMovie(mongoData)
	if err != nil {
		return
	}

	// 更新电影信息
	Mid , err := createMovie(mongoData, reShip)
	if err != nil {
		logs.Error(temporaryData.Number, "影片信息创建错误:", mongoData.Uid, err.Error())
		return
	}

	if len(mongoData.Magnet) == 0 {
		downloadError = append(downloadError, models.DownloadError{
			Title:   "磁链信息",
			Message: "无磁链",
		})
	}


	if len(downloadError) > 0 {
		var movieError = new(models.MovieError)
		movieError.Mid = Mid
		msg, err := json.Marshal(downloadError)
		if err != nil {
			logs.Error(mongoData.Uid, "json 压缩错误:", err.Error())
		}
		movieError.Message = string(msg)
		if err := movieError.Create(); err != nil {
			logs.Error(mongoData.Uid, "错误日志创建错误:", err.Error())
		}
	}

	// 完成导入
	err = temporaryData.PressSuccess()
	if err != nil {
		logs.Error(mongoData.Uid, "temporary_movie 状态修改错误:", temporaryData.Id, err.Error())
	}

	return
}

func getDirectorId(mongoData *mongo.MovieBase) (err error, ids []int) {

	if mongoData.Direct == "" {
		return nil , nil
	}

	name := strings.TrimSpace(mongoData.Direct)
	directID , err := rd.GetCashWithHash(models.DirectorNameWithIDHash, name, func() (id string, err error) {
		mc := new(models.MovieDirector)
		mc.Name = mongoData.Direct
		mc.MovieSum = 0
		mc.LikeSum = 0
		mc.Status = 1
		mc.Oid = 0

		if err := models.GetGormDb().Where("name = ?", mongoData.Direct).FirstOrCreate(&mc).Error ; err != nil{
			return "", err
		}

		return strconv.Itoa(mc.Id), nil
	})

	id, err := strconv.Atoi(directID)
	if err != nil {
		return
	}

	ids = append(ids, id)

	return
}

// 更新演员
func getActorId(mongoData *mongo.MovieBase) (err error, ids []int) {

	//logs.Debug(mv.Uid, "更新演员信息:", mv.Actor)
	actor := mongoData.ActorTypeChange()

	for _, v := range actor {

		if v["name"] == "" || v["name"] == "♀" || v["name"] == "♂" {
			continue
		}
		v["name"] = strings.TrimSpace(v["name"])

		actorId, err :=  rd.GetCashWithHash(models.ActorNameWithIDHash, v["name"], func() (string ,error) {
			var movieActor = new(models.MovieActor)
			movieActor.Sex = v["sex"]
			movieActor.Name = v["name"]
			movieActor.Photo = ""
			movieActor.Status = 1
			movieActor.LikeSum = 1
			movieActor.MovieSum = 1
			movieActor.SocialAccounts = ""
			movieActor.Oid = 0
			if err := models.GetGormDb().Where("name = ?", v["name"]).FirstOrCreate(&movieActor).Error; err != nil{
				return "",  err
			}
			return strconv.Itoa(movieActor.Id), nil
		})

		if err != nil {
			continue
		}

		id,err := strconv.Atoi(actorId)
		if err != nil {
			continue
		}

		ids = append(ids,id)
	}
	return
}

// 更新标签
func getLabelId(mongoData *mongo.MovieBase) (err error, ids []int) {

	videoSort := mongoData.VideoSortTypeChange()

	for _, labels := range videoSort {
		if labels == "" {
			continue
		}
		labels = strings.TrimSpace(labels)

		labelId, err := rd.GetCashWithHash(models.LabelNameWithIDHash, labels, func() (id string, err error) {
			var movieLabels models.MovieLabel
			movieLabels.Name = labels
			movieLabels.Cid = 0
			movieLabels.Status = 1
			movieLabels.Oid = 0
			movieLabels.Sort = 0
			movieLabels.ItemNum = 0
			movieLabels.LikeSum = 0
			if err := models.GetGormDb().Where("name = ?", labels).FirstOrCreate(&movieLabels).Error; err != nil {
				return "", err
			}
			id = strconv.Itoa(movieLabels.Id)
			return
		})

		id, err := strconv.Atoi(labelId)
		if err != nil {
			continue
		}

		ids = append(ids, id)
	}
	return
}

// 更新片商
func getFilmCompaniesId(mongoData *mongo.MovieBase) (err error, ids []int) {
	if mongoData.Producer == "" {
		return nil, nil
	}

	name := strings.TrimSpace(mongoData.Producer)
	filmID, err := rd.GetCashWithHash(models.FilmNameWithIDHash, name, func() (id string, err error) {
		mongoData.Producer = name
		mc := new(models.MovieFilmCompanies)
		mc.Name = mongoData.Producer
		mc.MovieSum = 0
		mc.LikeSum = 0
		mc.Status = 1
		mc.Oid = 0

		if err := models.GetGormDb().Where("name = ?", mongoData.Producer).FirstOrCreate(&mc).Error; err != nil {
			return "", err
		}
		id = strconv.Itoa(mc.Id)
		return
	})

	id, err := strconv.Atoi(filmID)
	if err != nil {
		return
	}

	ids = append(ids, id)

	return
}

// 更新系列
func getSeriesId(mongoData *mongo.MovieBase) (err error, ids []int) {
	if mongoData.Series == "" {
		return nil,nil
	}
	name := strings.TrimSpace(mongoData.Series)
	seriesID , err := rd.GetCashWithHash(models.SeriesNameWithIDHash, name , func() (id string, err error) {
		mongoData.Series = name
		ms := new(models.MovieSeries)
		ms.Name = mongoData.Series
		ms.MovieSum = 0
		ms.LikeSum = 0
		ms.Status = 1
		ms.Oid = 0

		if err := models.GetGormDb().Where("name = ?", mongoData.Series).FirstOrCreate(&ms).Error; err != nil {
			return "", err
		}
		return strconv.Itoa(ms.Id), nil
	})
	id ,err := strconv.Atoi(seriesID)
	if err != nil {
	  	return
	}

	ids = append(ids, id)

	return
}

/**
更新影片信息
tmpData 关联ID 数据
 */
func createMovie(mongoData *mongo.MovieBase, relShip RelationshipManager) ( Mid int, err error,) {

	movieCategory := new(models.MovieCategory)
	err, movieCategory = movieCategory.FindByNameWithRedis(mongoData.Group)
	if err != nil {
		return
	}

	// 根据 v.Number 查 mongo 中获取对应数据 ， javdb + javbus 等数据源的结合
	//var list []*mongo.MovieBase
	var isDownload = 1
	if len(mongoData.Magnet) > 0 {
		isDownload = 2
	}

	score := mongoData.ScoreTypeChange()
	scoreMan := mongoData.ScoreManTypeChange()

	var pictureMap []map[string]string
	for k, v := range mongoData.PreviewImg {
		var picture = make(map[string]string)
		picture["img"] = v
		// PreviewBigImg 不一定百分百存在，如果存在 PreviewBigImg 和  PreviewImg 等长，并且一一对应
		if len(mongoData.PreviewBigImg) > 0 {
			picture["big_img"] = getValue(k,  mongoData.PreviewBigImg)
		} else {
			picture["big_img"] = v
		}

		pictureMap = append(pictureMap, picture)
	}
	mspPictureJson, err := json.Marshal(pictureMap)
	mspPicture := string(mspPictureJson)

	var movie = models.Movie{
		Number:                mongoData.Uid,
		NumberSource:          mongoData.Uid,
		Name:                  mongoData.VideoTitle,
		Time:                  0,
		ReleaseTime:           mongoData.ReleaseTime,
		Issued:                "",
		Sell:                  mongoData.Sell,
		SmallCover:            mongoData.SmallCover, //  下载处理
		BigCove:               mongoData.BigCover,   //  下载处理
		Trailer:               mongoData.Trailer,    //  下载处理
		Map:                   mspPicture,    //  下载处理
		Score:                 0,
		ScorePeople:           0,
		CommentNum:            0, //  评论统计
		CollectionScore:       score,
		CollectionScorePeople: scoreMan,
		CollectionCommentNum:  0, //  评论数
		WanSee:                0,
		Seen:                  0,
		FluxLinkageNum:        0,  //  下载处理
		FluxLinkage:           "", //  下载处理
		StatusAudit:           1,  // '状态 1.待审核 2.通过 3.不通过'
		Status:                1,
		IsDownload:            isDownload, //  状态 1.不可下载  2.可下载
		IsSubtitle:            1,          // 状态 1.不含字幕  2.含字幕
		IsHot:                 1,          //状态 1.普通  2.热门
		IsShortComment:        1,          //  状态 1.不含短评  2.含短评
		IsUp:                  2,          // 状态 1.上架  2.下架
		NewCommentTime:        "",         // 最新评论时间
		FluxLinkageTime:       "",
		Oid:                   0,
		Cid:                   movieCategory.Id,
		Weight:                0,
	}

	//movie.Time = mv.VideoTime
	var valid = regexp.MustCompile("[0-9]")
	var resultValid [][]string
	resultValid = valid.FindAllStringSubmatch(mongoData.VideoTime, -1)
	var strTime string
	for _, v := range resultValid {
		strTime = strTime + v[0]
	}
	if strTime != "" {
		movie.Time, err = strconv.Atoi(strTime)
		if err != nil {
			return
		}
		movie.Time = movie.Time * 60
	} else {
		movie.Time = 3600
	}


	movie.FluxLinkageNum = len(mongoData.Magnet)
	movie.FluxLinkageTime = time.Now().Format("2006-01-02 15:04:05")

	fluxLinkage, err := json.Marshal(mongoData.Magnet)
	if err != nil {
		return
	}

	movie.FluxLinkage = string(fluxLinkage)

	// 写入 movie 表中


	if err = movie.Create(); err != nil {
		return
	}
	Mid = movie.Id

	// 建立与分类之间的关系
	mca := new(models.MovieCategoryAssociate)
	mca.Cid = movieCategory.Id
	mca.Mid = Mid
	mca.Status = 1
	mca.AssociateTime = time.Now().Format("2006-01-02 15:04:05")
	if err = models.GetGormDb().Where("mid = ? and cid = ?", Mid, movieCategory.Id).FirstOrCreate(&mca).Error; err != nil {
		return
	}

	if err = relShip.Update(Mid, movieCategory.Id) ; err != nil {
		return
	}

	return
}

