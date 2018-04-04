package models

import (
	. "fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com" // 通用函数包
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

// 设置全局变量：数据库名称、slite3数据库驱动
const (
	// 设置数据库路径
	_DB_NAME = "data/beeblog.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64     // 数据表字段id
	Title           string    // 分类名称
	Created         time.Time `orm:"index"` // 创建时间
	Views           int64     `orm:"index"` // 浏览次数
	TopicTime       time.Time `orm:"index"` // 发表时间
	TopicCount      int64     // 文章数量 数据表字段topic_count
	TopicLastUserId int64     // 最新操作时间
}

// 文章
type Topic struct {
	Id              int64     // id
	Uid             int64     // 用户id
	Title           string    // 文章标题
	Category        string    // 文章分类
	Lables          string    // 标签
	Content         string    `orm:"size(5000)"` // 文章内容
	Attachment      string    // 附件
	Created         time.Time `orm:"index"` // 创建时间
	Updated         time.Time `orm:"index"` // 更新时间
	Views           int64     `orm:"index"` // 浏览次数
	Author          string    // 作者
	ReplyTime       time.Time `orm:"index"` // 评论时间
	ReplyCount      int64     // 评论次数
	ReplyLastUserId int64     // 最新评论
}

// 评论
type Comment struct {
	Id      int64     // id
	Tid     int64     // 文章id
	Name    string    // 用户姓名
	Content string    `orm:"size(1000)"` // `orm:size(1000)`只供orm框架读取
	Created time.Time `orm:"index"`      // 评论时间
}

// 注册数据库
func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		// 递归创建所有目录（去除最后的文件名，权限默认）
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		// 创建数据库文件beeblog.db
		os.Create(_DB_NAME)
	}

	// 注册模型 --- 目录 文章 评论
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	//orm.RegisterDriver(_SQLITE3_DRIVER, orm.DR_Sqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

// 添加分类
func AddCategory(name string) error {
	// 创建orm对象
	o := orm.NewOrm()
	// 从category action获取赋值好的name值
	cate := &Category{Title: name}
	// 查询数据
	qs := o.QueryTable("category")
	// 分类已经存在则报错
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		Println("添加类型失败")
		return err
	}
	// 插入数据
	_, err = o.Insert(cate)
	if err != nil {
		Println("添加类型失败")
		return err
	}
	// 返回nil表示添加成功，否则添加失败
	return nil
}

// 删除分类
func DeleteCategory(id string) error {
	o := orm.NewOrm()
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	// 根据id删除分类
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

// 获取所有分类
func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)
	// 查询数据库
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

// 添加文章
func AddTopic(title, category, lable, content, attachment string) error {
	o := orm.NewOrm()
	// 处理标签
	lable = "$" + strings.Join(strings.Split(lable, " "), "#$") + "#"
	topic := &Topic{
		Title:      title,
		Category:   category,
		Lables:     lable,
		Content:    content,
		Attachment: attachment,
		Created:    time.Now(),
		Updated:    time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}

	// 更新分类统计
	// 创建一个category对象
	cate := &Category{Title: category}
	// cate := new(Category)
	// 查询数据库，返回该表的结果集
	qs := o.QueryTable("category")
	// 分类已经存在则报错
	err = qs.Filter("title", category).One(cate)
	if err == nil {

	} else {
		// 如果不存在我们就直接创建，并进行更新
		_, err := o.Insert(cate)
		if err != nil {
			Println("添加分类失败")
		}
	}
	cate.TopicCount++
	_, err = o.Update(cate)

	return err
}

// 获取文章
func GetTopic(tid string) (*Topic, error) {
	o := orm.NewOrm()
	// 获取文章id
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	// 获取失败 返回nil 并报错
	if err != nil {
		return nil, err
	}
	// 类似于空参构造
	topic := new(Topic)
	qs := o.QueryTable("topic")
	// 为topic对象 赋值
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	// 文章访问量+1
	topic.Views++
	_, err = o.Update(topic)

	topic.Lables = strings.Replace(strings.Replace(
		topic.Lables, "#", " ", -1), "$", "", -1)
	return topic, nil
}

// 修改文章
func ModifyTopic(tid, title, category, lable, content, attachment string) error {
	// 根据id获取文章
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	lable = "$" + strings.Join(strings.Split(lable, " "), "#$") + "#"

	var oldCate, oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category = category
		topic.Lables = lable
		topic.Content = content
		topic.Attachment = attachment
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}

	// 更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}

	// 删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return nil
}

// 删除文章
func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var oldCate string
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	return err
}

// 获取所有文章
func GetAllTopics(category, lable string, isDesc bool) (topics []*Topic, err error) {
	o := orm.NewOrm()

	topics = make([]*Topic, 0)

	qs := o.QueryTable("topic")
	// 文章排序
	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(lable) > 0 {
			qs = qs.Filter("lables__contains", "$"+lable+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)

	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

// 添加评论
func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
	}
	return err
}

// 获取所有评论
func GetAllReplies(tid string) (replies []*Comment, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Comment, 0)

	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}

// 删除评论
func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var tidNum int64
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_, err = o.Update(topic)
	}
	return err
}
