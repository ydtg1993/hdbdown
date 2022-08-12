package models

import (
	"gorm.io/gorm"
	"time"
)

/**
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`cid` int unsigned NOT NULL DEFAULT '0' COMMENT '类别ID',
`mid` int unsigned NOT NULL DEFAULT '0' COMMENT '影片ID',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`associate_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关联时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片类别关联表
 */
type MovieCategoryAssociate struct {
	Id            int    `json:"id" bson:"id" gorm:"primarykey"`
	Cid           int    `json:"cid" bson:"cid"`
	Mid           int    `json:"mid" bson:"mid"`
	Status        int    `json:"status" bson:"status"`
	AssociateTime string `json:"associate_time" bson:"associate_time"`
	CreatedAt     string `json:"created_at" bson:"created_at"`
	UpdatedAt     string `json:"updated_at" bson:"updated_at"`
}

/**
指定表名
 */
func (MovieCategoryAssociate) TableName() string {
	return "movie_category_associate"
}

func (d *MovieCategoryAssociate) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}

func (ma *MovieCategoryAssociate) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
