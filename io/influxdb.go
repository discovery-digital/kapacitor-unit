package io

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Influxdb service configurations
type Influxdb struct {
	Host   string
	Client http.Client
}

func NewInfluxdb(host string) Influxdb {
	return Influxdb{
		host,
		http.Client{},
	}
}

// Adds model data to influxdb
func (influxdb Influxdb) Data(data []string, db string, rp string) error {
	url := influxdb.Host + influxdb_write + "db=" + db + "&rp=" + rp
	for _, d := range data {
		resp, err := influxdb.Client.Post(url, "application/x-www-form-urlencoded",
			bytes.NewBuffer([]byte(d)))

		if err != nil {
			return err
		}

		if resp.StatusCode >= 400 {
			return errors.New(resp.Status)
		}
		glog.Info("DEBUG:: Influxdb added ["+d+"] to "+url)
	}
	return nil
}

// Creates db and rp where tests will run
func (influxdb Influxdb) Setup(db string, rp string) error {
	glog.Info("DEGUB:: Influxdb setup ", db+":"+rp)
	// If no retention policy is defined, use "autogen"
	if rp == "" {
		rp = "autogen"
	}
	q := "q=CREATE DATABASE \""+db+"\" WITH DURATION 1h REPLICATION 1 NAME \""+rp+"\""
	baseUrl := influxdb.Host + "/query"
	_, err := influxdb.Client.Post(baseUrl, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(q)))
	if err != nil {
		return err
	}
	return nil
}

func(influxdb Influxdb)Query(db string, query string) (map[string] interface{}, error ){
	baseUrl := influxdb.Host + "/query?db=" + db + "&q=" + url.QueryEscape(query)

	resp, err := influxdb.Client.Get(baseUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	j := make(map[string]interface{})
	err = json.Unmarshal(body, &j)

	return j, err
}

func (influxdb Influxdb) CleanUp(db string) error {
	q := "q=DROP DATABASE \""+db+"\""
	baseUrl := influxdb.Host + "/query"
	_, err := influxdb.Client.Post(baseUrl, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(q)))
	if err != nil {
		return err
	}
	glog.Info("DEBUG:: Influxdb cleanup database ", q)
	return nil
}
