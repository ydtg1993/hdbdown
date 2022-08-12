package models

/**
`id` bigint unsigned NOT NULL AUTO_INCREMENT,
`nid` int unsigned NOT NULL DEFAULT '0' COMMENT '番号ID',
`mid` int unsigned NOT NULL DEFAULT '0' COMMENT '影片ID',
`status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.弃用  ',
`associate_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '关联时间',
`created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
影片番号关联表
*/
type MovieNumberAssociate struct {
	Id            int    `json:"id" bson:"id"`
	Nid           int    `json:"nid" bson:"nid"`
	Mid           int    `json:"mid" bson:"mid"`
	Status        int    `json:"status" bson:"status"`
	AssociateTime string `json:"associate_time" bson:"associate_time"`
	CreatedAt     string `json:"created_at" bson:"created_at"`
	UpdatedAt     string `json:"updated_at" bson:"updated_at"`
}
