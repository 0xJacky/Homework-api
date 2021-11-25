package model

import "time"

type Homework struct {
	Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Upload      Upload    `json:"upload"`
}
