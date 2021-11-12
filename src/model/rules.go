package model

import (
	_ "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//todo test

type Rule struct {
	Aid                  int    `json:"aid" form:"aid" binding:"required" gorm:"Column:aid"`
	Platform             string `form:"platform" binding:"required" json:"platform" db:"platform" gorm:"Column:platform; Type:varchar(128)"`
	DownloadUrl          string `form:"download_url" binding:"required" json:"download_url" db:"download_url" gorm:"Column:download_url; Type:varchar(1024)" `
	UpdateVersionCode    string `form:"update_version_code" binding:"required" json:"update_version_code" gorm:"Column:update_version_code; Type:varchar(128)"`
	Md5                  string `form:"md_5" binding:"required" json:"md_5" gorm:"Column:md5; Type:varchar(128)"`
	DeviceIdList         string `form:"device_id_list" binding:"required" json:"device_id_list" gorm:"Column:device_id_list; Type:text"`
	MaxUpdateVersionCode string `form:"max_update_version_code" binding:"required" json:"max_update_version_code" gorm:"Column:max_update_version_code; Type:varchar(128)"`
	MinUpdateVersionCode string `form:"min_update_version_code" binding:"required" json:"min_update_version_code" gorm:"Column:min_update_version_code; Type:varchar(128)"`
	MaxOsApi             int    `form:"max_os_api" binding:"required" json:"max_os_api" gorm:"Column:max_os_api"`
	MinOsApi             int    `form:"min_os_api" binding:"required" json:"min_os_api" gorm:"Column:min_os_api"`
	CpuArch              string `form:"cpu_arch" binding:"required" json:"cpu_arch" gorm:"Column:cpu_arch; Type:varchar(32)"`
	Channel              string `form:"channel" binding:"required" json:"channel" gorm:"Column:channel; Type:varchar(128)"`
	Title                string `form:"title" binding:"required" json:"title" gorm:"Column:title; Type:varchar(256)"`
	UpdateTips           string `form:"update_tips" binding:"required" json:"update_tips" gorm:"Column:update_tips; Type:varchar(1024)"`
}

type NewRule struct {
	gorm.Model
	Rule
	Enable       bool        `gorm:"Column: enable"`
	DeviceIdList []WhiteList `gorm:"many2many:rule2device_list;"`
}

func (NewRule) TableName() string {
	return "rules"
}

func rule2newRule(rule Rule) NewRule {
	whiteList := ParseWhiteList(rule.DeviceIdList)
	return NewRule{Rule: rule, DeviceIdList: *whiteList, Enable: true}
}

func GetRules(deviceId string) *[]Rule {
	Mysql()
	defer Db.Close()

	var newRules = make([]NewRule, 0)
	_ = Db.Preloads(&WhiteList{}).Find(&newRules).Where("device_name = ?", deviceId)
	rules := make([]Rule, 0)
	for index := 0; index < len(newRules); index++ {
		rules = append(rules, newRules[index].Rule)
	}
	return &rules
}
func AddRule(rule Rule) uint {
	Mysql()
	defer Db.Close()

	Db.AutoMigrate(&NewRule{})
	Db.AutoMigrate(&WhiteList{})
	var newRule = rule2newRule(rule)
	if err := Db.Create(&newRule).Error; err != nil {
		return 0
	}
	return newRule.Model.ID
}
