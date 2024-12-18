export interface Projects {
    projectID: number;
    projectName: string;  // Required
    projectDescription?: string|null; // Optional
    projectStartDate: string; // ISO date string (e.g., "2024-01-01")
    projectDueDate: string;   // ISO date string
}

export interface tasks {
    taskID: number;
    taskName: string;           // Required
    taskDescription?: string;   // Optional
    taskStatus?: string;        // Optional
    taskStartDate?: string;     // Optional, ISO date string
    taskDueDate?: string;       // Optional, ISO date string
    parentProjectID: number;    // Required, Foreign Key (Project.ProjectID)
    assignedUsername: string;   // Required, Foreign Key (User.Username)
    approved?: boolean;         // Optional
}

export interface PertInput {
    parentTaskID: number;        // Required, Foreign Key (Task.TaskID)
    predecessorTaskId?: number;  // Optional, Foreign Key (Task.TaskID)
    optimistic: number;          // Required
    pessimistic: number;         // Required
    mostLikely: number;          // Required
}


export interface CpmInput {
    taskID: number;              // Required, Foreign Key (Task.TaskID)
    earliestStart: number;       // Required
    earliestFinish: number;      // Required
    latestStart: number;         // Required
    latestFinish: number;        // Required
    slackTime: number;           // Required
    criticalPath?: boolean;      // Optional, defaults to false
    parentProjectId: number;     // Required
    dependencies?: number[];     // Optional
}
  
