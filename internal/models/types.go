package models

import (
	"context"
	"time"

	"github.com/guregu/null/v5"
)

type ReadDatabase[T  Analytic ] interface{
    Exist()(bool,error)
    InsertResult(context.Context,int,Result)error
    GetData(context.Context,int)([]*T,error)
    GetResult(context.Context,int)(Result,error)
}

type Analytic interface{
    Cpm|Pert
}

// User represents a user in the database
type User struct {
    Username       string    `json:"username"`       // Primary Key
    Email          string    `json:"email" validate:"required,email"`
    Password       string    `json:"password" validate:"required,min=8"`
    HashedPassword string    `json:"hashedPassword,omitempty"`
    Created        string    `json:"created,omitempty"`
}

type Invitation struct{
    ID                  string         `json:"id"`
    Username            string         `json:"username"`       
    ProjectID           int            `json:"projectID,omitempty"`           // Primary Key
    Status              string           `json:"status"`
    Role                string         `json:"role"`
}

type DisplayInvitations struct {
	ProjectID           int    `json:"projectID"`
	ProjectName         string `json:"projectName"`
	ProjectDescription  string `json:"projectDescription"`
	OwnerName           string `json:"ownerName"`
	InvitationID        int    `json:"invitationID"`
	Role                string `json:"role"`
	Username            string `json:"invitationUsername"`
}

type Update struct {
    ID             int              `db:"id"`
    ProjectID      int              `db:"projectid"`
    Msg            string           `db:"msg"`
    CreatedAt      time.Time        `db:"createdat"`
    CreatedBy      string           `db:"createdby"`      // references users.username
    TargetType     string           `db:"targettype"`     // either "all" or "user"
    TargetUsername null.String      `db:"targetusername"` // nullable, only set if TargetType = "user"
}


// Project represents a project in the database
type Project struct {
    ProjectID          int            `json:"projectID,omitempty"`           // Primary Key
    ProjectName        string         `json:"projectName" validate:"required"`
    ProjectDescription null.String    `json:"projectDescription,omitempty"`
    ProjectStartDate   null.Time      `json:"projectStartDate,omitempty"`
    ProjectDueDate     null.Time      `json:"projectDueDate,omitempty"`
    Ownername          string         `json:"ownername,omitempty"`           // Foreign Key (User.Username)
}

// Manager represents a manager in the database
type Manager struct {
    Manager   string `json:"manager"`    // Primary Key, Foreign Key (User.Username)
    ProjectID int    `json:"projectId"`  // Foreign Key (Project.ProjectID)
}

// Task represents a task in the database
type Task struct {
    TaskID              int            `json:"taskID,omitempty"`               // Primary Key
    TaskName            string         `json:"taskName" validate:"required"`
    TaskDescription     null.String    `json:"taskDescription,omitempty"`
    TaskStatus          null.String    `json:"taskStatus"`
    TaskStartDate       null.Time      `json:"taskStartDate,omitempty"`
    TaskDueDate         null.Time      `json:"taskDueDate,omitempty"`
    ParentProjectID     int            `json:"parentProjectId"`      // Foreign Key (Project.ProjectID)
    AssignedUsername    null.String    `json:"assignedUsername" validate:"required"` // Foreign Key (User.Username)
    Approved            null.Bool      `json:"approved"`
    TaskCompletedDate   null.Time      `json:"taskCompletedDate"`
    TaskApprovedDate    null.Time      `json:"taskApprovedDate"`
    Createdby           string 
}

// Pert represents a PERT record in the database
type Pert struct {
    ParentTaskID         int                `json:"parentTaskID"`         // Primary Key, Foreign Key (Task.TaskID)
    PredecessorTaskID    null.Int64         `json:"predecessorTaskId,omitempty"` // Foreign Key (Task.TaskID)
    Optimistic           int                `json:"optimistic" validate:"required"`
    Pessimistic          int                `json:"pessimistic" validate:"required"`
    MostLikely           int                `json:"mostLikely" validate:"required"`
    ParentProjectID      int
}


type Cpm struct {
    TaskID          int                 `json:"taskId"`             // Primary Key, Foreign Key (Task.TaskID)
    ParentProjectID int                 `json:"parentProjectID"`
    Dependencies    []int               `json:"dependencies"`  
    Duration        int                 `json:"duration"`
}

type Result struct{         
    Result map[string]any
}