package models

import (
	"gorm.io/gorm"
	"time"
)

/*
movie_series_associate 影片系列关联表
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`series_id` int unsigned NOT NULL DEFAULT '0' COMMENT '导演ID',
`mid` int unsigned NOT NULL DEFAULT '0' COMMENT '影片ID',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`associate_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关联时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
*/
type MovieSeriesAssociate struct{
	Id  int `json:"id" bson:"id" gorm:"primarykey"`
	SeriesId int `json:"series_id" bson:"series_id"`
	Mid int `json:"mid" bson:"mid"`
	Status int `json:"status" bson:"status"`
	AssociateTime string `json:"associate_time" bson:"associate_time"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}

func (d *MovieSeriesAssociate) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}

/**
指定表名
 */
func (MovieSeriesAssociate) TableName() string {
	return "movie_series_associate"
}

func (ma *MovieSeriesAssociate) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}


