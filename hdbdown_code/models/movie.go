package models

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"hdbdown/rd"
	"strings"
	"time"
)

// 影片图片处理 redis 队列 key
const MoviePicturePress = "movie_picture_press"

/*
movie
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `number` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '番号',
  `number_source` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '原番号',
  `name` varchar(750) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '影片名称',
  `time` int unsigned DEFAULT NULL COMMENT '播放时长/秒',
  `release_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '发布/发行时间',
  `issued` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '发行',
  `sell` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '卖家',
  `small_cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '小封面',
  `big_cove` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '大封面',
  `trailer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '预告片',
  `map` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'json 数组 其他组图-预览图',
  `score` float unsigned DEFAULT NULL COMMENT '评分-计算过后的',
  `score_people` int unsigned NOT NULL DEFAULT '0' COMMENT '评分人数-冗余',
  `comment_num` int unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
  `collection_score` float unsigned NOT NULL COMMENT '评分-采集',
  `collection_score_people` int unsigned NOT NULL DEFAULT '0' COMMENT '评分人数-冗余-采集',
  `collection_comment_num` int unsigned NOT NULL DEFAULT '0' COMMENT '评论数-采集',
  `wan_see` int unsigned NOT NULL DEFAULT '0' COMMENT '想看数量-冗余',
  `seen` int unsigned NOT NULL DEFAULT '0' COMMENT '看过数量-冗余',
  `flux_linkage_num` int unsigned NOT NULL DEFAULT '0' COMMENT '磁链信息数',
  `flux_linkage` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'json 数组 磁链信息',
  `status_audit` tinyint DEFAULT '1' COMMENT '状态 1.待审核 2.通过 3.不通过',
  `status` tinyint DEFAULT '1' COMMENT '状态 1.正常  2.禁用  ',
  `is_download` tinyint DEFAULT '1' COMMENT '状态 1.不可下载  2.可下载  ',
  `is_subtitle` tinyint DEFAULT '1' COMMENT '状态 1.不含字幕  2.含字幕  ',
  `is_hot` tinyint DEFAULT '1' COMMENT '状态 1.普通  2.热门  ',
  `is_short_comment` tinyint DEFAULT '1' COMMENT '状态 1.不含短评  2.含短评  ',
  `is_up` tinyint DEFAULT '1' COMMENT '状态 1.上架  2.下架  ',
  `new_comment_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '最新评论时间 -冗余',
  `flux_linkage_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '磁链更新时间',
  `oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
  `cid` int DEFAULT '0' COMMENT '分类id',
  `weight` int DEFAULT '0' COMMENT '加权分',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
*/
type Movie struct {
	Id                    int     `json:"id" bson:"id" gorm:"primarykey"`
	Number                string  `json:"number" bson:"number"`
	NumberSource          string  `json:"number_source" bson:"number_source"`
	Name                  string  `json:"name" bson:"name"`
	Time                  int     `json:"time" bson:"time"`
	ReleaseTime           string  `json:"release_time" bson:"release_time"`
	Issued                string  `json:"issued" bson:"issued"`
	Sell                  string  `json:"sell" bson:"sell"`
	SmallCover            string  `json:"small_cover" bson:"small_cover"`
	BigCove               string  `json:"big_cove" bson:"big_cove"`
	Trailer               string  `json:"trailer" bson:"trailer"`
	Map                   string  `json:"map" bson:"map"`
	Score                 float32 `json:"score" bson:"score"`
	ScorePeople           int     `json:"score_people" bson:"score_people"`
	CommentNum            int     `json:"comment_num" bson:"comment_num"`
	CollectionScore       float32 `json:"collection_score" bson:"collection_score"`
	CollectionScorePeople float32 `json:"collection_score_people" bson:"collection_score_people"`
	CollectionCommentNum  int     `json:"collection_comment_num" bson:"collection_comment_num"`
	WanSee                int     `json:"wan_see" bson:"wan_see"`
	Seen                  int     `json:"seen" bson:"seen"`
	FluxLinkageNum        int     `json:"flux_linkage_num" bson:"flux_linkage_num"`
	FluxLinkage           string  `json:"flux_linkage" bson:"flux_linkage"`
	StatusAudit           int     `json:"status_audit" bson:"status_audit"`
	Status                int     `json:"status" bson:"status"`
	IsDownload            int     `json:"is_download" bson:"is_download"`
	IsSubtitle            int     `json:"is_subtitle" bson:"is_subtitle"`
	IsHot                 int     `json:"is_hot" bson:"is_hot"`
	IsShortComment        int     `json:"is_short_comment" bson:"is_short_comment"`
	IsUp                  int     `json:"is_up" bson:"is_up"`
	NewCommentTime        string  `json:"new_comment_time" bson:"new_comment_time"`
	FluxLinkageTime       string  `json:"flux_linkage_time" bson:"flux_linkage_time"`
	Oid                   int     `json:"oid" bson:"oid"`
	Cid                   int     `json:"cid" bson:"cid"`
	Weight                int     `json:"weight" bson:"weight"`
	CreatedAt             string  `json:"created_at" bson:"created_at"`
	UpdatedAt             string  `json:"updated_at" bson:"updated_at"`
	//Magnet []mongo.MagnetMode
}

