package models

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
)

//同步电影对象
type CollectionMovie struct {
	Id           int
	Number       string //番号
	NumberSource string //源番号
	NumberName   string //番号名称
	Name         string //影片名称

	SourceSite string //来源网站
	SourceUrl  string //来源路径
	Director   string //导演名称
	Sell       string //卖家
	Time       string //播放时长（秒）

	ReleaseTime string //发布时间
	SmallCover  string //小封面
	BigCove     string //大封面
	Trailer     string //预告片
	Map         string //组图，数组

	Series        string //系列
	FilmCompanies string //片商
	Issued        string //发行
	Actor         string //演员，数组
	Category      string //类别

	Label       string  //标签，数组
	Score       float64 //积分
	ScorePeople int     //评分人数
	CommentNum  int     //评论数
	Comment     string  //评论，数组

	ActualSource   string //验证网址
	FluxLinkageNum int    //磁链数量
	FluxLinkage    string //磁链,数组
	IsDownload     int    //是否可下载
	IsSubtitle     int    //是否包含字幕

	IsNew           int    //是否最新
	AbnormalDataId  string //异常数据id组，数组
	Status          int    //状态
	ResourcesStatus int    //资源下载状态
	ResourcesInfo   string //下载成功的数据，数组

	ResourcesOddInfo string //未成功的数据，数组
	AdminId          int    //处理的管理员id
	DisSum           int    //处理计数
	OriginalId       int    //源数据ID
	Ctime            string //同步爬取时间

	Utime     string //同步更新时间
	CreatedAt string //创建时间
	UpdatedAt string //更新时间
}

/**
* 从数据库中读取列表,每次读取10条
* mid			id
* updatedAt		最后更新时间
* limit 		每次读取多少条
 */
func (d *CollectionMovie) Lists(mid, limit int, restart bool, rTime string) (lastid int, res []*CollectionMovie) {

	q := `SELECT id,small_cover,big_cove,trailer,map
		FROM collection_movie 
		where id>? and resources_status=1 order by id asc limit %d;`

	if restart == true {
		q = `SELECT id,small_cover,big_cove,trailer,map
		FROM collection_movie 
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
		dd := new(CollectionMovie)
		er := rows.Scan(&dd.Id, &dd.SmallCover, &dd.BigCove, &dd.Trailer, &dd.Map)

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
func (d *CollectionMovie) Total(restart bool, rTime string) int {
	res := 0

	q := `SELECT count(0) as nums FROM collection_movie where resources_status=1; `

	if restart == true {
		q = `SELECT count(0) as nums FROM collection_movie where created_at >='` + rTime + `'; `
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
 */
func (d *CollectionMovie) Save(sId, RStatus, rInfo, ROinfo string) bool {

	result, err := DB.Exec("update collection_movie set resources_status=?,resources_info=?,resources_odd_info=? where id=? limit 1;", RStatus, rInfo, ROinfo, sId)

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

/**
* 通过番号，修改状态
 */
func (d *CollectionMovie) SaveWithNumber(number, RStatus string) bool {

	result, err := DB.Exec("update collection_movie set resources_status=?,status=1 where number=? limit 1;", RStatus, number)

	if err != nil {
		logs.Error("更新数据库错误！->", RStatus, number, err.Error())
		return false
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return false
	}
	return true
}
