package router

import (
	_ "Project/docs"
	"Project/middlewares"
	"Project/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Roter() *gin.Engine {
	r := gin.Default()

	//swag 配置
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//路由规则

	// 问题
	r.GET("/ping", service.Ping)
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	//用户
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)

	//排行榜
	r.GET("/rank-list", service.GetRankList)

	//提交记录
	r.GET("/submit-list", service.GetSubmitList)

	//管理员私有方法
	adminGroup := r.Group("/admin", middlewares.AuthAdminCheck())
	//问题
	adminGroup.POST("/problem-create", service.CreateProblem)
	adminGroup.PUT("/problem-modify", service.ModifyProblem)
	//分类
	adminGroup.GET("/category-list", service.GetCategoryList)
	adminGroup.POST("/category-create", service.CreateCategory)
	adminGroup.PUT("/category-modify", service.ModifyCategory)
	adminGroup.DELETE("/category-delete", service.DeleteCategory)

	//用户私有方法
	userGroup := r.Group("/user", middlewares.AuthUserCheck())
	userGroup.POST("/problem-submit", service.SubmitProblem())

	return r
}
