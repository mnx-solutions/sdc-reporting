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
	Usage     float64
	Timestamp time.Time
	Processed bool
}

type InvoiceLineItem struct {
	UUID         string
	Alias        string
	UsageCharge  float64
	FromTime     time.Time
	UntilTime    time.Time
	UsageMinutes float64
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
		v1.GET("/usage/owner/:ownerid", fetchUsageOwner)
		v1.GET("/usage/owner/:ownerid/machine/:machineid", fetchUsageMachine)
		v1.GET("/invoice/owner/:ownerid", fetchInvoiceLineItems)
		v1.POST("/usage/owner/:ownerid/process", processUsageOwner)
	}
	router.Run()
	defer db.Close()

}

// fetchAllTodo fetch all todos
func fetchUsageOwner(c *gin.Context) {
	ownerUUID := c.Param("ownerid")

	var total_usage float64
	var watcherData []RawWatcherData

	db.Where("owner_uuid = ? and processed = 0", ownerUUID).Find(&watcherData)

	if len(watcherData) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No data found!!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range watcherData {
		total_usage = item.Usage + total_usage
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "usage": total_usage})
}

func fetchUsageMachine(c *gin.Context) {
	ownerUUID := c.Param("owernid")
	vmUUID := c.Param("machineid")

	var total_usage float64
	var watcherData []RawWatcherData

	db.Where("owner_uuid = ? and uuid = ? and processed = 0", ownerUUID, vmUUID).Find(&watcherData)

	if len(watcherData) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No data found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range watcherData {
		total_usage = item.Usage + total_usage
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "usage": total_usage})
}

func fetchInvoiceLineItems(c *gin.Context) {
	ownerUUID := c.Param("ownerid")

	var total_usage float64
	var watcherData []RawWatcherData
	var usageCharge float64
	var fromTime time.Time
	var untilTime time.Time
	var invoiceLineItems []InvoiceLineItem

	lineItems := make(map[string]InvoiceLineItem)

	db.Where("owner_uuid = ? and processed = 0", ownerUUID).Find(&watcherData)

	if len(watcherData) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No data found!!!"})
		return
	}

	for _, item := range watcherData {
		if _, err := lineItems[item.UUID]; err {

			_tmpLineItem := lineItems[item.UUID]
			_newUsageCharge := float64(_tmpLineItem.UsageCharge) + usageCharge
                        fromTime = _tmpLineItem.FromTime
                        untilTime = _tmpLineItem.UntilTime

			if item.Timestamp.Before(_tmpLineItem.FromTime) {
				fromTime = item.Timestamp
 
			}
                        

			if item.Timestamp.After(_tmpLineItem.UntilTime) {
				untilTime = item.Timestamp
			}

			lineItems[item.UUID] = InvoiceLineItem{UUID: item.UUID, Alias: item.Alias, UsageCharge: _newUsageCharge, FromTime: fromTime, UntilTime: untilTime}

		} else {
			lineItems[item.UUID] = InvoiceLineItem{UUID: item.UUID, Alias: item.Alias, UsageCharge: item.Usage, FromTime: item.Timestamp, UntilTime: item.Timestamp}
		}
		total_usage = item.Usage + total_usage
	}

	for _, item := range lineItems {
		invoiceLineItems = append(invoiceLineItems, item)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "line_tems": invoiceLineItems})
}

func processUsageOwner(c *gin.Context) {
	ownerUUID := c.Param("owernid")
	var watcherData RawWatcherData
	db.Model(&watcherData).Where("owner_uuid = ? and processed = 0", ownerUUID).Update("processed", 1)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "processed": true})
}
