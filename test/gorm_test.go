package test

import (
	"Project/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormTest(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/online_judge?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(nil)
	}
	data := make([]*models.ProblemBasic, 0)
	err = db.Find(&data).Error
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range data {
		fmt.Printf("Problem ==> %v \n\n", v)
	}
}

func TestGormTest1(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/online_judge?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&models.TestCase{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGormTest2(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/online_judge?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// 插入CategoryBasic数据
	//category := &models.CategoryBasic{
	//	Identity: "category1",
	//	Name:     "Category 1",
	//}
	//err = db.Create(category).Error
	//if err != nil {
	//	t.Fatal(err)
	//}

	// 插入ProblemCategory数据
	//problemCategory := &models.ProblemCategory{
	//	ProblemId:  "1",
	//	CategoryId: "2",
	//}
	//err = db.Create(problemCategory).Error
	//if err != nil {
	//	t.Fatal(err)
	//}

	//插入Category数据
	//	categoryBasic := &models.CategoryBasic{
	//	Identity: "category2",
	//	Name:     "Category2",
	//	ParentId: 0,
	//}

	//插入user数据
	//data := &models.UserBasic{
	//	Name:     "name1",
	//	Password: "password1",
	//	Phone:    "phone1",
	//	Email:    "email1@email.com",
	//}

	//插入submit数据\
	data := &models.SubmitBasic{
		Identity:        "submit1",
		ProblemIdentity: "1",
		UserIdentity:    "user1",
		Path:            "lkasdj/lkjl",
		Status:          "-1",
	}
	err = db.Create(data).Error
	if err != nil {
		t.Fatal(err)
	}
}
