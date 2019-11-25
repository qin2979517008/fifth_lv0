package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
//账号和用户信息的映射关系
var account map[string]User = make(map[string]User)

func main() {
	router := gin.Default()
	//通过HTTP方法绑定路由规则和路由函数
	//curl http://localhost:8080/hello
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	// 冒号:加上一个参数名组成路由参数。
	// curl http://localhost:8080/user/hello
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	router.GET("/query/hello", func(c *gin.Context) {
		name := c.DefaultQuery("name", "guest")
		c.JSON(200, gin.H{
			"message": "success",
			"name" : name,
		})
	})
	// http的报文体传Body，
	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nickname", "guest")

		c.JSON(http.StatusOK, gin.H{
			"status":  gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"message": message,
			"nickname":    nick,
		})
	})

	//参数绑定
	//c.Bind()	: 自动推断content-type是x-www-form-urlencoded表单还是json的参数。
	router.POST("/login", func(c *gin.Context) {
		var user User

		err := c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		if v, ok := account[user.Username]; ok && v.Password == user.Password {
			c.JSON(http.StatusOK, gin.H{
				"username":   v.Username,
				"message": "你已经登陆成功",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message" : "账号或者密码有误",
			})
		}
	})

	router.POST("/register", func(c *gin.Context) {
		var user User
		err := c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		username := user.Username
		if _, ok := account[username]; ok {
			message := "用户名" + username + "已存在"
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": message,
			})
		} else {
			account[username] = user
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": "注册成功",
			})
		}
	})

	//管理组织分组api
	someGroup := router.Group("/hello")
	{
		someGroup.GET("/getting", getting)
		someGroup.POST("/posting", posting)
	}

	router.GET("/auth/signin", func(c *gin.Context) {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "123",
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		c.String(http.StatusOK, "Login successful")
	})

	router.Run(":8080")
}

func posting(c *gin.Context)  {
	username := c.DefaultPostForm("username", "guest")//可设置默认值
	msg := c.PostForm("msg")
	title := c.PostForm("title")
	fmt.Printf("username is %s, msg is %s, title is %s\n", username, msg, title)
}

func getting(c *gin.Context)  {
	name := c.DefaultQuery("name", "Guest") //可设置默认值
	// 是 c.Request.URL.Query().Get("lastname") 的简写
	lastname := c.Query("lastname")
	c.String(http.StatusOK, "Hello %s \n", name);
	fmt.Printf("Hello %s \n", name)
	fmt.Printf("Hello %s \n", lastname)
}

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password" bdinding:"required"`
}
