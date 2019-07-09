// Keep data about a task to be tested and interface to run all task's tests
package model

import (
	"encoding/json"
	"github.com/gpestana/kapacitor-unit/io"
	"io/ioutil"
	"strings"
)

// FS configurations, namely path where TICKscripts are located
type Task struct {
	Name   string
	Path   string
	Script string
	Type   string
	Dbrps  DbRps
}

type Dbrp struct {
	Db string `json:"db"`
	Rp string `json:"rp"`
}

type DbRps []Dbrp

type TaskPayload struct {
	Script string `json:"script"`
	DbRps  DbRps  `json:"dbrps"`
	Id     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func LoadTask(task Task, dbrps DbRps, k io.Kapacitor) error {
	script := task.Script
	if task.Type == "batch" {
		script = k.BatchReplaceEvery(script)
	}
	payload := TaskPayload{
		Script: script,
		DbRps:  dbrps,
		Id:     task.Name,
		Type:   task.Type,
		Status: "enabled",
	}

	j, err := json.Marshal(payload);
	if err != nil {
		return err
	}

	return k.PostTask(j)
}

// Task constructor
func NewTask(tc TestConf, config Config) (*Task, error) {

	task := Task{
		Name: tc.TaskName,
		Path: config.ScriptsDir, // do we need this?
		Type: tc.Type,
	}

	if s, e := readScriptFromFs(tc.TaskName, config.ScriptsDir); e != nil {
		return nil, e
	} else {
		task.Script = s
	}

	task.Dbrps = make(DbRps, 0)
	if tc.SeedDb != "" {
		task.Dbrps = append(task.Dbrps, Dbrp{Db: tc.SeedDb, Rp: tc.SeedRp})
	}
	if tc.AssertDb != "" {
		task.Dbrps = append(task.Dbrps, Dbrp{Db: tc.AssertDb, Rp: tc.AssertRp})
	}

	return &task, nil
}

func readScriptFromFs(n string, p string) (string, error) {
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	s, err := ioutil.ReadFile(p + n)
	return string(s), err
}
