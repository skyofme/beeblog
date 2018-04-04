// 管理员登录功能 --- 登录 校验 cookie
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	// 判断是否为退出操作
	if c.Input().Get("exit") == "true" {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.Ctx.SetCookie("pwd", "", -1, "/")
		c.Redirect("/", 302)
		return
	}
	// 不是则跳转到登录界面
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	// 获取表单信息
	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autoLogin") == "on"

	// 验证用户名及密码
	if uname == beego.AppConfig.String("adminName") &&
		pwd == beego.AppConfig.String("adminPass") {
		maxAge := 0
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		// 设置cookie 保存用户及密码
		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", pwd, maxAge, "/")
	}

	c.Redirect("/", 302)
	return
}

// 密码校验
func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}

	uname := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}

	pwd := ck.Value
	return uname == beego.AppConfig.String("adminName") &&
		pwd == beego.AppConfig.String("adminPass")
}
