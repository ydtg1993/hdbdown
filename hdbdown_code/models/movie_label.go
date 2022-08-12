package models

import (
	"gorm.io/gorm"
	"time"
)

const LabelNameWithIDHash = "Label_name_hash"

/**
`id` int unsigned NOT NULL AUTO_INCREMENT,
`cid` int unsigned NOT NULL DEFAULT '0' COMMENT '父标签ID 0 表示顶级标签',
`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '标签名称',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
`sort` int DEFAULT '0' COMMENT '排序，正序',
`item_num` int DEFAULT '0' COMMENT '下一级的数量',
`like_sum` int DEFAULT '0' COMMENT '收藏数量',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
movie_label
影片标签表
*/
type MovieLabel struct {
	Id        int    `json:"id" bson:"id" gorm:"primarykey"`
	Cid       int    `json:"cid" bson:"cid"`
	Name      string `json:"name" bson:"name"`
	Status    int    `json:"status" bson:"status"`
	Oid       int    `json:"oid" bson:"oid"`
	Sort      int    `json:"sort" bson:"sort"`
	ItemNum   int    `json:"item_num" bson:"item_num"`
	LikeSum   int    `json:"like_sum" bson:"like_sum"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}

func (d *MovieLabel) Create() (err error) {
	err = GetGormDb().Create(&d).Error
	return
}

func (ma *MovieLabel) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return
}

func (MovieLabel) TableName() string {
	return "movie_label"
}
