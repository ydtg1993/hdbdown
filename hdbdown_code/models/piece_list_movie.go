package models


/*
piece_list_movie
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `plid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '片单ID',
  `mid` int unsigned NOT NULL DEFAULT '0' COMMENT '影片ID',
  `status` tinyint DEFAULT '1' COMMENT ' 1.正常  2.删除  ',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
*/
type PieceListMovie struct{
	Id         int    `json:"id" bson:"id"`
	Mid        int    `json:"mid" bson:"mid"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
	//UpdatedAt  string `json:"updated_at" bson:"updated_at"`
	Status    int `json:"status" bson:"status"`
	Plid      int `json:"plid" bson:"plid"`
	//UpdatedAt int `json:"updated_at" bson:""`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	
}