package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"testing"
	"time"
)

// Q1快递费用的计算函数
func calCost(weight float64) int {
	base := 18

	maxWeight := 100.0
	level := int(math.Ceil(weight))
	level = int(math.Min(float64(level), maxWeight))

	var total int
	for i := 1; i <= level; i++ {
		cost := base + (i-1)*5
		insurance := int(float64(total) * 0.01)
		if i == 1 {
			insurance = 0
		}
		total = int(math.Round(float64(cost + insurance)))
	}
	return total
}

// Q2生成对应的测试代码和测试订单数据
// 1使用sqlite数据库存储订单数据
type Order struct {
	gorm.Model
	ID        int
	UID       int
	Weight    float64
	CreatedAt time.Time
}

// 2生成1000个用户id
func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// 3生成100000条订单记录并插入数据库中
// 1)id不允许重复
// 2)uid从第2步生成的用户id中随机选择
func randomID() int {
	return 0
}

// 3)weight随机生成，但要保证在所有订单中，计费重量的分布权重大致为1/W，例如2KG订单的数量与8KG单的数量比约为(1/2):(1/8)=4:1
// 4)请用注释描述清楚所使用的权重算法原理
func randomWeight() float64 {
	return 0
}
func generateOrders(numOrders, numUsers int) []Order {

	var orders []Order

	for i := 0; i < numOrders; i++ {
		userID := randomID()

		// 生成随机重量，使得权重大致符合1/W的分布
		weight := randomWeight()
		createTime := time.Now()

		order := Order{
			UID:       userID,
			Weight:    weight,
			CreatedAt: createTime,
		}

		orders = append(orders, order)
	}

	return orders
}

// 4.程序提供查询功能，在命令行输入任意1个用户id，则显示此用户所有订单，并计算出此用户的总费用
func allTotal(id int) int {
	var orders []Order
	//根据id搜索菜单获取orders
	var sum int
	for _, order := range orders {
		curCost := calCost(order.Weight)
		sum += curCost
	}
	return sum
}

func TestEval(t *testing.T) {
	weight := 2.5 // 以KG为单位的实际重量
	cost := calCost(weight)
	fmt.Printf("实际重量为 %.2f KG 的快递费用为 %d 元\n", weight, cost)
}
func TestSqlite(t *testing.T) {
	//fmt.Println(11111111111111111111)
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// test connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected111!")
}
func TestOrm(t *testing.T) {
	// 建立数据库连接
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Order{})

	db.Create(&Order{UID: 122, Weight: 2.45, CreatedAt: time.Now()})

	var order Order
	db.First(&order, "UID = ?", 122)
	fmt.Println(order.Weight)
	//db.Delete(&order, 1)
}
