// 附件上传功能
package controllers

import (
	"io"
	"net/url"
	"os"

	"github.com/astaxie/beego"
)

type AttachController struct {
	beego.Controller
}

func (c *AttachController) Get() {
	// 获取文件路径
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:])
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	// 读取文件
	f, err := os.Open(filePath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	// 函数结束最后关闭文件流
	defer f.Close()
	// 复制文件
	_, err = io.Copy(c.Ctx.ResponseWriter, f)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
}
