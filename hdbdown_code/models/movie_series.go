package models

import (
	"gorm.io/gorm"
	"time"
)

const SeriesNameWithIDHash = "series_name_hash"

/**
`id` int unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '系列名称',
`movie_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '影片数量-冗余',
`like_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '收藏数量-冗余',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片系列表
movie_series
*/
type MovieSeries struct {
	Id        int    `json:"id" bson:"id" gorm:"primarykey"`
	Name      string `json:"name" bson:"name"`
	MovieSum  int    `json:"movie_sum" bson:"movie_sum"`
	LikeSum   int    `json:"like_sum" bson:"like_sum"`
	Status    int    `json:"status" bson:"status"`
	Oid       int    `json:"oid" bson:"oid"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}
func (d *MovieSeries) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}

/**
指定表名
*/
func (MovieSeries) TableName() string {
	return "movie_series"
}

func (ma *MovieSeries) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
