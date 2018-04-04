// 评论功能 --- 增删
package controllers

import (
	"beeblog/models"

	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

// 添加评论
func (c *ReplyController) Add() {
	// 获取页面中返回的name属性tid对应的value
	tid := c.Input().Get("tid")
	err := models.AddReply(tid,
		c.Input().Get("nickname"), c.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}

	// 重定向到浏览的文章当中
	c.Redirect("/topic/view/"+tid, 302)
}

// 删除评论
func (c *ReplyController) Delete() {
	if !checkAccount(c.Ctx) {
		return
	}
	tid := c.Input().Get("tid")
	err := models.DeleteReply(c.Input().Get("rid"))
	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic/view/"+tid, 302)
}
