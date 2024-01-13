package db

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var client influxdb2.Client

type InfluxClient struct {
	url    string
	token  string
	org    string
	bucket string
}

func (influx *InfluxClient) initDB() influxdb2.Client {
	return influxdb2.NewClient(influx.url, influx.token)
}

func (influx *InfluxClient) getClient() influxdb2.Client {
	if client == nil {
		client = influx.initDB()
	}
	return client
}

func (influx *InfluxClient) WritePoint(point *write.Point) error {
	client := influx.getClient()

	writeApi := client.WriteAPIBlocking(influx.org, influx.bucket)
	fmt.Println("Writing point to influx db")
	return writeApi.WritePoint(context.Background(), point)
}

func (influx *InfluxClient) Close() {
	if client != nil {
		client.Close()
		client = nil
	}
}

func NewInfluxClient(url string, token string, org string, bucket string) *InfluxClient {
	return &InfluxClient{
		url:    url,
		token:  token,
		org:    org,
		bucket: bucket,
	}
}
