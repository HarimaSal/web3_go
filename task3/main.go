package main

import (
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	//err := initDB()
	//if err != nil {
	//	fmt.Printf("init DB failed, err:%v\n", err)
	//	return
	//}
	gormDB, err := initGormDB()
	if err != nil {
		fmt.Printf("init Gorm DB failed, err:%v\n", err)
		return
	}
	//gorm_task2(gormDB)
	gorm_task3(gormDB)
}

var db *sqlx.DB

/* sqlx初始化DB */
func initDB() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func initGormDB() (*gorm.DB, error) {
	// 数据库连接配置
	dsn := "root:root@tcp(localhost:3308)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	log.Println("数据库表创建成功！")
	return gormDb, err
}

/* SQL语句 */
// 题目1：CRUD
func SQL_task1() {
	/*  编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。*/
	//insert into students(name, age, grade) values('张三', 20, '三年级');
	/*  编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。*/
	//select * from students where age > 18;
	/*  编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。*/
	//update students set grade = '四年级' where name = '张三';
	/*  编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。*/
	//delete students where age < 15;

}

// 题目2：事务语句
func SQL_task2() {
	/*
		-- 开启事务
		start transaction;

		-- 从账户A向账户B转账100
		update accounts set balance = balance - 100 where id = '1'; -- A账户
		update accounts set balance = balance + 100 where id = '2'; -- B账户
		-- 检查账户1余额是否充足
		select balance into @balance from accounts where id = '1';
		if @balance < 0 then
		 -- 余额不足，回滚
		 rollback;
		else
		 -- 余额充足，提交，并记入transactions表
		 insert into transactions(from_account_id, to_account_id, amount) values('1','2', 100);
		 commit;
		end;

	*/
}

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

/* Sqlx入门 */
// 题目1：使用SQL扩展库进行查询
func sqlx_task1() {
	//-- 映射到切片
	sqlStr1 := "select id, name, department, salary from employees where department=?"
	var emp []Employee
	err := db.Select(&emp, sqlStr1, "技术部")
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	//-- 映射到结构体
	sqlStr2 := "select id, name, department, salary from employees where salary order by salary desc limit 1"
	var emp2 Employee
	err2 := db.Get(&emp2, sqlStr2)
	if err2 != nil {
		fmt.Printf("query failed, err:%v\n", err2)
		return
	}
}

// 题目2：实现类型安全映射
func sqlx_task2() {
	sqlStr := "select id, title, author, price from books where price > 50"
	var books []Book
	_ = db.Select(&books, sqlStr)
}

type Book struct {
	ID     int
	Title  string
	Author string
	Price  float64
}

/* 进阶gorm */
// 题目1：模型定义。编写Go代码，使用Gorm创建这些模型对应的数据库表
func gorm_task1() {
	// 已定义好结构体
}

// 题目2：关联查询
func gorm_task2(gormDb *gorm.DB) {
	//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	var user User
	gormDb.Where("username = ?", "小王").Preload("Posts.Comments").Find(&user)
	printUserWithPrettyJSON(user)
	fmt.Println("----")
	//编写Go代码，使用Gorm查询评论数量最多的文章信息
	var post Post
	err := gormDb.Raw("select * from post where id = (select post_id from comment c group by c.post_id order by count(*) desc limit 1)").
		Scan(&post).Error
	if err != nil {
		log.Printf("query failed, err:%v\n", err)
		return
	}
	printUserWithPrettyJSON(post)
}

// 题目3：钩子函数
func gorm_task3(gormDb *gorm.DB) {
	// 创建post
	user := User{
		Username: "小王",
		Email:    "wang@example.com",
		Password: "123456",
		Posts: []Post{
			{
				Title:   "我的第一篇博客",
				Content: "这是第一篇博客的内容",
			},
		},
	}
	gormDb.Create(&user)
	// 删除 comment
	var comment Comment
	gormDb.First(&comment, 3)
	gormDb.Delete(&comment)
	// 注意，gormDb.Delete(comment, 3)，这种通过主键直接删除的方式，GORM 不会先加载完整的 Comment对象
}

// User 用户模型（一对多关系：拥有多篇文章）
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Username string `gorm:"type:longtext"`
	Email    string `gorm:"type:longtext"`
	Password string `gorm:"type:longtext"`
	Posts    []Post `gorm:"foreignKey:UserID"` // 用户拥有的文章列表
}

func (User) TableName() string {
	return "user"
}

// Post 文章模型（一对多关系：属于一个用户，拥有多条评论）
type Post struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Title         string    `gorm:"type:longtext"`
	Content       string    `gorm:"type:longtext"`
	UserID        uint      `gorm:"type:int unsigned"` // 外键，关联User表的ID
	Counts        int64     `gorm:"type:int unsigned"`
	CommentStatus string    `gorm:"type:longtext"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
}

func (Post) TableName() string {
	return "post"
}

// Comment 评论模型（多对一关系：属于一篇文章）
type Comment struct {
	//gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Content string `gorm:"type:longtext"`     // 评论内容非空
	UserID  uint   `gorm:"type:int unsigned"` // 外键，关联User表的ID
	PostID  uint   `gorm:"type:int unsigned"` // 外键，关联Post表的ID
}

func (Comment) TableName() string {
	return "comment"
}

func printUserWithPrettyJSON[T any](data T) {
	s, err := prettyjson.Marshal(data)
	if err != nil {
		log.Printf("prettyjson error: %v", err)
		return
	}
	fmt.Println(string(s))
}

// Post 创建前增加评论数
func (p *Post) BeforeCreate(db *gorm.DB) (err error) {
	p.Counts++
	return
}

// Comment 在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"
func (c *Comment) BeforeDelete(db *gorm.DB) (err error) {
	var post Post
	postTx := db.Model(&Post{}).Where("id = ?", c.PostID)
	postTx.Find(&post)
	if post.Counts <= 1 {
		updates := map[string]interface{}{
			"CommentStatus": "无评论",
			"counts":        0,
		}
		postTx.Updates(updates)
	}
	return
}

func (c *Comment) AfterDelete(db *gorm.DB) (err error) {
	var post Post
	postTx := db.Model(&Post{}).Where("id = ?", c.PostID)
	postTx.Update("counts", post.Counts-1)
	return
}
