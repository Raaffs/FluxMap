def calculateCpm(tasks):
    # Ensure dependencies are lists, never None
    for task in tasks:
        if task.get("dependencies") is None:
            task["dependencies"] = []

    taskDict = {task["taskId"]: task for task in tasks if task is not None and "taskId" in task}

    # Build successors list for backward pass
    for task in tasks:
        task["successors"] = []

    for task in tasks:
        for depId in task["dependencies"]:
            if depId in taskDict:
                taskDict[depId]["successors"].append(task["taskId"])

    # Initialize CPM fields
    for task in tasks:
        task["earliestStart"] = 0
        task["earliestFinish"] = 0
        task["latestStart"] = -1
        task["latestFinish"] = -1
        task["totalFloat"] = 0
        task["freeFloat"] = 0
        task["independentFloat"] = 0
        task["isCriticalPath"] = False

    # Forward pass function
    def forward(task):
        if task["dependencies"]:
            valid_deps = [dep for dep in task["dependencies"] if dep in taskDict]
            if valid_deps:
                task["earliestStart"] = max(taskDict[dep]["earliestFinish"] for dep in valid_deps)
            else:
                task["earliestStart"] = 0
        else:
            task["earliestStart"] = 0
        task["earliestFinish"] = task["earliestStart"] + task["duration"]

    # Run forward pass until stable
    changed = True
    while changed:
        changed = False
        for task in tasks:
            old_es = task["earliestStart"]
            forward(task)
            if task["earliestStart"] != old_es:
                changed = True

    projectFinishTime = max(task["earliestFinish"] for task in tasks)

    def backward(task):
        valid_succs = [succ for succ in task.get("successors", []) if succ in taskDict]
        if not valid_succs:
            task["latestFinish"] = projectFinishTime
        else:
            task["latestFinish"] = min(taskDict[succ]["latestStart"] for succ in valid_succs)
        task["latestStart"] = task["latestFinish"] - task["duration"]

    changed = True
    while changed:
        changed = False
        for task in reversed(tasks):
            old_lf = task["latestFinish"]
            backward(task)
            if task["latestFinish"] != old_lf:
                changed = True

    # Calculate floats and critical path flag, clean up
    for task in tasks:
        task["totalFloat"] = task["latestStart"] - task["earliestStart"]
        valid_succs = [succ for succ in task.get("successors", []) if succ in taskDict]
        if valid_succs:
            task["freeFloat"] = min(taskDict[succ]["earliestStart"] for succ in valid_succs) - task["earliestFinish"]
        else:
            task["freeFloat"] = task["totalFloat"]
        task["independentFloat"] = max(0, task["freeFloat"] - task["totalFloat"])
        task["isCriticalPath"] = (task["totalFloat"] == 0)
        del task["successors"]  # remove temp field

    return tasks
