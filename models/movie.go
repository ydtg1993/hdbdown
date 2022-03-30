package models

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
)

//同步电影对象
type Movie struct {
	Id           int
	Number       string //番号
	NumberSource string //源番号
	Name         string //影片名称
	Sell         string //卖家

	Time        string //播放时长（秒）
	ReleaseTime string //发布时间
	SmallCover  string //小封面
	BigCove     string //大封面
	Trailer     string //预告片

	Map       string //组图，数组
	Oid       string //源id
	Cid       string //类别id
	CreatedAt string //创建时间
	UpdatedAt string //更新时间
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
