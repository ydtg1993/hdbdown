package models

import (
	"gorm.io/gorm"
	"hdbdown/global/orm"
	"hdbdown/models/base"
	"hdbdown/tools/rd"
)

/**
CREATE TABLE `admin_config` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_config_name_unique` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
*/
type AdminConfig struct {
	base.Model
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (AdminConfig) TableName() string {
	return "admin_config"
}

const AdminConfigLists = "admin_config_lists"

func (a AdminConfig) GetValue(key string) (err error, data string) {
	data, err = rd.GetCashWithHash(AdminConfigLists, key, func() (s string, err error) {
		var config AdminConfig
		if err = orm.Eloquent.Model(AdminConfig{}).Where("name =?", key).First(&config).Error; err != nil && err != gorm.ErrRecordNotFound {
			return
		}
		s = config.Value
		return
	})
	return
}
