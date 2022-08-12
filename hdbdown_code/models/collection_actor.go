package models

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
)

//同步演员对象
type CollectionActor struct {
	Id             int
	Name           string //演员名称
	Photo          string //演员照片
	Sex            string //演员性别
	SocialAccounts string //社交账户（保留） 数组

	MovieSum     int    //影片数量
	Category     string //类别，数组
	ActualSource string //验证网站
	Interflow    string //社交账户，数组
	Status       int    //状态

	ResourcesStatus  int    //资源下载状态
	ResourcesInfo    string //下载成功的数据，数组
	ResourcesOddInfo string //未成功的数据，数组
	Source           string //来源
	AdminId          int    //处理的管理员id

	OriginalId int    //源数据ID
	CreatedAt  string //创建时间
	UpdatedAt  string //更新时间
}

/**
* 从数据库中读取列表,每次读取10条
* mid			id
* updatedAt		最后更新时间
* limit 		每次读取多少条
 */
func (d *CollectionActor) Lists(mid, limit int, restart bool, rTime string) (lastid int, res []*CollectionActor) {

	q := `SELECT id,photo
		FROM collection_actor 
		where id>? and resources_status=1 and photo<>'' order by id asc limit %d;`

	if restart == true {
		q = `SELECT id,photo
		FROM collection_actor 
		where id>? and created_at>='` + rTime + `' and photo<>'' order by id asc limit %d;`
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
		dd := new(CollectionActor)
		er := rows.Scan(&dd.Id, &dd.Photo)

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
 */
func (d *CollectionActor) Total(restart bool, rTime string) int {
	res := 0

	q := `SELECT count(0) as nums FROM collection_actor where resources_status=1 and photo<>''; `

	if restart == true {
		q = `SELECT count(0) as nums FROM collection_actor where created_at >='` + rTime + `' and photo<>''; `
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
func (d *CollectionActor) Save(sId, RStatus, rInfo, ROinfo string) bool {

	result, err := DB.Exec("update collection_actor set resources_status=?,resources_info=?,resources_odd_info=? where id=? limit 1;", RStatus, rInfo, ROinfo, sId)

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
