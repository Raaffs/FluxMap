import numpy as np
from scipy.stats import norm

# Function to compute mean and standard deviation for task duration
def compute_distribution(task):
    duration_mean = (task["earliestFinish"] + task["earliestStart"]) / 2
    duration_std = (task["earliestFinish"] - task["earliestStart"]) / 6  # Approx. std dev using range / 6 rule
    return duration_mean, duration_std

# Function to calculate floats and distributions for tasks
def get_task_distributions(tasks):
    task_distributions = []
    overall_mean = 0
    overall_variance = 0

    for task in tasks:
        mean, std_dev = compute_distribution(task)
        task_distributions.append({
            "taskId": task["taskId"],
            "mean": mean,
            "std_dev": std_dev,
            "total_float": task["latestStart"] - task["earliestStart"],
            "free_float": compute_free_float(task, tasks),
            "independent_float": task["latestStart"] - task["earliestFinish"] - task["successorSlack"]
        })
        
        # Sum means and variances for overall distribution
        overall_mean += mean
        overall_variance += std_dev ** 2

    # Overall project mean and standard deviation
    overall_std = np.sqrt(overall_variance)
    task_distributions.append({
        "taskId": "overall",
        "mean": overall_mean,
        "std_dev": overall_std
    })

    return task_distributions

def compute_free_float(task, tasks):
    # Find the earliest start time of the next task
    next_tasks = [t for t in tasks if task["taskId"] in t.get("dependencies", [])]
    if not next_tasks:
        return float('inf')  # No successors mean infinite free float

    earliest_start_next = min(nt["earliestStart"] for nt in next_tasks)
    return task["earliestFinish"] - earliest_start_next

def get_critical_path_distributions(tasks):
    critical_path_distributions = []
    critical_path_mean = 0
    critical_path_variance = 0

    for task in tasks:
        if task.get("criticalPath", False):  # Check if the task is on the critical path
            mean, std_dev = compute_distribution(task)
            critical_path_distributions.append({
                "taskId": task["taskId"],
                "mean": mean,
                "std_dev": std_dev,
                "total_float": task["latestStart"] - task["earliestStart"],
                "free_float": compute_free_float(task, tasks),
                "independent_float": task["latestStart"] - task["earliestFinish"] - task["successorSlack"]
            })

            critical_path_mean += mean
            critical_path_variance += std_dev ** 2

    critical_path_std = np.sqrt(critical_path_variance)
    critical_path_distributions.append({
        "taskId": "criticalPathOverall",
        "mean": critical_path_mean,
        "std_dev": critical_path_std
    })

    return critical_path_distributions

# Example of tasks data
tasks = [
    {
        "taskId": "A",
        "earliestStart": 0,
        "earliestFinish": 4,
        "latestStart": 0,
        "dependencies": [],
        "successorSlack": 0,
        "criticalPath": True
    },
    {
        "taskId": "B",
        "earliestStart": 4,
        "earliestFinish": 8,
        "latestStart": 4,
        "dependencies": ["A"],
        "successorSlack": 0,
        "criticalPath": True
    },
    {
        "taskId": "C",
        "earliestStart": 2,
        "earliestFinish": 6,
        "latestStart": 4,
        "dependencies": ["A"],
        "successorSlack": 2,
        "criticalPath": False
    },
    {
        "taskId": "D",
        "earliestStart": 8,
        "earliestFinish": 10,
        "latestStart": 8,
        "dependencies": ["B", "C"],
        "successorSlack": 0,
        "criticalPath": True
    }
]

# Usage
task_distributions = get_task_distributions(tasks)
critical_path_distributions = get_critical_path_distributions(tasks)

print("Task Distributions:")
for dist in task_distributions:
    print(dist)

print("\nCritical Path Distributions:")
for dist in critical_path_distributions:
    print(dist)
