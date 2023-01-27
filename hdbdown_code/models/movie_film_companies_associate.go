package models

import (
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"time"
)

/**
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`film_companies_id` int unsigned NOT NULL DEFAULT '0' COMMENT '导演ID',
`mid` int unsigned NOT NULL DEFAULT '0' COMMENT '影片ID',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`associate_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关联时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片片商关联表
movie_film_companies_associate
*/
type MovieFilmCompaniesAssociate struct {
	base.Model
	Mid             int    `json:"mid" bson:"mid"`
	FilmCompaniesId int    `json:"film_companies_id" bson:"film_companies_id"`
	Status          int    `json:"status" bson:"status"`
	AssociateTime   string `json:"associate_time" bson:"associate_time"`
}

/**
指定表名
*/
func (MovieFilmCompaniesAssociate) TableName() string {
	return "movie_film_companies_associate"
}

func (d *MovieFilmCompaniesAssociate) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *MovieFilmCompaniesAssociate) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
