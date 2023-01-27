package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"time"
)

type DownloadError struct {
	Title   string `json:"title" bson:"title"`
	Message string `json:"message" bson:"message"`
}

/*
CREATE TABLE `movie_error` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mid` int(11) DEFAULT NULL,
  `message` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
MovieError
*/
type MovieError struct {
	base.Model
	Mid     int    `json:"mid" bson:"mid"`
	Message string `json:"message" bson:"message"`
}

/**
指定表名
*/
func (MovieError) TableName() string {
	return "movie_error"
}

func (d *MovieError) Create() (err error) {
	err = orm.Eloquent.Create(&d).Error
	return
}

func (ma *MovieError) BeforeCreate(tx *gorm.DB) (err error) {
	ma.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ma.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return
}

func (me *MovieError) Add(Mid int, title string, message string) (err error) {

	if err := orm.Eloquent.Where("mid =?", Mid).First(&me).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	var errs []DownloadError

	if me.Mid == Mid {
		if me.Message != "" {
			err = json.Unmarshal([]byte(me.Message), &errs)
			if err != nil {
				return
			}
		}

		isHas := false
		for _, v := range errs {
			if v.Title == title {
				isHas = true
				v.Message = message
			}
		}

		if isHas == false {
			var er = DownloadError{
				Title:   title,
				Message: message,
			}
			errs = append(errs, er)
		}

		resultError, err := json.Marshal(errs)
		if err != nil {
			return err
		}

		me.Message = string(resultError)

		if err := orm.Eloquent.Model(&me).Where("mid =?", Mid).Update("message", resultError).Error; err != nil {
			return err
		}
	} else {
		var er = DownloadError{
			Title:   title,
			Message: message,
		}
		errs = append(errs, er)

		resultError, err := json.Marshal(errs)
		if err != nil {
			return err
		}

		me.Mid = Mid
		me.Message = string(resultError)
		if err := me.Create(); err != nil {
			return err
		}
	}

	return

}
