package models

import (
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"time"
)

/**
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`mid` int unsigned NOT NULL DEFAULT '0' COMMENT '影片ID',
`uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID/回复的用户ID',
`cid` int unsigned NOT NULL DEFAULT '0' COMMENT '归属评论ID 0表示顶级评论',
`collection_id` int unsigned NOT NULL DEFAULT '0' COMMENT '来源采集ID - 只有source_type 为3时生效',
`comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '评论记录',
`score` float DEFAULT '0' COMMENT '评分0代表没有评分',
`type` tinyint DEFAULT '1' COMMENT '评论类型 1.评论  2.回复  ',
`source_type` tinyint DEFAULT '1' COMMENT '来源类型 1.用户评论 2.虚拟用户评论 3.采集评论 ',
`nickname` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '模拟 用户昵称 --目前只有采集的才会',
`avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '模拟 头像 -- 目前只有采集的才会',
`reply_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '回复的目标用户ID',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.删除  ',
`oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
`like` int unsigned NOT NULL DEFAULT '0' COMMENT '赞',
`dislike` int unsigned NOT NULL DEFAULT '0' COMMENT '踩',
`comment_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
`audit` tinyint DEFAULT '1' COMMENT '审核状态：1.正常 0.待审核 2.不通过',
`m_like` int NOT NULL DEFAULT '0' COMMENT '人工加赞',
`m_dislike` int NOT NULL DEFAULT '0' COMMENT '人工点踩',
影片评论表
movie_comment
*/
type MovieComment struct {
	base.Model
	Mid          int    `json:"mid" bson:"mid"`
	Uid          int    `json:"uid" bson:"uid"`
	Cid          int    `json:"cid" bson:"cid"`
	CollectionId int    `json:"collection_id" bson:"collection_id"`
	Comment      string `json:"comment" bson:"comment"`
	Score        string `json:"score" bson:"score"`
	Type         int    `json:"type" bson:"type"`
	SourceType   string `json:"source_type" bson:"source_type"`
	Nickname     string `json:"nickname" bson:"nickname"`
	Avatar       string `json:"avatar" bson:"avatar"`
	ReplyUid     string `json:"reply_uid" bson:"reply_uid"`
	Status       int    `json:"status" bson:"status"`
	Oid          string `json:"oid" bson:"oid"`
	Like         string `json:"'like'" bson:"'like'"`
	Dislike      string `json:"dislike" bson:"dislike"`
	CommentTime  string `json:"comment_time" bson:"comment_time"`
	Audit        string `json:"audit" bson:"audit"`
	MLike        string `json:"m_like" bson:"m_like"`
	MDislike     string `json:"m_dislike" bson:"m_dislike"`
}

/**
指定表名
*/
func (MovieComment) TableName() string {
	return "movie_comment"
}

func (d *MovieComment) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *MovieComment) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}
