// io package is responsable for all interactions with disk and external
// services such as Kapacitor and Influxdb. Its main goal is to read tasks from
// disk, load, read and delete tasks from kapacitor as well as check the alert
// logs. It also is responsible for loading and deleting model data into
// Influxdb as well as creating and deleting the necessary model databases.
package io

const (
	kapacitor_write = "/kapacitor/v1/write?"
	influxdb_write = "/write?"
	tasks = "/kapacitor/v1/tasks"
	logs = "/kapacitor/v1preview/logs?task="
)
