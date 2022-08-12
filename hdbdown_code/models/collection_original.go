package models

import (
	"gorm.io/gorm"
)

/**
collection_original
`id` int unsigned NOT NULL AUTO_INCREMENT,
`oid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'mongoDb id',
`number` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '番号',
`db_name` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '所属mongodb表',
`dis_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '处理计数 为防止意外一个影片数据会处理两次',
`data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'json 数据',
`status` tinyint DEFAULT '1' COMMENT ' 1.未处理  2.已处理 3.需要重新处理  ',
`ctime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '采集方的创建时间 最大值用于从采集那边的筛选条件',
`utime` datetime DEFAULT NULL COMMENT '更新时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
*/
type CollectionOriginal struct {
	Id  int `json:"id" bson:"id" gorm:"primarykey"`
	Oid       string `json:"oid" bson:"oid"`
	Number    string `json:"number" bson:"number"`
	DbName    string `json:"db_name" bson:"db_name"`
	DisSum    int    `json:"dis_sum" bson:"dis_sum"`
	Data      string `json:"data" bson:"data"`
	Status    int    `json:"status" bson:"status"`
	Ctime     string `json:"ctime" bson:"ctime"`
	Utime     string `json:"utime" bson:"utime"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}

/**
指定表名
 */
func (CollectionOriginal) TableName() string {
	return "collection_original"
}


func (d *CollectionOriginal) Total() int64 {
	var total int64
	res := GetGormDb().Model(&d).Count(&total)
	if res.Error != nil {
		return 0
	}
	return total
}

func (d *CollectionOriginal) Lists(limit int, last int) (lastid int, data []*CollectionOriginal) {

	var rest *gorm.DB
	if last != 0  {
		rest = GetGormDb().Where("id > ?" , last).Limit(limit).Find(&data)
	}else{
		rest = GetGormDb().Limit(limit).Find(&data)
	}

	if rest.Error != nil {
		return
	}

	lastData := data[len(data)- 1]

	return lastData.Id, data
}
