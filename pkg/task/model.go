package task

import (
	"time"

	"dojo-api/db"
)

// TaskResponse reflects the task structure used in API responses
type TaskResponse struct {
	ID         string        `json:"taskId"`
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	ExpireAt   time.Time     `json:"expireAt"`
	Type       db.TaskType   `json:"type"`
	TaskData   db.JSON       `json:"taskData"`
	Status     db.TaskStatus `json:"status"`
	MaxResults int           `json:"maxResults"`
}

type SortField string

const (
	SortCreatedAt    SortField = "createdAt"
	SortNumResult    SortField = "numResult"
	SortHighestYield SortField = "highestYield"
)

type Pagination struct {
	Page       int `json:"pageNumber"`
	Limit      int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
}

type TaskPagination struct {
	Tasks      []TaskResponse `json:"tasks"`
	Pagination Pagination     `json:"pagination"`
}

type CreateTaskRequest struct {
	Title        string      `json:"title"`
	Body         string      `json:"body"`
	ExpireAt     interface{} `json:"expireAt"`
	TaskData     []TaskData  `json:"taskData"`
	MaxResults   int         `json:"maxResults"`
	TotalRewards float64     `json:"totalRewards"`
}
type SubmitTaskResultRequest struct {
	TaskId       string                 `json:"taskId"`
	DojoWorkerId string                 `json:"dojoWorkerId"`
	ResultData   map[string]interface{} `json:"resultData"`
}

type TaskData struct {
	Prompt    string          `json:"prompt"`
	Dialogue  []Message       `json:"dialogue,omitempty"`
	Responses []ModelResponse `json:"responses,omitempty"`
	Task      TaskType        `json:"task"`
	Criteria  []Criteria      `json:"criteria"`
}

type TaskType string

const (
	TaskTypeCodeGen     TaskType = TaskType(db.TaskTypeCodeGeneration)
	TaskTypeDialogue    TaskType = TaskType(db.TaskTypeDialogue)
	TaskTypeTextToImage TaskType = TaskType(db.TaskTypeTextToImage)
)

type ModelResponse struct {
	Model      string      `json:"model"`
	Completion interface{} `json:"completion"`
}

type Message struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type Criteria struct {
	Type    CriteriaType `json:"type"`
	Options []string     `json:"options,omitempty"`
	Min     float64      `json:"min,omitempty"`
	Max     float64      `json:"max,omitempty"`
}

type CriteriaType string

const (
	CriteriaTypeRanking     CriteriaType = "ranking"
	CriteriaTypeMultiSelect CriteriaType = "multi-select"
	CriteriaTypeScore       CriteriaType = "score"
)
