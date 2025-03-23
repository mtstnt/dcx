package main

import "time"

type File struct {
	Encoding string `json:"Encoding"`
	Content  string `json:"Content"`
}

type Test struct {
	Input  string `json:"Input"`
	Output string `json:"Output"`
}

type Configuration struct {
	Strictness  string `json:"Strictness"`
	Timeout     int    `json:"Timeout"`
	MemoryLimit int    `json:"MemoryLimit"`
}

type SubmitRequest struct {
	Files          map[string]File `json:"Files"`
	Image          string          `json:"Image"`
	Tests          []Test          `json:"Tests"`
	Configurations Configuration   `json:"Configurations"`
}

type ExecutionStatus struct {
	Status    string    `json:"Status"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type SubmitResponse struct {
	ExecutionGroupID string            `json:"ExecutionGroupID"`
	Executions       map[string]string `json:"Executions"` // ExecutionID -> Status
}
