package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhubin/common"
)

func HelloGin(c *gin.Context){
	c.String(http.StatusOK, "Hello World")
}

// 绑定JSON的例子 ({"user": "manu", "password": "123"})
func AA(c *gin.Context) {
	var json common.Login
	if c.BindJSON(&json) == nil {
		if json.User == "manu" && json.Password == "123" {
			c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	}
}

// 绑定普通表单的例子 (user=manu&password=123)
func BB (c *gin.Context) {
	var form common.Login
	// 根据请求头中 content-type 自动推断.
	if c.Bind(&form) == nil {
		if form.User == "manu" && form.Password == "123" {
			c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	}
}

// 绑定多媒体表单的例子 (user=manu&password=123)
func CC (c *gin.Context) {
	var form common.LoginForm
	// 你可以显式声明来绑定多媒体表单：
	// c.BindWith(&form, binding.Form)
	// 或者使用自动推断:
	if c.Bind(&form) == nil {
		if form.User == "manu" && form.Password == "123" {
			c.JSON(200, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}
