package common

type Status string

const (
	TEMPLATE Status = "TEMPLATE"
	NEW      Status = "NEW"
	STARTED  Status = "STARTED"
	FINISHED Status = "FINISHED"
)
