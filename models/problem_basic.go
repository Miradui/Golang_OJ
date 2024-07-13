package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"colum:identity;type:varchar(36);" json:"identity"`
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	Title             string             `gorm:"cloum:title;type:varchar(255);" json:"title"`
	Content           string             `gorm:"cloum:content;type:text;" json:"content"`
	MaxRuntime        int                `gorm:"cloumn:max_runtime;type:int(11);" json:"max_runtime"`
	MaxMem            int                `gorm:"cloumn:max_mem;type:int(11);" json:"max_mem"`
	TestCases         []*TestCase        `gorm:"foreignKey:problem_identity;references:identity"`
}

func (table ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ? ", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )", categoryIdentity)
	}
	return tx
}
