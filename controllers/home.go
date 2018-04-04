// 主页 --- 显示所有文章及分类
package controllers

import (
	"beeblog/models"

	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	// 设置主页高亮
	c.Data["IsHome"] = true
	// 设置模版
	c.TplName = "home.html"
	// 检查是否登录
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	// 获取所有文章
	topics, err := models.GetAllTopics(
		c.Input().Get("cate"), c.Input().Get("lable"), true)
	if err != nil {
		beego.Error(err)
	}
	// 设置文章属性及属性值
	c.Data["Topics"] = topics
	// 获取所有分类
	categories, err := models.GetAllCategories()

	// 输出用于测试
	//c.Data["categories"] = categories
	//c.Ctx.WriteString(categories)

	if err != nil {
		beego.Error(err)
	}
	// 返回所有分类
	c.Data["Categories"] = categories
}
