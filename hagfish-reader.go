package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

type BillingPlan struct {
	ID        uint `gorm:"primary_key"`
	BillingID string
	Price     int32
	Name      string
}

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

type WatcherData struct {
	ServerUUID     string    `json:"server_uuid"`
	DatacenterName string    `json:"datacenter_name"`
	Timestamp      time.Time `json:"timestamp"`
	Type           string    `json:"type"`
	UUID           string    `json:"uuid"`
	Config         struct {
		Name       string `json:"name"`
		Attributes struct {
			CreateTimestamp time.Time `json:"create-timestamp"`
			DatasetUUID     string    `json:"dataset-uuid"`
			BillingID       string    `json:"billing-id"`
			OwnerUUID       string    `json:"owner-uuid"`
			Alias           string    `json:"alias"`
		} `json:"attributes"`
	} `json:"config"`
	NetworkUsage struct {
		Net0 struct {
			SentBytes     int64     `json:"sent_bytes"`
			ReceivedBytes int64     `json:"received_bytes"`
			CounterStart  time.Time `json:"counter_start"`
		} `json:"net0"`
		Net1 struct {
			SentBytes     int64     `json:"sent_bytes"`
			ReceivedBytes int64     `json:"received_bytes"`
			CounterStart  time.Time `json:"counter_start"`
		} `json:"net1"`
	} `json:"network_usage"`
}

var dbName = os.Getenv("DB_NAME")
var dbUser = os.Getenv("DB_USER")
var dbPass = os.Getenv("DB_PASS")
var dbHost = os.Getenv("DB_HOST")
var dbPort = os.Getenv("DB_PORT")
var dbDriver = os.Getenv("DB_DRIVER")
var DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbUser, dbPass, dbHost, dbPort, dbName)

var db, err = gorm.Open("mysql", DBURL)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// db.LogMode(true)

	panicOnError(err)
	defer db.Close()

	db.AutoMigrate(&RawWatcherData{})

	billingData := make(map[string]float64)

	rows, err := db.Model(&BillingPlan{}).Select("billing_id, price").Rows() // (*sql.Rows, error)
	defer rows.Close()

	for rows.Next() {
		var billingPlan BillingPlan
		db.ScanRows(rows, &billingPlan)
		billingData[billingPlan.BillingID] = float64(billingPlan.Price)
	}

	filename := os.Args[1]
	var watcherData WatcherData
	var perMinute float64

	gzLogFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	gz, err := gzip.NewReader(gzLogFile)
	if err != nil {
		log.Fatal(err)
	}

	defer gz.Close()
	defer gzLogFile.Close()

	s := bufio.NewScanner(gz)

	mapWatcherData := make(map[string]RawWatcherData)
	for s.Scan() {
		if err := json.Unmarshal(s.Bytes(), &watcherData); err != nil {
			// fmt.Println()
		}
		usageType := watcherData.Type
		if usageType == "usage" {
			ownerUUID := watcherData.Config.Attributes.OwnerUUID
			billingID := watcherData.Config.Attributes.BillingID
			vmUUID := watcherData.UUID
			alias := watcherData.Config.Attributes.Alias
			timestamp := watcherData.Timestamp.Truncate(time.Second)
			// net0SentBytes := watcherData.NetworkUsage.Net0.SentBytes
			// net0ReceivedBytes := watcherData.NetworkUsage.Net0.ReceivedBytes
			// net1SentBytes := watcherData.NetworkUsage.Net1.SentBytes
			// net1ReceivedBytes := watcherData.NetworkUsage.Net1.ReceivedBytes

			if val, err := billingData[billingID]; err {
				perMinute = val / 720 / 60 / 100
			} else {
				// set 0.0.  This allows for unknown billingID's
				perMinute = 0.0
			}

			if _, err := mapWatcherData[vmUUID]; err {
				_tmpRawWatcherData := mapWatcherData[vmUUID]
				_newUsage := float64(_tmpRawWatcherData.Usage) + perMinute
				mapWatcherData[vmUUID] = RawWatcherData{OwnerUUID: ownerUUID, UUID: vmUUID, BillingID: billingID, Usage: _newUsage, Timestamp: timestamp, Alias: alias}
			} else {
				mapWatcherData[vmUUID] = RawWatcherData{OwnerUUID: ownerUUID, UUID: vmUUID, BillingID: billingID, Usage: perMinute, Timestamp: timestamp, Alias: alias}
			}
		}

	}
	if s.Err() != nil {
		fmt.Println("e")
	}

	for _, v := range mapWatcherData {

		rawWatcherData := RawWatcherData{}
		res := db.Where("owner_uuid = ? and uuid = ? and billing_id = ? and timestamp = ?", v.OwnerUUID, v.UUID, v.BillingID, v.Timestamp).Find(&rawWatcherData)
		if res.RecordNotFound() {
			newRawWatcherData := RawWatcherData{OwnerUUID: v.OwnerUUID, UUID: v.UUID, BillingID: v.BillingID, Usage: v.Usage, Timestamp: v.Timestamp, Alias: v.Alias}
			db.Create(&newRawWatcherData)
		} else if db.Error != nil {
			panic("error:" + res.Error.Error())
		}
	}

}
