package common

type Status string

const (
	TEMPLATE Status = "TEMPLATE"
	NEW      Status = "NEW"
	STARTED  Status = "STARTED"
	FINISHED Status = "FINISHED"
)

type OrderType string

const (
	PARALLEL_ANY_OF    OrderType = "PARALLEL_ANY_OF"
	PARALLEL_ALL_OF    OrderType = "PARALLEL_ALL_OF"
	SEQUENTIAL_ALL_OFF OrderType = "SEQUENTIAL_ALL_OFF"
)

type Decision string

const (
	UNKNOWN Decision = "UNKNOWN"
	ACCEPT  Decision = "ACCEPT"
	REJECT  Decision = "REJECT"
)
