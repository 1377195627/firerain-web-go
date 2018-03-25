package main

import (
	"flag"
	"github.com/firerainos/firerain-web-go/api"
	"github.com/firerainos/firerain-web-go/core"
	"os"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/firerainos/firerain-web-go/userCenter"
	"strconv"
)

var port = flag.Int("p", 8080, "port")

func main() {
	flag.Parse()

	err := core.ParseConf("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("请配置config.json")
			os.Exit(0)
		}
		panic(err)
	}

	db, err := core.GetSqlConn()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&api.List{})
	db.AutoMigrate(&userCenter.User{})
	db.AutoMigrate(&userCenter.Group{})
	db.AutoMigrate(&userCenter.Permission{})
	db.Close()

	router := gin.Default()

	store := sessions.NewCookieStore([]byte("firerain"))
	router.Use(sessions.Sessions("firerain-session", store))

	apiRouter := router.Group("/api")

	apiRouter.GET("/installer/package/:de", api.GetPackage)

	apiRouter.POST("/login", api.Login)

	apiRouter.POST("/list/add", api.AddList)
	apiRouter.GET("/list/list", checkLoginMiddleware, api.GetList)
	apiRouter.DELETE("/list/delete", checkLoginMiddleware, api.DelList)
	apiRouter.GET("/list/pass", checkLoginMiddleware, api.PassList)

	packageRouter := apiRouter.Group("/package",checkLoginMiddleware)

	packageRouter.POST("/add")
	packageRouter.DELETE("/delete")
	packageRouter.GET("/list")

	itemRouter := apiRouter.Group("/item",checkLoginMiddleware)

	itemRouter.POST("/add")
	itemRouter.DELETE("/delete")
	itemRouter.GET("/list")

	uCenterRouter := apiRouter.Group("/userCenter")

	uCenterRouter.POST("/user/add",api.AddUser)
	uCenterRouter.DELETE("/user/delete",api.DeleteUser)
	uCenterRouter.GET("/user/list",api.GetUser)

	uCenterRouter.POST("/group/add",api.AddGroup)
	uCenterRouter.DELETE("/group/delete",api.DeleteGroup)
	uCenterRouter.GET("/group/list",api.GetGroup)

	router.Run(":" + strconv.Itoa(*port))
}

func checkLoginMiddleware(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("username") != core.Conf.Username {
		ctx.JSON(200, gin.H{
			"code":    101,
			"message": "unauthorized",
		})

		ctx.Abort()
	}

	ctx.Next()
}
