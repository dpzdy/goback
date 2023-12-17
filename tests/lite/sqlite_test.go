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
func TestEval(t *testing.T) {
	weight := 2.5 // 以KG为单位的实际重量
	cost := calCost(weight)
	fmt.Printf("实际重量为 %.2f KG 的快递费用为 %d 元\n", weight, cost)
}

// ********************************************************************
// Q2生成对应的测试代码和测试订单数据
// 1使用sqlite数据库存储订单数据
type Order struct {
	gorm.Model
	ID        int
	UID       int
	Weight    float64
	CreatedAt time.Time
}

// ********************************************************************
// 2生成1000个用户id
type User struct {
	gorm.Model
	ID  int
	UID int
}

func batchInsertUser(users []*User) error {
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
	err := batchInsertUser(ids)
	if err != nil {
		fmt.Println(err)
	}

}

// ********************************************************************
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
func TestGenUID(t *testing.T) {
	genUID()
}

// 3)weight随机生成，但要保证在所有订单中，计费重量的分布权重大致为1/W，例如2KG订单的数量与8KG单的数量比约为(1/2):(1/8)=4:1
// 4)请用注释描述清楚所使用的权重算法原理
/*
生成权重方法
1.根据比例计算整体大小，即计算res = 1+1/2+1/3+...+1/100
2.分别计算1/x在整体中的比例，即(1/x)/res,再乘订单总数numOrders，即为对应登记下的订单个数
3.进行两层循环，外层是每个等级的个数，内层是分别生成对应个数的随机数
*/
func genWeight(numOrders int) []float64 {
	res := 0.0
	//1.根据比例计算整体大小，即计算1+1/2+1/3+...+1/100
	var ratioArr []float64
	for i := 1; i <= 100; i++ {
		cur := 1 / float64(i)
		res += cur
		ratioArr = append(ratioArr, cur)
	}
	ordersNum := float64(numOrders)
	all := 0.0
	//2.分别计算1/x在整体中的比例，即(1/x)/res,再乘订单总数numOrders，即为对应登记下的订单个数
	var allArr []int
	for i := 1; i <= 100; i++ {
		cur := (1.0 / float64(i)) * ordersNum / res
		curInt := math.Round(cur)
		all += curInt
		allArr = append(allArr, int(curInt))

	}
	//3.进行两层循环，外层是每个等级的个数，内层是分别生成对应个数的随机数
	var weights []float64
	for i, cnts := range allArr {
		for j := 0; j < cnts; j++ {
			randomFloat := rand.Float64() + float64(i)
			weightS := fmt.Sprintf("%.2f", randomFloat)
			weightF, err := strconv.ParseFloat(weightS, 3)
			if err != nil {
				fmt.Println(err)
			}
			weights = append(weights, weightF)
		}
	}
	return weights
}
func batchInsertOrder(orders []*Order) error {
	tx := NewDb().Begin()
	defer tx.Rollback()
	for _, order := range orders {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
	}
	return tx.Commit().Error
}
func genOrders(numOrders int) {
	genUIDMap()
	weights := genWeight(numOrders)
	var orders []*Order
	for i := 0; i < len(weights); i++ {
		userID := randomUID()
		createTime := time.Now()
		order := &Order{
			UID:       userID,
			Weight:    weights[i],
			CreatedAt: createTime,
		}
		orders = append(orders, order)
	}
	err := batchInsertOrder(orders)
	if err != nil {
		fmt.Println(err)
	}
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
	fmt.Println(len(orders))
	return total(orders)
}
func TestAllCost(t *testing.T) {
	UID := 230002
	fmt.Println(allCost(UID))
}

// ********************************************************************
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
