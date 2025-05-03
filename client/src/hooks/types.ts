export interface Projects {
  projectID?: number; // Matches `omitempty`
  projectName: string; // Required field
  projectDescription?: string | null; // Matches `null.String`
  projectStartDate?: string | null; // Matches `null.Time`
  projectDueDate?: string | null; // Matches `null.Time`
  ownername: string; // Required field
}

export interface tasks {
    taskID: number;
    taskName: string;           // Required
    taskDescription?: string;   // Optional
    taskStatus?: string|"";        // Optional
    taskStartDate?: string|null;     // Optional, ISO date string
    taskDueDate?: string|null;       // Optional, ISO date string
    parentProjectID: number;    // Required, Foreign Key (Project.ProjectID)
    assignedUsername: string;   // Required, Foreign Key (User.Username)
    approved?: boolean|false;         // Optional
    taskCompletedDate?:  string|null
    taskApprovedDate?:  string|null
    createdBy:string;
}

export interface PertData {
    parentTaskID: number;        // Required, Foreign Key (Task.TaskID)
    predecessorTaskId?: number|null;  // Optional, Foreign Key (Task.TaskID)
    optimistic: number;          // Required
    pessimistic: number;         // Required
    mostLikely: number;          // Required
}

interface PertTaskResult {  
    mean:   number;
    stddev: number;
    taskId: number;
  }
  
export  interface ApiResponse {
    data: PertData[];
    result: {
      criticalPath: number[];
      taskResults: PertTaskResult[];
    };
}


export interface CpmApiResponse {
    taskId: number;
    parentProjectID: number;
    dependencies: number[];
    duration: number;
    earliestStart: number;
    earliestFinish: number;
    latestStart: number;
    latestFinish: number;
    totalFloat: number;
    freeFloat: number;
    independentFloat: number;
    isCriticalPath: boolean;
  }
  
export interface CpmResult {
    Result: CpmApiResponse[];
}

export interface DisplayInvitations {
  projectID: number;
  projectName: string;
  projectDescription: string;
  ownerName: string;
  invitationID: number;
  role: string;
  invitationUsername: string;
}

export interface Invitation {
  id: string;
  username: string;
  projectID?: number; 
  accepted: boolean;
  role: string;
}

export interface Update {
  id: number;
  projectID: number;
  msg: string;
  createdAt: string; // ISO timestamp string
  createdBy: string; // references users.username
  targetType: 'all' | 'user';
  targetUsername: string | null; // nullable, only set if targetType = 'user'
  projectName: string;
  projectDescription: string;
}

// let u :Update[]=[]

// const grouped: Record<number, Update[]> = {};
// u.forEach(update => {
//   (grouped[update.projectID] ||= []).push(update);
// });