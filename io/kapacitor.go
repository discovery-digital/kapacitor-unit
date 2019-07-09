package io

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Status struct {
	Data map[string]map[string]interface{} `json:"stats"`
}

// Kapacitor service configurations
type Kapacitor struct {
	Host   string
	Client http.Client
}

func NewKapacitor(host string) Kapacitor {
	return Kapacitor{
		host,
		http.Client{},
	}
}

// Loads a task
func (k Kapacitor) Load(f map[string]interface{}) error {
	glog.Info("DEBUG:: Kapacitor loading task: ", f["id"])
	// Replaces '.every()' if type of script is batch
	if f["type"] == "batch" {
		str, ok := f["script"].(string)
		if ok != true {
			return errors.New("Task Load: script is not of type string")
		}
		f["script"] = k.BatchReplaceEvery(str)

		glog.Info("DEBUG:: batch script after replace: ", f["script"])
	}

	j, err := json.Marshal(f)
	if err != nil {
		return err
	}

	return k.PostTask(j)
}

func (k Kapacitor) PostTask(payload []byte) error {
	u := k.Host + tasks
	res, err := k.Client.Post(u, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		r, _ := ioutil.ReadAll(res.Body)
		return errors.New(res.Status + ":: " + string(r))
	}
	return nil
}


// Deletes a task
func (k Kapacitor) Delete(id string) error {
	u := k.Host + tasks + "/" + id
	r, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	_, err = k.Client.Do(r)
	if err != nil {
		return err
	}
	glog.Info("DEBUG:: Kapacitor deleted task: ", id)
	return nil
}

// Adds model data to kapacitor
func (k Kapacitor) Data(data []string, db string, rp string) error {
	u := k.Host + kapacitor_write + "db=" + db + "&rp=" + rp
	for _, d := range data {
		resp, err := k.Client.Post(u, "application/x-www-form-urlencoded",
			bytes.NewBuffer([]byte(d)))

		var buffer []byte
		resp.Body.Read(buffer)
		fmt.Println(buffer)

		if err != nil {
			return err
		}
		glog.Info("DEBUG:: Kapacitor added data: ", d)
	}
	return nil
}

func (k Kapacitor) ListenLog(id string) {
	url := k.Host + logs + id

	response, err := k.Client.Get(url)

	if err != nil {
		return //nil, err
	}

	logs := make([][]byte, 0)

	reader := bufio.NewReader(response.Body)
	for {
		line, _ := reader.ReadBytes('\n')

		fmt.Println(string(line))

		logs = append(logs, line)
	}

	return //logs
}

// Gets task alert status
func (k Kapacitor) Status(id string) (map[string]int, error) {
	glog.Info("DEBUG:: Kapacitor fetching status of: ", id)
	u := k.Host + tasks + "/" + id
	res, err := k.Client.Get(u)
	if err != nil {
		return nil, err
	}
	var s Status
	b, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	f := make(map[string]int)
	var sa interface{}
	for key, value := range s.Data["node-stats"] {
		if strings.HasPrefix(key, "alert") {
			sa = value
			for k, val := range sa.(map[string]interface{}) {
				switch v := val.(type) {
				case float64:
					f[k] += int(v)
				default:
					return nil, errors.New("kapacitor.status: wrong response from service")
				}
			}
		}
	}
	if sa == nil {
		return nil, errors.New("kapacitor.status: expected alert.* key to be found on stats")
	}
	return f, nil
}

// Replaces '.every(*)' for the batch request to be performed every 1s to speed up the model
func (k Kapacitor) BatchReplaceEvery(s string) string {
	re := regexp.MustCompile("every\\((.*?)\\)")
	return re.ReplaceAllString(s, "every(1s)")
}
