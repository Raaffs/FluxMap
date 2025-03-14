def calculateCpm(tasks):
    """
    Calculate CPM fields for tasks.

    :param tasks: List of tasks where each task is a dictionary with:
                  - taskId: int
                  - projectId: int
                  - dependencies: list of task IDs this task depends on
                  - duration: int
    :return: List of tasks with calculated fields (earliestStart, earliestFinish, etc.)
    """
    # Create a dictionary for quick lookup by taskId
    taskDict = {task["taskId"]: task for task in tasks}

    # Add required fields for CPM calculations
    for task in tasks:
        task["earliestStart"] = 0
        task["earliestFinish"] = 0
        task["latestStart"] = -1  # Use -1 as a placeholder for infinity
        task["latestFinish"] = -1  # Use -1 as a placeholder for infinity
        task["totalFloat"] = 0
        task["freeFloat"] = 0
        task["independentFloat"] = 0
        task["isCriticalPath"] = False

    # Forward pass to calculate earliestStart (ES) and earliestFinish (EF)
    for task in tasks:
        if task["dependencies"]:
            task["earliestStart"] = max(
                taskDict[dep]["earliestFinish"] for dep in task["dependencies"]
            )
        task["earliestFinish"] = task["earliestStart"] + task["duration"]

    # Backward pass to calculate latestStart (LS) and latestFinish (LF)
    projectFinishTime = max(task["earliestFinish"] for task in tasks)
    for task in reversed(tasks):
        if not any(t["taskId"] in task["dependencies"] for t in tasks):
            task["latestFinish"] = projectFinishTime
        else:
            task["latestFinish"] = min(
                taskDict[dep]["latestStart"] for dep in task["dependencies"]
            )
        task["latestStart"] = task["latestFinish"] - task["duration"]

    # Calculate floats and determine the critical path
    for task in tasks:
        task["totalFloat"] = task["latestStart"] - task["earliestStart"]
        task["freeFloat"] = min(
            (taskDict[dep]["earliestStart"] - task["earliestFinish"])
            for dep in task["dependencies"]
        ) if task["dependencies"] else task["totalFloat"]
        task["independentFloat"] = max(
            0,
            task["freeFloat"] - task["totalFloat"]
        )
        task["isCriticalPath"] = task["totalFloat"] == 0

    # Return the updated task list with replaced placeholders for infinity/None values
    return tasks
