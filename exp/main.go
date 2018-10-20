package main

import (
	// _ "github.com/lib/pq"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "lenslocked_dev"
)

type User struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)

	//	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}
	err = us.Create(&user)
	if err != nil {
		panic(err)
	}
	// Verify that the user has a Remember and RememberHash
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}

	// Now verify that we can lookup a user with that remember token
	user2, err := us.ByRemember(user.Remember)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *user2)

	// db.LogMode(true)

	// db.AutoMigrate(&User{})

	// name, email := getInfo()
	// u := &User{
	// 	Name:  name,
	// 	Email: email,
	// }
	// if err = db.Create(u).Error; err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", u)

	// var u User
	// u.Email = "chris@tsubaiso.jp"

	// db.Preload("Orders").First(&u)
	// // db.Preload("Orders").First(&u, id)

	// if db.Error != nil {
	// 	panic(db.Error)
	// }
	// fmt.Println(u)

	// createOrder(db, u, 1001, "Fake Description #1")
	// createOrder(db, u, 9999, "Fake Description #2")
	// createOrder(db, u, 8800, "Fake Description #3")

	// fmt.Println("Email:", u.Email)
	// fmt.Println("Number of orders:", len(u.Orders))
	// fmt.Println("Orders:", u.Orders)

	// fmt.Println(rand.String(10))
	// fmt.Println(rand.RememberToken())

	// hmac := hash.NewHMAC("my-secret-key")
	// fmt.Println(hmac.Hash("this is my string to hash"))
	// db.Where("id <= ?", maxId).First(&u)
	// if db.Error != nil {
	// 	panic(db.Error)
	// }
	// fmt.Println(u)

	// db.Where(u).First(&u)
	// if db.Error != nil {
	// 	panic(db.Error)
	// }
	// fmt.Println(u)

	// var users []User
	// db.Find(&users)
	// if db.Error != nil {
	// 	panic(db.Error)
	// }
	// fmt.Println("Retrieved", len(users), "users.")
	// fmt.Println(users)
	// 	err = db.Ping()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Successfully connected to database.")

	// 	var id int
	// 	var name, email string

	// 	row := db.QueryRow(`
	// SELECT id, name, email
	// FROM users
	// WHERE id=$1`, 1)
	// 	err = row.Scan(&id, &name, &email)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("ID:", id, "Name:", name, "Email:", email)

	// 	fmt.Println("-------------------")

	// 	rows, err := db.Query(`
	// SELECT id, name, email FROM users WHERE email = $1 OR ID > $2`, "jon@calhoun.io", 4)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for rows.Next() {
	// 		rows.Scan(&id, &name, &email)
	// 		fmt.Println("ID:", id, "Name:", name, "Email:", email)
	// 	}
	// 	fmt.Println("-------------------")

	// 	for i := 1; i < 6; i++ {
	// 		userId := 1
	// 		if i > 3 {
	// 			userId = 2
	// 		}
	// 		amount := 1000 * i
	// 		description := fmt.Sprintf("USB Adapter x%d", i)

	// 		err = db.QueryRow(`
	// INSERT into orders (user_id, amount, description) VALUES ($1, $2, $3) RETURNING id`,
	// 			userId, amount, description).Scan(&id)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Println("Created an order with the ID:", id)
	// 	}
	// 	db.Close()

}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	return name, email
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}
