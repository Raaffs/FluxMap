package models

import (
    "database/sql"
)

// User represents a user in the database
type User struct {
    Username       string    `json:"username"`       // Primary Key
    Email          string    `json:"email" validate:"required,email"`
    Password       string    `json:"password" validate:"required,min=8"`
    HashedPassword string    `json:"hashedPassword,omitempty"`
    Created        string    `json:"created,omitempty"`
}

// Project represents a project in the database
type Project struct {
    ProjectID          int            `json:"projectID,omitempty"`           // Primary Key
    ProjectName        string         `json:"projectName" validate:"required"`
    ProjectDescription sql.NullString `json:"projectDescription,omitempty"`
    ProjectStartDate   sql.NullTime   `json:"projectStartDate,omitempty"`
    ProjectDueDate     sql.NullTime   `json:"projectDueDate,omitempty"`
    Ownername          string         `json:"ownername"`           // Foreign Key (User.Username)
}

// Manager represents a manager in the database
type Manager struct {
    Manager   string `json:"manager"`    // Primary Key, Foreign Key (User.Username)
    ProjectID int    `json:"projectId"`  // Foreign Key (Project.ProjectID)
}

// Task represents a task in the database
type Task struct {
    TaskID           int            `json:"taskID,omitempty"`               // Primary Key
    TaskName         string         `json:"taskName" validate:"required"`
    TaskDescription  sql.NullString `json:"taskDescription,omitempty"`
    TaskStatus       sql.NullString `json:"taskStatus"`
    TaskStartDate    sql.NullTime   `json:"taskStartDate,omitempty"`
    TaskDueDate      sql.NullTime   `json:"taskDueDate,omitempty"`
    ParentProjectID  int            `json:"parentProjectId"`      // Foreign Key (Project.ProjectID)
    AssignedUsername sql.NullString `json:"assignedUsername" validate:"required"` // Foreign Key (User.Username)
}

// Pert represents a PERT record in the database
type Pert struct {
    ParentTaskID         int            `json:"parentTaskId"`         // Primary Key, Foreign Key (Task.TaskID)
    PredecessorTaskID    sql.NullInt64  `json:"predecessorTaskId,omitempty"` // Foreign Key (Task.TaskID)
    Optimistic           int            `json:"optimistic" validate:"required"`
    Pessimistic          int            `json:"pessimistic" validate:"required"`
    MostLikely           int            `json:"mostLikely" validate:"required"`
}

// Cpm represents a CPM record in the database
type Cpm struct {
    TaskID         int    `json:"taskId"`             // Primary Key, Foreign Key (Task.TaskID)
    EarliestStart  int    `json:"earliestStart" validate:"required"`
    EarliestFinish int    `json:"earliestFinish" validate:"required"`
    LatestStart    int    `json:"latestStart" validate:"required"`
    LatestFinish   int    `json:"latestFinish" validate:"required"`
    SlackTime      int    `json:"slackTime" validate:"required"`
    CriticalPath   bool   `json:"criticalPath" default:"false"`
}