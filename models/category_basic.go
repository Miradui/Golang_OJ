package models

import "gorm.io/gorm"

type CategoryBasic struct {
	gorm.Model
	Identity string `gorm:"cloumn:identity;type:varchar(36);" json:"identity"`
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`  // 分类名称
	ParentId int    `gorm:"column:parent_id;type:int;" json:"parent_id"` //
}

func (table *CategoryBasic) TableName() string {
	return "category_basic"
}
