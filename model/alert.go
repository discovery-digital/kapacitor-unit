// Responsible for setting up, run, gather results and tear down a model. It
// exposes the method model.Run(), which saves the model results in the AlertTest
// struct or fails.
package model


import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gpestana/kapacitor-unit/io"
	"regexp"
	"time"
)



// AlertTest is a model that asserts that the expected alerts are triggered
type AlertTest struct {
	Name     string
	TaskName string `yaml:"task_name,omitempty"`
	Data     []string
	RecId    string `yaml:"recording_id"`
	Expects  Result
	Result   Result
	Db       string
	Rp       string
	Type     string
	Task     Task
}

func (t *AlertTest) Init(c TestConf) error {
	t.Name = c.Name
	t.Db = c.Db
	t.Rp = c.Rp
	t.TaskName = c.TaskName
	t.Expects = c.Expects
	t.Data = c.Data
	t.Type = c.Type
	return nil
}

// Method exposed to start the model. It sets up the model, adds the model data,
// fetches the triggered alerts and saves it. It also removes all artifacts
// (database, retention policy) created for the model.
func (t *AlertTest) Run(k io.Kapacitor, i io.Influxdb) error {

	err := t.setup(k, i)
	defer t.teardown(k, i) //defer teardown so it gets run incase of early termination

	if err != nil {
		return err
	}
	err = t.addData(k, i)
	if err != nil {
		return err
	}

	t.wait()

	err = t.results(k)
	if err != nil {
		return err
	}

	return nil
}

func (t AlertTest) GetName() string {
	return t.TaskName
}

func (t AlertTest) GetType() string {
	return t.Type
}

func (t *AlertTest) SetTask(tsk Task) {
	t.Task = tsk
}

func (t AlertTest) Passed() bool {
	return t.Result.Passed
}

func (t AlertTest) String() string {
	if t.Result.Error == true {
		return fmt.Sprintf("TEST %v (%v) ERROR: %v", t.Name, t.TaskName, t.Result.String())
	} else {
		return fmt.Sprintf("TEST %v (%v) %v", t.Name, t.TaskName, t.Result.String())
	}
}

// Adds model data
func (t *AlertTest) addData(k io.Kapacitor, i io.Influxdb) error {
	switch t.Type {
	case "stream":
		// adds data to kapacitor
		err := k.Data(t.Data, t.Db, t.Rp)
		if err != nil {
			return err
		}
	case "batch":
		// adds data to InfluxDb
		err := i.Data(t.Data, t.Db, t.Rp)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validates if individual model configuration is correct
func (t *AlertTest) Validate() error {
	glog.Info("DEBUG:: validate model: ", t.Name)
	if len(t.Data) > 0 && t.RecId != "" {
		m := "Configuration file cannot define a recording_id and line protocol data input for the same model case"
		r := Result{0, 0, 0, m, false, true}
		t.Result = r
	}
	return nil
}

// Creates all necessary artifacts in database to run the model
func (t *AlertTest) setup(k io.Kapacitor, i io.Influxdb) error {
	glog.Info("DEBUG:: setup model: ", t.Name)
	switch t.Type {
	case "batch":
		err := i.Setup(t.Db, t.Rp)
		if err != nil {
			return err
		}
	}

	// Loads model task to kapacitor
	f := map[string]interface{}{
		"id":     t.TaskName,
		"type":   t.Type,
		"script": t.Task.Script,
		"status": "enabled",
	}

	dbrp, _ := regexp.MatchString(`(?m:^dbrp \"\w+\"\.\"\w+\"$)`, t.Task.Script)
	if !dbrp {
		f["dbrps"] = []map[string]string{{"db": t.Db, "rp": t.Rp}}
	}

	err := k.Load(f)
	if err != nil {
		return err
	}
	return nil
}

func (t *AlertTest) wait() {
	switch t.Type {
	case "batch":
		// If batch script, waits 3 seconds for batch queries being processed
		fmt.Println("Processing batch script " + t.TaskName + "...")
		time.Sleep(3 * time.Second)
	}
}

// Deletes data, database and retention policies created to run the model
func (t *AlertTest) teardown(k io.Kapacitor, i io.Influxdb) {
	glog.Info("DEBUG:: teardown model: ", t.Name)
	switch t.Type {
	case "batch":
		err := i.CleanUp(t.Db)
		if err != nil {
			glog.Error("Error performing teardown in cleanup. error: ", err)
		}
	}
	err := k.Delete(t.TaskName)
	if err != nil {
		glog.Error("Error performing teardown in delete error: ", err)
	}
}

// Fetches status of kapacitor task, stores it and compares expected model result
// and actual result model
func (t *AlertTest) results(k io.Kapacitor) error {
	s, err := k.Status(t.Task.Name)
	if err != nil {
		return err
	}
	t.Result = NewResult(s)
	t.Result.Compare(t.Expects)
	return nil
}
