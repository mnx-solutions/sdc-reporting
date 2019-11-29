package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"time"
)

type RawWatcherData struct {
	ID        uint `gorm:"primary_key"`
	OwnerUUID string
	UUID      string
	BillingID string
	Alias     string
	Price     float64
	Timestamp time.Time
	Processed bool
}

var dbName = os.Getenv("DB_NAME")
var dbUser = os.Getenv("DB_USER")
var dbPass = os.Getenv("DB_PASS")
var dbHost = os.Getenv("DB_HOST")
var dbPort = os.Getenv("DB_PORT")
var dbDriver = os.Getenv("DB_DRIVER")
var DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbUser, dbPass, dbHost, dbPort, dbName)

var db, err = gorm.Open("mysql", DBURL)

func main() {

	if err != nil {
		fmt.Println(err)
		fmt.Println("err open databases")
		fmt.Println(DBURL)
		return
	}
	defer db.Close()

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		v1.GET("/balance/owner/:id", fetchBalanceOwner)
		v1.GET("/balance/owner/:id/vm/:vmid", fetchBalanceVm)
	}
	router.Run()
	defer db.Close()

}

// fetchAllTodo fetch all todos
func fetchBalanceOwner(c *gin.Context) {
	ownerUUID := c.Param("id")

	var totalCost float64
	var watcherData []RawWatcherData

	db.Where("owner_uuid = ? and processed = 0", ownerUUID).Find(&watcherData)

	if len(watcherData) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No data found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range watcherData {
		totalCost = item.Price + totalCost
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "total_cost": totalCost})
}

func fetchBalanceVm(c *gin.Context) {
	ownerUUID := c.Param("id")
        vmUUID := c.Param("vmid")

        var totalCost float64
        var watcherData []RawWatcherData

        db.Where("owner_uuid = ? and uuid = ? and processed = 0", ownerUUID, vmUUID).Find(&watcherData)

        if len(watcherData) <= 0 {
                c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No data found!"})
                return
        }

        //transforms the todos for building a good response
        for _, item := range watcherData {
                totalCost = item.Price + totalCost
        }
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "total_cost": totalCost})
}

