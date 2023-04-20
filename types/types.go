package types

type Level int
type Status int

const (
	LOW Level = iota
	MEDIUM
	HIGH
	CRITICAL
)

const (
	CREATED Status = iota
	IN_PROGRESS
	TESTING
	COMPLETED
	CLOSED
)

type Task struct {
	//Id          string `json:"Moid" bson:"_id"`
	Name        string `json:"Name" bson:"name"`
	Description string `json:"Description" bson:"description"`
	Priority    Level  `json:"Priority" bson:"priority"`
	EndDate     string `json:"EndDate" bson:"endDate"`
	Status      Status `json:"Status" bson:"status"`
}

type TaskBase struct {
	Id   string `json:"TaskId" bson:"_id"`
	Task `bson:",inline"`
}