/**
指定表名
*/
func (Movie) TableName() string {
	return "movie"
}

func (m *Movie) AutoSuccess(id int) error {
	if err := GetGormDb().Where("id =?",id).Updates(Movie{
		StatusAudit: 2,
		Status:      1,
		IsUp:        1,
	}).Error; err != nil {
		return err
	}
	return nil
}

/**
钩子函数
*/
func (m *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	m.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	if m.ReleaseTime != "" {

		ts, err := dateparse.ParseLocal(strings.TrimSpace(m.ReleaseTime))
		if err != nil {
			fmt.Println(m.Number, "at 时间错误:", m.ReleaseTime)
			m.ReleaseTime = ""
		} else {
			m.ReleaseTime = ts.Format("2006-01-02 15:04:05")
		}
	}

	if m.ReleaseTime == "" {
		m.ReleaseTime = time.Now().Format("2006-01-02 15:04:05")
	}
	if m.NewCommentTime == "" {
		m.NewCommentTime = time.Now().Format("2006-01-02 15:04:05")
	}
	if m.FluxLinkageTime == "" {
		m.FluxLinkageTime = time.Now().Format("2006-01-02 15:04:05")
	}


	return
}

func (m *Movie) AfterCreate(tx *gorm.DB) (err error) {
	err = rd.RPush(MoviePicturePress, m.Number)
	if err != nil {
		return
	}
	return
}

//Magnet []mongo.MagnetMode

func (m *Movie) FindByNumber(number string) (err error) {
	res := GetGormDb().Where("number = ?", number).Find(&m)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err = res.Error
		return
	}
	return
}

/**
* 从数据库中读取列表,每次读取10条
* mid			id
* updatedAt		最后更新时间
* limit 		每次读取多少条
 */
func (d *Movie) Lists(mid, limit int, restart bool, rTime string) (lastid int, res []*Movie) {

	q := `SELECT id,number,small_cover,big_cove,trailer,map
		FROM movie 
		where id>? and is_up=1 order by id asc limit %d;`

	if restart == true {
		q = `SELECT id,number,small_cover,big_cove,trailer,map
		FROM movie 
		where id>? and created_at>='` + rTime + `' order by id asc limit %d;`
	}

	q = fmt.Sprintf(q, limit)

	rows, err := DB.Query(q, mid)
	if err != nil {
		logs.Error("sql error->", q, mid, err.Error())
		return lastid, res
	}
	defer rows.Close()

	//扫描数据
	for rows.Next() {
		dd := new(Movie)
		er := rows.Scan(&dd.Id, &dd.Number, &dd.SmallCover, &dd.BigCove, &dd.Trailer, &dd.Map)

		if er != nil {
			logs.Error("scan row error->", er.Error())
		}

		lastid = dd.Id

		res = append(res, dd)
	}

	return lastid, res
}

/**
* 总记录数
* param		restart		是否指定时间重写获取
 */
func (d *Movie) Total(restart bool, rTime string) int {
	res := 0

	q := `SELECT count(0) as nums FROM movie where is_up=1; `

	if restart == true {
		q = `SELECT count(0) as nums FROM movie where created_at >='` + rTime + `'; `
	}

	row, err := DB.Query(q)
	if err != nil {
		logs.Error("sql error->", q, err.Error())
		return res
	}
	defer row.Close()

	if row.Next() == true {
		row.Scan(&res)
	}
	return res
}

/**
* 更新数据
* @param	sId			影片id
* @param	RStatus 	是否上架,1=上架，2=下架
 */
func (d *Movie) Save(sId, RStatus string) bool {

	result, err := DB.Exec("update movie set is_up=? where id=? limit 1;", RStatus, sId)

	if err != nil {
		logs.Error("更新数据库错误！->", err.Error())
		return false
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return false
	}
	return true
}

func (d *Movie) Create() (err error) {
	if err = GetGormDb().Create(&d).Error; err != nil {
		return
	}
	return
}

/**
重复性检查
*/
func (d *Movie) Exists(number string) bool {
	var num int
	var q string
	q = `select count(1) as n from movie where number =%d`

	q = fmt.Sprintf(q, number)
	row, err := DB.Query(q)
	if err != nil {
		logs.Error("sql error->", q, err.Error())
		return false
	}

	defer row.Close()

	if row.Next() == true {
		row.Scan(&num)
	}

	if num > 0 {
		return true
	}

	return false
}
