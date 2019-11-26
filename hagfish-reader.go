package main
import (
    "fmt"
    "os"
    "log"
    "bufio"
    "encoding/json"
    "time"
//    "io/ioutil"
)

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

func main() {

    filename := os.Args[1]
    var watcherData WatcherData

    jsonFile, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer jsonFile.Close()
    s := bufio.NewScanner(jsonFile)
    for s.Scan() {
        if err := json.Unmarshal(s.Bytes(), &watcherData); err != nil {
           fmt.Println(err)
        }
        usageType := watcherData.Type
        if usageType == "usage" {
            ownerUUID := watcherData.Config.Attributes.OwnerUUID
            billingID := watcherData.Config.Attributes.BillingID
            vmUUID := watcherData.UUID
            alias := watcherData.Config.Attributes.Alias
            timestamp := watcherData.Timestamp
            net0SentBytes := watcherData.NetworkUsage.Net0.SentBytes
            net0ReceivedBytes := watcherData.NetworkUsage.Net0.ReceivedBytes
            // net1SentBytes := watcherData.NetworkUsage.Net1.SentBytes
            // net1ReceivedBytes := watcherData.NetworkUsage.Net1.ReceivedBytes

            fmt.Println(ownerUUID, vmUUID, billingID, alias, net0SentBytes, net0ReceivedBytes,  timestamp)
        }
     }
    if s.Err() != nil {
        fmt.Println("e")
    }
}

