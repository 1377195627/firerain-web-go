package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/firerainos/firerain-web-go/userCenter"
)

type User struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	user := User{}

	if err := ctx.Bind(&user); err != nil {
		ctx.JSON(400, gin.H{
			"code":    105,
			"message": err.Error(),
		})
		return
	}

	session := sessions.Default(ctx)

	u, err := userCenter.GetUserByName(user.Username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    105,
			"message": err.Error(),
		})
	}

	password := userCenter.EncryptionPassword(u.Username,user.Password,u.Email)

	if u.Password == password {
		session.Set("username", user.Username)
		session.Save()

		ctx.JSON(200, gin.H{
			"code": "0",
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    100,
			"message": "username or password errors",
		})
	}
}
