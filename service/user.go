package service

import (
	"Project/define"
	"Project/helper"
	"Project/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetUserDetail
// @Summary Get Detail of Users
// @Tags User
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}
	data := new(models.UserBasic)
	err := models.DB.Where("identity = ?", identity).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Detail By Identity:" + identity + "Error: " + err.Error(),
		})
		return
	}
	if data.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})

}

// Login
// @Summary Login
// @Tags User
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}

	password = helper.GetMd5(password)

	data := new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ? ", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名患者密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get UserBasic Error:" + err.Error(),
		})
		return
	}

	token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// SendCode
// @Summary SendCode
// @Tags User
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	fmt.Println(email)
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱为空",
		})
		return
	}
	code := helper.GenerateVerificationCode()
	models.RDB.Set(c, email, code, time.Second*300)
	err := helper.SendCodeByEmail(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Send Code Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

// Register
// @Summary Register
// @Tags User
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string true "phone"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	if email == "" || code == "" || name == "" || password == "" || phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}

	// 验证码验证
	redisCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码已过期",
		})
		return
	}
	if redisCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}

	//判断用户是否已经存在
	var cnt int64
	err = models.DB.Where("email = ?", email).Model(new(models.UserBasic)).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Error:" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
		})
		return
	}

	//插入数据
	data := &models.UserBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Email:    email,
	}

	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create User Error:" + err.Error(),
		})
		return
	}

	//生成token
	token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// GetRankList
// @Summary Get RankList
// @Tags User
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
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
	var list = make([]*models.UserBasic, 0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("pass_num DESC, submit_num ASC").
		Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Get Rank List Error:" + err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
