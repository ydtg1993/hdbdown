package models

import (
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"hdbdown/tools/rd"
	"time"
)

// 演员处理 redis 队列 key
const MovieActorPress = "movie_actor_press"
const ActorNameWithIDHash = "actor_name_hash"

/**
`id` int unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '演员名称',
`photo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '演员照片',
`sex` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '演员性别',
`oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
`social_accounts` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'json 数组 社交账户',
`movie_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '影片数量-冗余',
`like_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '收藏数量-冗余',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片演员表
movie_actor
*/
type MovieActor struct {
	base.Model
	Name           string `json:"name" bson:"name"`
	Photo          string `json:"photo" bson:"photo"`
	Sex            string `json:"sex" bson:"sex"`
	Oid            int    `json:"oid" bson:"oid"`
	SocialAccounts string `json:"social_accounts" bson:"social_accounts"`
	MovieSum       int    `json:"movie_sum" bson:"movie_sum"`
	LikeSum        int    `json:"like_sum" bson:"like_sum"`
	Status         int    `json:"status" bson:"status"`
}

/**
指定表名
*/
func (MovieActor) TableName() string {
	return "movie_actor"
}

func (d *MovieActor) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *MovieActor) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *MovieActor) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *MovieActor) AfterCreate(tx *gorm.DB) (err error) {
	err = rd.RPush(MovieActorPress, ma.Name)
	if err != nil {
		return
	}
	return
}
