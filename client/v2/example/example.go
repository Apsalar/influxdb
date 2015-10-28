package client_example

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/influxdb/influxdb/client/v2"
)

func ExampleNewClient() client.Client {
	u, _ := url.Parse("http://localhost:8086")

	// NOTE: this assumes you've setup a user and have setup shell env variables,
	// namely INFLUX_USER/INFLUX_PWD. If not just omit Username/Password below.
	client := client.NewClient(client.Config{
		URL:      u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	})
	return client
}

func ExampleNewCustomClient() client.Client {
	u, _ := url.Parse("http://localhost:8086")

	client := client.NewClient(client.Config{
		URL: u,
		// Pass in a *http.Client
		HTTPClient: &http.Client{
			Timeout: time.Duration(5) * time.Second,
		},
	})
	return client
}

func ExampleWrite() {
	// Make client
	u, _ := url.Parse("http://localhost:8086")
	c := client.NewClient(client.Config{
		URL: u,
	})

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "BumbleBeeTuna",
		Precision: "s",
	})

	// Create a point and add to batch
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		panic(err.Error())
	}
	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)
}

// Write 1000 points
func ExampleWrite1000() {
	sampleSize := 1000

	// Make client
	u, _ := url.Parse("http://localhost:8086")
	clnt := client.NewClient(client.Config{
		URL: u,
	})

	rand.Seed(42)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "systemstats",
		Precision: "us",
	})

	for i := 0; i < sampleSize; i++ {
		regions := []string{"us-west1", "us-west2", "us-west3", "us-east1"}
		tags := map[string]string{
			"cpu":    "cpu-total",
			"host":   fmt.Sprintf("host%d", rand.Intn(1000)),
			"region": regions[rand.Intn(len(regions))],
		}

		idle := rand.Float64() * 100.0
		fields := map[string]interface{}{
			"idle": idle,
			"busy": 100.0 - idle,
		}

		pt, err := client.NewPoint(
			"cpu_usage",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			println("Error:", err.Error())
			continue
		}
		bp.AddPoint(pt)
	}

	err := clnt.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleQuery() {
	// Make client
	u, _ := url.Parse("http://localhost:8086")
	c := client.NewClient(client.Config{
		URL: u,
	})

	q := client.Query{
		Command:   "SELECT count(value) FROM shapes",
		Database:  "square_holes",
		Precision: "ns",
	}
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		log.Println(response.Results)
	}
}

func ExampleCreateDatabase() {
	// Make client
	u, _ := url.Parse("http://localhost:8086")
	c := client.NewClient(client.Config{
		URL: u,
	})

	q := client.Query{
		Command: "CREATE DATABASE telegraf",
	}
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		log.Println(response.Results)
	}
}
