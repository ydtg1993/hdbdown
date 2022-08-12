package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"hdbdown/rd"
	"strings"
	"time"
)

const CategoryList = "category_list"

/**
`id` int unsigned NOT NULL AUTO_INCREMENT,
 `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '类别名称',
 `status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
 `oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
 `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
 `show` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0显示 1不显示',
 `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
*/
type MovieCategory struct {
	Id        int    `json:"id" bson:"id" gorm:"primarykey"`
	Name      string `json:"name" bson:"name"`
	Status    int    `json:"status" bson:"status"`
	Oid       string `json:"oid" bson:"oid"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	UpdatedAt string `json:"UpdatedAt" bson:"UpdatedAt"`
}

/**
指定表名
*/
func (MovieCategory) TableName() string {
	return "movie_category"
}

func (d *MovieCategory) FirstByName(name string) (err error) {
	name = d.SwitchName(name)
	res := GetGormDb().Where("name = ?", name).First(&d)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err = res.Error
		return
	}
	return
}

func (d *MovieCategory) SwitchName(name string) string {
	switch name {
	case "有碼":
		return "有码"
	case "無碼":
		return "无码"
	case "國產":
		return "国产"
	case "歐美":
		return "欧美"
	default:
		return strings.ToUpper(name)
	}
}

func (d *MovieCategory) FindWithRedis() (err error, list []*MovieCategory) {

	data, err := rd.GetCash(CategoryList, func() string {

		if err = GetGormDb().Where("status = ?", 1).Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
			return ""
		}

		jsonStr, err := json.Marshal(list)
		if err != nil {
			return ""
		}

		return string(jsonStr)
	}, time.Second*3600)

	if data != "" {
		err = json.Unmarshal([]byte(data), &list)
		if err != nil {
			return
		}
		return
	}

	return
}

func (d *MovieCategory) FindByNameWithRedis(name string) (err error, data *MovieCategory) {
	name = d.SwitchName(name)
	err, list := d.FindWithRedis()
	if err != nil {
		return
	}

	for _, v := range list {
		if v.Name == name {
			return nil, v
		}
	}

	return nil, list[0]
}
