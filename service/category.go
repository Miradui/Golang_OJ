package service

import (
	"Project/define"
	"Project/helper"
	"Project/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Summary Get list of category
// @Tags Admin
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param authorization header string true "authorization"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /admin/category-list [get]
func GetCategoryList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Page Parse Error:", err)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("Size Parse Error:", err)
		return
	}
	offset := (page - 1) * size
	var count int64

	keyword := c.Query("keyword")

	list := make([]*models.CategoryBasic, 0)
	tx := models.DB.Model(new(models.CategoryBasic)).Where("name like ? ", "%"+keyword+"%")
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Category List Error:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// CreateCategory
// @Summary CreateCategory
// @Tags Admin
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parent_id formData string true "parent_id"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /admin/problem-create [post]
func CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	parentId, err := strconv.Atoi(c.PostForm("parent_id"))
	if err != nil {
		log.Println("ParentId Parse Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ParentId Parse Error",
		})
		return
	}
	category := models.CategoryBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		ParentId: parentId,
	}
	err = models.DB.Create(&category).Error
	if err != nil {
		log.Println("Create Category Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create Category Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Create Category Success",
	})

}

// ModifyCategory
// @Summary ModifyCategory
// @Tags Admin
// @Param authorization header string true "authorization"
// @Param category body models.CategoryBasic true "category"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /admin/problem-modify [put]
func ModifyCategory(c *gin.Context) {
	var category models.CategoryBasic
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Data binding error",
		})
		return
	}

	if category.Name == "" || category.Identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Name or Identity is empty",
		})
		return
	}

	err := models.DB.Model(&models.CategoryBasic{}).Where("identity = ?", category.Identity).Updates(&category).Error
	if err != nil {
		log.Println("Modify Category Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Modify Category Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Modify Category Success",
	})
}

// DeleteCategory
// @Summary DeleteCategory
// @Tags Admin
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /admin/problem-delete [delete]
func DeleteCategory(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Identity is empty",
		})
		return
	}

	tx := models.DB.Begin()

	// ProblemCategory
	var count int64
	err := tx.Where("category_id = {SELECT id FROM category_basic WHERE identity = ? LIMIT 1 }",
		identity).Count(&count).Delete(&models.ProblemCategory{}).Error
	if err != nil {
		tx.Rollback()
		log.Println("Delete ProblemCategory Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Delete ProblemCategory Error",
		})
		return
	}
	if count > 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Category is used",
		})
		return
	}

	// Category
	err = tx.Where("identity = ?", identity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		tx.Rollback()
		log.Println("Delete Category Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Delete Category Error",
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Delete Category Success",
	})
}
