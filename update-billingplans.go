package main

import (
	"fmt"
	"github.com/chargebee/chargebee-go"
	planAction "github.com/chargebee/chargebee-go/actions/plan"
	"github.com/chargebee/chargebee-go/models/plan"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// {u'6c1329e0-e983-ea40-8d2d-98f9ee035d3c': {'price': 14000, 'name': u'c1.large'}
type BillingPlan struct {
	ID        uint `gorm:"primary_key"`
	BillingID string
	Price     int32
	Name      string
}

var dbName = os.Getenv("DB_NAME")
var dbUser = os.Getenv("DB_USER")
var dbPass = os.Getenv("DB_PASS")
var dbHost = os.Getenv("DB_HOST")
var dbPort = os.Getenv("DB_PORT")
var dbDriver = os.Getenv("DB_DRIVER")

var chargebeeKey = os.Getenv("CHARGEBEE_KEY")
var chargebeeSite = os.Getenv("CHARGEBEE_SITE")

var DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbUser, dbPass, dbHost, dbPort, dbName)

var db, err = gorm.Open("mysql", DBURL)

func main() {

	fmt.Println("Configuring chargebee.. ")
	// chargebee.Configure("key", "site")
	chargebee.Configure(chargebeeKey, chargebeeSite)

	// db.LogMode(true)

	defer db.Close()

	fmt.Println("Performing AutoMigrate().. ")
	db.AutoMigrate(&BillingPlan{})

	fmt.Println("running planAction.List().. ")
	res, err := planAction.List(&plan.ListRequestParams{Limit: chargebee.Int32(100)}).ListRequest()
	if err != nil {
		panicOnError(err)
	} else {
		for i := range res.List {
			Plan := res.List[i].Plan
			billingPlan := BillingPlan{}
			res := db.Where("billing_id = ?", Plan.Id).First(&billingPlan)
			if res.RecordNotFound() {
				newBillingPlan := BillingPlan{BillingID: Plan.Id, Price: Plan.Price, Name: Plan.Name}
				fmt.Println("Creating", newBillingPlan)
				db.Create(&newBillingPlan)
			} else if db.Error != nil {
				panic("error:" + res.Error.Error())
			} else {
				fmt.Println("Billing Plan Exists: ", billingPlan)
			}
		}
	}

}
