package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Student 结构体表示一个学生记录
type Student struct {
	//gorm.Model
	Id    int `gorm:"unique"`
	Name  string
	Age   int
	Grade string
}

// 显示菜单选项
func showMenu() {
	fmt.Println("****************************")
	fmt.Println("********* 1. 添加学生信息 *********")
	fmt.Println("********* 2. 显示学生信息 *********")
	fmt.Println("********* 3. 修改学生信息 *********")
	fmt.Println("********* 4. 删除学生信息 *********")
	fmt.Println("********* 5. 获取单个学生信息 *********")
	fmt.Println("********* 6. 退出 *********")
	fmt.Println("****************************")
}

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Silent,
		},
	)
	// 打开或创建数据库
	db, err := gorm.Open(sqlite.Open("./students.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}

	// 自动迁移模式
	db.AutoMigrate(&Student{})

	var choice int
	for {
		showMenu()
		fmt.Print("请输入选项: ")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			// 添加学生
			addStudent(db)

		case 2:
			// 获取所有学生并显示
			getStudents(db)

		case 3:
			// 更新学生信息
			updateStudent(db)

		case 4:
			// 删除学生
			deleteStudent(db)

		case 5:
			// 获取单个学生信息

			getStudent(db)

		case 6:
			// 退出
			fmt.Println("退出系统，谢谢使用！")
			return

		default:
			fmt.Println("无效选项，请重新输入。")
		}
	}
}

// 添加一个新的学生记录到数据库
func addStudent(db *gorm.DB) {
	var student Student

	for {
		fmt.Println("学号:")
		fmt.Scan(&student.Id)
		var existingStudent Student
		result := db.Take(&existingStudent, "id = ?", student.Id)
		if result.Error == nil {
			fmt.Println("添加学生信息失败: 学号已存在，请重新输入")
		} else if result.Error == gorm.ErrRecordNotFound {
			break
		} else {
			log.Fatal("检查学号时出错: ", result.Error)
		}
	}

	fmt.Println("姓名:")
	fmt.Scan(&student.Name)
	fmt.Println("年龄:")
	fmt.Scan(&student.Age)
	fmt.Println("成绩:")
	fmt.Scan(&student.Grade)
	result := db.Create(&student)
	if result.Error != nil {
		log.Fatal("添加学生失败", result.Error)
	}
	fmt.Println("学生信息添加成功")
}

// 获取所有学生记录
func getStudents(db *gorm.DB) []Student {

	var students []Student
	result := db.Find(&students)
	if result.Error != nil {
		log.Println("获取学生信息失败: ", result.Error)
	}

	fmt.Println("所有学生信息:")
	for _, student := range students {
		fmt.Printf("学号: %d, 姓名: %s, 年龄: %d, 成绩: %s\n", student.Id, student.Name, student.Age, student.Grade)
	}
	return students
}

// 根据ID获取单个学生记录
func getStudent(db *gorm.DB) Student {
	var student Student
	for {
		fmt.Println("请输入学号：")
		fmt.Scan(&student.Id)
		result := db.First(&student, student.Id)
		if result.Error != nil {
			log.Println("未找到该学生！")
		} else {
			break
		}
	}

	fmt.Printf("学号: %d, 姓名: %s, 年龄: %d, 成绩: %s\n", student.Id, student.Name, student.Age, student.Grade)
	return student
}

// 更新现有的学生记录
func updateStudent(db *gorm.DB) {
	var student Student
	for {
		fmt.Print("输入要更新的学生ID: ")
		fmt.Scan(&student.Id)
		result := db.First(&student, student.Id)
		if result.Error != nil {
			fmt.Println("未找到该学生，请重新输入！")
		} else {
			break
		}
	}

	for {
		var newId int
		fmt.Print("输入新学生ID: ")
		fmt.Scan(&newId)
		if newId != student.Id {
			var existingStudent Student
			result := db.First(&existingStudent, newId)
			if result.Error == nil {
				fmt.Println("更新学生信息失败: 新学号已存在，请重新输入新学号")
			} else {
				//删除原来的记录
				db.Delete(&student, student.Id)

				student.Id = newId
				break
			}
		} else {
			break
		}
	}

	fmt.Print("输入新学生姓名: ")
	fmt.Scan(&student.Name)
	fmt.Print("输入新学生年龄: ")
	fmt.Scan(&student.Age)
	fmt.Print("输入新学生成绩: ")
	fmt.Scan(&student.Grade)

	result := db.Save(&student)
	if result.Error != nil {
		fmt.Println("更新学生信息失败: ", result.Error)
	} else {
		fmt.Println("学生信息更新成功")
	}
}

// 根据ID删除学生记录
func deleteStudent(db *gorm.DB) {
	var student Student
	for {
		fmt.Print("输入学生ID: ")
		fmt.Scan(&student.Id)
		var existingStudent Student
		result := db.First(&existingStudent, student.Id)
		if result.Error != nil {
			fmt.Println("未找到该学生,请重新输入！")
		} else {
			db.First(student.Id)
			fmt.Println("学生信息删除成功")
			break
		}
	}
}
