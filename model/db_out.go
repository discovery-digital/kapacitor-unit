package model

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gpestana/kapacitor-unit/io"
	"time"
)

// DbOUtTest
type DbOutTest struct {
	Name       string
	TaskName   string
	SeedData   []string
	SeedDb     string
	SeedRp     string
	SeedType   string
	AssertDb   string
	AssertRp   string
	AssertData []string
	TestType   string
	Task       Task
}

func (d *DbOutTest) Init(conf TestConf) error {
	d.Name = conf.Name
	d.TaskName = conf.TaskName
	d.SeedData = conf.SeedData
	d.SeedDb = conf.SeedDb
	d.SeedRp = conf.SeedRp
	d.AssertDb = conf.AssertDb
	d.AssertRp = conf.AssertRp
	d.AssertData = conf.AssertData
	d.TestType = conf.Type
	return nil
}

func (d DbOutTest) setup(k io.Kapacitor, i io.Influxdb) error {

	if err := i.Setup(d.SeedDb, d.SeedRp); err != nil {
		return err
	}

	if err := i.Setup(d.AssertDb, d.AssertRp); err != nil {
		return err
	}

	dbrps := DbRps{{Db: d.SeedDb, Rp: d.SeedRp}, {Db: d.AssertDb, Rp: d.AssertRp}}

	if err := LoadTask(d.Task, dbrps, k); err != nil {
		return err
	}

	return nil
}

// Adds model data
func (t *DbOutTest) addData(k io.Kapacitor, i io.Influxdb) error {
	switch t.TestType {
	case "stream":
		// adds data to kapacitor
		err := k.Data(t.SeedData, t.SeedDb, t.SeedRp)
		if err != nil {
			return err
		}
	case "batch":
		// adds data to InfluxDb
		err := i.Data(t.SeedData, t.SeedDb, t.SeedRp)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deletes data, database and retention policies created to run the model
func (t *DbOutTest) teardown(k io.Kapacitor, i io.Influxdb) {
	glog.Info("DEBUG:: teardown model: ", t.Name)
	switch t.TestType {
	case "batch":
		err := i.CleanUp(t.SeedDb)
		if err != nil {
			glog.Error("Error performing teardown in cleanup. error: ", err)
		}
	}

	if err := i.CleanUp(t.AssertDb); err != nil {
		glog.Error("Error performing teardown in cleanup. error: ", err)
	}

	if err := k.Delete(t.TaskName); err != nil {
		glog.Error("Error performing teardown in delete error: ", err)
	}
}

func (t DbOutTest) Run(kapacitor io.Kapacitor, influxdb io.Influxdb) error {

	defer t.teardown(kapacitor, influxdb)
	if err := t.setup(kapacitor, influxdb); err != nil {
		return err
	}

	if err := t.addData(kapacitor, influxdb); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	data, err := influxdb.Query(t.AssertDb, "SELECT * FROM \"sum\"")
	_ = data
	fmt.Printf("DB: ", data)
	return err
}

func (DbOutTest) Validate() error {
	return nil
}

func (d DbOutTest) GetName() string {
	return d.Name
}

func (d DbOutTest) GetType() string {
	return d.TestType
}

func (d *DbOutTest) SetTask(task Task) {
	d.Task = task
}

func (DbOutTest) Passed() bool {
	return true
}
