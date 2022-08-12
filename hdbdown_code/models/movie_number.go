package models

/**
`id` int unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '番号名称',
`movie_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '影片数量',
`like_sum` int unsigned NOT NULL DEFAULT '0' COMMENT '收藏数量-冗余',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`oid` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据ID',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片番号表
 */
type MovieNumber struct {
	Id        int    `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	MovieSum  int    `json:"movie_sum" bson:"movie_sum"`
	LikeSum   int    `json:"like_sum" bson:"like_sum"`
	Status    int    `json:"status" bson:"status"`
	Oid       int    `json:"oid" bson:"oid"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}
