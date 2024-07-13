package service

import (
	"Project/define"
	"Project/helper"
	"Project/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Summary Get list of problems
// @Tags Problem
// @Accept  json
// @Produce  json
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
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
	categoryIdentity := c.Query("category_identity")

	list := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
	err = tx.Count(&count).Omit("content").Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
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

// GetProblemDetail
// @Summary Get Detail of problems
// @Tags Problem
// @Accept  json
// @Produce  json
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Get ProblemDetail Error" + err.Error(),
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// CreateProblem
// @Summary CreateProblem
// @Tags Admin
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Param max_runtime formData int true "max_runtime"
// @Param max_mem formData int true "max_mem"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /admin/problem-create [post]
func CreateProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	maxRuntime, err := strconv.Atoi(c.PostForm("max_runtime"))
	if err != nil {
		log.Println("MaxRuntime Parse Error:", err)
		return
	}
	maxMem, err := strconv.Atoi(c.PostForm("max_mem"))
	if err != nil {
		log.Println("MaxMem Parse Error:", err)
		return
	}
	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}
	identity := helper.GetUUID()
	data := models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}

	problemCategories := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		i, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Atoi Error:", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}
		problemCategories = append(problemCategories, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(i),
		})
	}
	data.ProblemCategories = problemCategories

	testCaseArr := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if caseMap["input"] == "" || caseMap["output"] == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "输入或输出为空",
			})
			return
		}
		testCaseIdentity := helper.GetUUID()
		testCaseArr = append(testCaseArr, &models.TestCase{
			Identity:        testCaseIdentity,
			ProblemIdentity: data.Identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		})
	}
	data.TestCases = testCaseArr

	err = models.DB.Create(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Create Error" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data.Identity,
	})
}

// ModifyProblem
// @Summary ModifyProblem
// @Tags Admin
// @Param authorization header string true "authorization"
// @Param problem body models.ProblemBasic true "problem"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /admin/problem-modify [put]
func ModifyProblem(c *gin.Context) {
	// @Param problem body models.ProblemBasic true "problem"
	var problem models.ProblemBasic
	if err := c.BindJSON(&problem); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Data binding error",
		})
		return
	}

	if problem.Title == "" || problem.Identity == "" || problem.Content == "" || problem.MaxRuntime == 0 ||
		problem.MaxMem == 0 || len(problem.ProblemCategories) == 0 || len(problem.TestCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Element is empty",
		})
		return
	}

	//problem关联分类,problem关联测试用例
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		//删除原有关联
		err := tx.Where("problem_id = ?", problem.ID).Delete(&models.ProblemCategory{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("problem_identity = ?", problem.Identity).Delete(&models.TestCase{}).Error
		if err != nil {
			return err
		}
		//创建新关联
		for _, category := range problem.ProblemCategories {
			err = tx.Create(&category).Error
			if err != nil {
				return err
			}
		}
		for _, testCase := range problem.TestCases {
			err = tx.Create(&testCase).Error
			if err != nil {
				return err
			}
		}
		//更新problem
		err = tx.Model(&models.ProblemBasic{}).Where("identity = ?", problem.Identity).Updates(&problem).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Modify Problem Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Modify Problem Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Modify Problem Success",
	})
}
