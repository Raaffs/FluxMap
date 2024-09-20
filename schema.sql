-- Switch to mapmyprojectv2 database
drop table tasks cascade;
-- Users table
CREATE TABLE users (
    username VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    hashedPassword CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    userStatus BOOLEAN DEFAULT TRUE -- active/inactive status
);
-- Projects table
CREATE TABLE projects(
    projectID SERIAL PRIMARY KEY,
    projectName VARCHAR(255) NOT NULL,
    projectDescription TEXT,
    projectStartDate TIMESTAMP NOT NULL,
    projectDueDate TIMESTAMP,
    ownername VARCHAR(255) NOT NULL,
    FOREIGN KEY (ownername) REFERENCES users(username) ON DELETE CASCADE
);
    
create table managers(
	manager VARCHAR(255) PRIMARY KEY,
	projectID	INT,
	FOREIGN KEY (manager) REFERENCES users(username) ON DELETE cascade,
	FOREIGN KEY (projectID) REFERENCES projects(projectID) ON DELETE CASCADE
);


-- Tasks table
CREATE TABLE tasks(
    taskID SERIAL primary KEY,
    taskName VARCHAR(255) NOT NULL,
    taskDescription TEXT,
    taskStatus VARCHAR(255),
    taskStartDate TIMESTAMP,
    taskDueDate TIMESTAMP,
    parentProjectID INTEGER NOT NULL,
    assignedUsername VARCHAR(255),
    FOREIGN KEY (assignedUsername) REFERENCES users(username) ON DELETE SET NULL,
    FOREIGN KEY (parentProjectID) REFERENCES projects(projectID) ON DELETE CASCADE
);

-- PERT table with composite key references
CREATE TABLE pert(
    parentTaskID INTEGER NOT NULL,
    predecessorTaskID INTEGER,
    predecessorProjectID INTEGER,
    optimistic INTEGER NOT NULL,
    pessimistic INTEGER NOT NULL,
    mostLikely INTEGER NOT NULL,
    PRIMARY KEY (parentTaskID),
    FOREIGN KEY (parentTaskID) REFERENCES tasks(taskID) ON DELETE CASCADE,
    FOREIGN KEY (predecessorTaskID) REFERENCES tasks(taskID) ON DELETE SET NULL
);

-- CPM table with composite key references
CREATE TABLE cpm(
    taskID INTEGER NOT NULL,
    earliestStart INTEGER NOT NULL,
    earliestFinish INTEGER NOT NULL,
    latestStart INTEGER NOT NULL,
    latestFinish INTEGER NOT NULL,
    slackTime INTEGER NOT NULL,
    criticalPath BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (taskID),
    FOREIGN KEY (taskID) REFERENCES tasks(taskID) ON DELETE CASCADE
);