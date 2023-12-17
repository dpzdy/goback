package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func NewDb() *gorm.DB {
	// 建立数据库连接
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&User{})
	return db
}

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
type User struct {
	gorm.Model
	ID  int
	UID int
}

func batchInsert(users []*User) error {
	tx := NewDb().Begin()
	defer tx.Rollback()

	for _, user := range users {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
func genUID() {
	ids := make([]*User, 0)
	maxUser := 1000
	pre := "23"
	for i := 1; i <= maxUser; i++ {
		index := fmt.Sprintf("%04d", i)
		curUID := pre + index
		intUID, err := strconv.Atoi(curUID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ids = append(ids, &User{UID: intUID})
	}
	batchInsert(ids)

}

// 3生成100000条订单记录并插入数据库中
// 1)id不允许重复
// 2)uid从第2步生成的用户id中随机选择
func randomInt() int {
	return 1 + rand.Intn(1000)
}

var idMap map[int]int

func genUIDMap() {
	db := NewDb()
	idMap = make(map[int]int, 0)
	var users []User
	db.Find(&users)
	for _, item := range users {
		idMap[item.ID] = item.UID
	}
}
func randomUID() int {
	ranID := randomInt()
	return idMap[ranID]
}

// 3)weight随机生成，但要保证在所有订单中，计费重量的分布权重大致为1/W，例如2KG订单的数量与8KG单的数量比约为(1/2):(1/8)=4:1
// 4)请用注释描述清楚所使用的权重算法原理
func lcm(end int) int {
	result := 1
	for n := 2; n <= end; n++ {
		result = lcmTwoNumber(result, n)
	}
	return result
}

func lcmTwoNumber(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func genWeight(numOrders int) []float64 {
	return nil
}
func genOrders(numOrders int) []Order {
	genUIDMap()
	db := NewDb()
	var orders []Order
	weights := genWeight(numOrders)
	for i := 0; i < numOrders; i++ {
		userID := randomUID()

		// 生成随机重量，使得权重大致符合1/W的分布
		createTime := time.Now()

		order := Order{
			UID:       userID,
			Weight:    weights[i],
			CreatedAt: createTime,
		}

		db.Create(&order)
	}
	return orders
}
func TestGenOrders(t *testing.T) {
	numOrders := 100000
	genOrders(numOrders)
}

// ********************************************************************
// 4.程序提供查询功能，在命令行输入任意1个用户id，则显示此用户所有订单，并计算出此用户的总费用
func total(orders []Order) int {
	//根据id搜索菜单获取orders
	var sum int
	for _, order := range orders {
		curCost := calCost(order.Weight)
		sum += curCost
	}
	return sum
}
func allCost(UID int) int {
	db := NewDb()
	var orders []Order
	db.Find(&orders, "UID = ?", UID)
	return total(orders)
}
func TestAllCost(t *testing.T) {
	UID := 230002
	fmt.Println(allCost(UID))
}

// ********************************************************************
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

	db := NewDb()
	//db.Create(&Order{UID: 122, Weight: 2.45, CreatedAt: time.Now()})

	var order Order
	db.First(&order, "UID = ?", 122)
	fmt.Println(order.Weight)
	//db.Delete(&order, 1)
}

func TestGenUID(t *testing.T) {
	//lcmMax := lcm(100)
	res := 0.0
	var ratioArr []float64
	for i := 1; i <= 100; i++ {
		cur := 1 / float64(i)
		res += cur
		ratioArr = append(ratioArr, cur)
	}
	fmt.Println(ratioArr)

	ordersNum := 100000.0
	all := 0.0
	var allArr []float64
	for i := 1; i < 100; i++ {
		cur := (1.0 / float64(i)) * ordersNum / res
		curInt := math.Round(cur)
		all += curInt
		allArr = append(allArr, curInt)

	}
	fmt.Println(all)
	fmt.Println(allArr[0] / allArr[49])

}
