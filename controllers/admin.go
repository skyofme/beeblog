// 管理员功能
package controllers

import (
	"github.com/astaxie/beego"
)

type AdminController struct {
	// 有一些默认的实现方法，实现了controller这个接口（因为重写了该接口的GET方法）
	beego.Controller // 嵌入字段
}

// get方式访问控制器，调用重写的get方法，将模板名传给接收者
func (c *AdminController) Get() {
	// 设置模版文件为login页面
	c.TplName = "login.html"
}
