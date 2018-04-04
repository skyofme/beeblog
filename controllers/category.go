// 分类管理功能 -- 增删
package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	// 检查是否有操作，从页面接受name为op的value
	op := c.Input().Get("op")
	switch op {
	// 进行添加操作
	case "add":
		// 获取分类名称
		name := c.Input().Get("name")
		if len(name) == 0 {
			break
		}
		// 添加分类
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		// 添加失败则显示302错误
		c.Redirect("/category", 302)
		return
	// 进行删除操作
	case "del":
		// 获取目录id
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		// 删除分类
		err := models.DeleteCategory(id)
		if err != nil {
			beego.Error(err)
		}
		// 删除失败则显示302错误
		c.Redirect("/category", 302)
		return
	}
	// 目录高亮
	c.Data["IsCategory"] = true
	// 跳转模板
	c.TplName = "category.html"
	// 检查是否登录
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var err error
	// 获取所有分类，并设置属性及属性值
	c.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
}
