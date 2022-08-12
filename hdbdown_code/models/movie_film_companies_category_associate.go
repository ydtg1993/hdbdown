package models

import (
	"gorm.io/gorm"
	"time"
)

/*
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `cid` int unsigned NOT NULL DEFAULT '0' COMMENT '类别ID',
  `film_companies_id` int unsigned NOT NULL DEFAULT '0' COMMENT '片商ID',
  `status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
  `associate_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关联时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
movie_film_companies_category_associate
影片片商类别关联表
*/
type MovieFilmCompaniesCategoryAssociate struct {
	Id              int    `json:"id" bson:"id" gorm:"primarykey"`
	Cid             int    `json:"cid" bson:"cid"`
	FilmCompaniesId int    `json:"film_companies_id" bson:"film_companies_id"`
	Status          int    `json:"status" bson:"status"`
	AssociateTime   string `json:"associate_time" bson:"associate_time"`
	CreatedAt       string `json:"created_at" bson:"created_at"`
	UpdatedAt       string `json:"updated_at" bson:"updated_at"`
}

/**
指定表名
*/
func (MovieFilmCompaniesCategoryAssociate) TableName() string {
	return "movie_film_companies_category_associate"
}

func (ma *MovieFilmCompaniesCategoryAssociate) BeforeCreate(tx *gorm.DB) (err error) {
	ma.Status = 1
	ma.AssociateTime = time.Now().Format("2006-01-02 15:04:05")
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (d *MovieFilmCompaniesCategoryAssociate) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}
