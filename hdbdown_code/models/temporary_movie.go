package models

import (
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"time"
)

//  `status` tinyint DEFAULT '1' COMMENT ' 1.未处理  2.已处理 3.数据异常',
const StatusProcessed = 2
const StatusNotProcessed = 1
const StatusNotUnusual = 3

// `is_update` tinyint DEFAULT '1' COMMENT '是否更新, 1 不需要 , 2 需要更新',
const NeedUpdate = 2
const NoNeedUpdate = 1

/**
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `db_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据源',
  `number` char(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '番号',
  `status` tinyint DEFAULT '1' COMMENT ' 1.未处理  2.已处理 3.需要重新处理  ',
  `is_update` tinyint DEFAULT '1' COMMENT '是否更新, 1 不需要 , 2 需要更新',
  `ctime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '采集方的创建时间 最大值用于从采集那边的筛选条件',
  `utime` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
temporary_movie
采集原始数据表
*/
type TemporaryMovie struct {
	base.Model
	Number   string `json:"number" bson:"number"`
	DbName   string `json:"db_name" bson:"db_name"`
	Status   int    `json:"status" bson:"status"`
	IsUpdate int    `json:"is_update" bson:"is_update"`
	Ctime    string `json:"ctime" bson:"ctime"`
	Utime    string `json:"utime" bson:"utime"`
}

func (d *TemporaryMovie) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

/**
指定表名
*/
func (TemporaryMovie) TableName() string {
	return "temporary_movie"
}

func (ma *TemporaryMovie) BeforeCreate(tx *gorm.DB) (err error) {
	if ma.Ctime == "" {
		ma.Ctime = time.Now().Format("2006-01-02 15:04:05")
	}

	if ma.Utime == "" {
		ma.Utime = ma.Ctime
	}

	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *TemporaryMovie) BeforeUpdate(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *TemporaryMovie) BeforeSave(tx *gorm.DB) (err error) {
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (ma *TemporaryMovie) GetDataByNumber(number string) (err error) {
	err = orm.Eloquent.Where("number = ?", number).First(&ma).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

/**
标记处理完成
*/
func (ma *TemporaryMovie) PressSuccess() error {
	rest := orm.Eloquent.Model(&ma).Where("id = ?", ma.Id).Update("status", StatusProcessed)
	if rest.Error != nil {
		return rest.Error
	}
	return nil
}

func (ma *TemporaryMovie) IncompleteData() error {
	rest := orm.Eloquent.Model(&ma).Where("id = ?", ma.Id).Update("status", StatusNotUnusual)
	if rest.Error != nil {
		return rest.Error
	}
	return nil
}

/**
统计临时表的数量
*/
func (ma *TemporaryMovie) ListOfNotProcessedCount() (err error, total int64) {

	res := orm.Eloquent.Model(ma).Where("status = ?", StatusNotProcessed).Count(&total)
	if res.Error != nil {
		err = res.Error
		return
	}
	return
}

/**
获取第一条数据的 ID
*/
func (ma *TemporaryMovie) ListOfNotProcessedLastId() (err error, lastId int) {
	var temp TemporaryMovie
	res := orm.Eloquent.Model(ma).Where("status = ?", StatusNotProcessed).Order("id asc").Limit(1).Select([]string{"id"}).Find(&temp)
	if res.Error != nil {
		err = res.Error
		return
	}

	return
}

/**
查询 status = 1 的数据
*/
func (ma *TemporaryMovie) ListOfNotProcessed(lastId int, limit int) (err error, last int, data []*TemporaryMovie) {
	// select * from temporary_movie where status = 1 and id > lastId order by id asc limit 500
	if lastId == 0 {
		err = orm.Eloquent.Model(ma).Where("status = ?", StatusNotProcessed).Order("id asc").Limit(limit).Find(&data).Error
	} else {
		err = orm.Eloquent.Model(ma).Where("status = ? and id > ?", StatusNotProcessed, lastId).Order("id asc").Limit(limit).Find(&data).Error
	}
	if err != nil {
		return
	}

	last = lastId
	if len(data) > 0 {
		obj := data[len(data)-1]
		last = obj.Id
	}
	return
}

func (ma *TemporaryMovie) ListOfNeedUpdateCount() (err error, total int64) {
	err = orm.Eloquent.Model(ma).Where("is_update = ? and status = ?", NeedUpdate, StatusProcessed).Count(&total).Error
	if err != nil {
		return
	}
	return
}

func (ma *TemporaryMovie) ListOfNeedUpdate(lastId int, limit int) (err error, last int, data []*TemporaryMovie) {
	// select * from temporary_movie where is_update = 2 and id > lastId and status = 2 order by id asc limit 500
	if lastId == 0 {
		err = orm.Eloquent.Model(ma).Where("is_update = ?", NeedUpdate).Order("id asc").Limit(limit).Find(&data).Error
	} else {
		err = orm.Eloquent.Model(ma).Where("is_update = ? and id > ? and status = ?", NeedUpdate, lastId, StatusProcessed).Order("id asc").Limit(limit).Find(&data).Error
	}

	if err != nil {
		return
	}

	last = lastId
	if len(data) > 0 {
		obj := data[len(data)-1]
		last = obj.Id
	}
	return
}

func (ma *TemporaryMovie) ListOfAll(lastId int, limit int) (err error, last int, data []TemporaryMovie) {
	// select * from temporary_movie where id > lastId order by id asc limit 500
	if lastId == 0 {
		err = orm.Eloquent.Model(ma).Order("id asc").Limit(limit).Find(&data).Error
	} else {
		err = orm.Eloquent.Model(ma).Where("id > ?", lastId).Order("id asc").Limit(limit).Find(&data).Error
	}
	if err != nil {
		return
	}

	last = lastId
	if len(data) > 0 {
		obj := data[len(data)-1]
		last = obj.Id
	}
	return
}
