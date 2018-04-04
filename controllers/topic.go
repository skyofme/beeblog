// 文章相关功能 --- 添加 删除 修改 查看
package controllers

import (
	"path"
	"strings"

	"github.com/astaxie/beego"

	"beeblog/models"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsTopic"] = true
	c.TplName = "topic.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	topics, err := models.GetAllTopics("", "", false)
	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"] = topics
}

func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	// 解析表单
	tid := c.Input().Get("tid")
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	category := c.Input().Get("category")
	lable := c.Input().Get("lable")

	// 获取附件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}

	var attachment string
	if fh != nil {
		// 保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		err = models.AddTopic(title, category, lable, content, attachment)
	} else {
		err = models.ModifyTopic(tid, title, category, lable, content, attachment)
	}

	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic", 302)
}

func (c *TopicController) Add() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "topic_add.html"
	c.Data["IsLogin"] = true
}

func (c *TopicController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(c.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic", 302)
}

func (c *TopicController) Modify() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "topic_modify.html"

	tid := c.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
	c.Data["IsLogin"] = true
}

func (c *TopicController) View() {
	c.TplName = "topic_view.html"

	reqUrl := c.Ctx.Request.RequestURI
	i := strings.LastIndex(reqUrl, "/")
	tid := reqUrl[i+1:]
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Lables"] = strings.Split(topic.Lables, " ")

	replies, err := models.GetAllReplies(tid)
	if err != nil {
		beego.Error(err)
		return
	}

	c.Data["Replies"] = replies
	c.Data["IsLogin"] = checkAccount(c.Ctx)
}
