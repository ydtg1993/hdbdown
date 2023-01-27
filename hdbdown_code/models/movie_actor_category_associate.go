package models

import (
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"time"
)

/**
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`cid` int unsigned NOT NULL DEFAULT '0' COMMENT '类别ID',
`aid` int unsigned NOT NULL DEFAULT '0' COMMENT '演员ID',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`associate_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关联时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片演员类别关联表
movie_actor_category_associate
*/
type MovieActorCategoryAssociate struct {
	base.Model
	Cid           int    `json:"cid" bson:"cid"`
	Aid           int    `json:"aid" bson:"aid"`
	Status        int    `json:"status" bson:"status"`
	AssociateTime string `json:"associate_time" bson:"associate_time"`
}

/**
指定表名
*/
func (MovieActorCategoryAssociate) TableName() string {
	return "movie_actor_category_associate"
}

func (d *MovieActorCategoryAssociate) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *MovieActorCategoryAssociate) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
