import numpy as np
from scipy.stats import norm

# Define the CPM data for tasks
tasks = [
    {
        "taskId": 1,
        "earliestStart": 2,
        "earliestFinish": 5,
        "latestStart": 3,
        "latestFinish": 6,
        "slackTime": 1,
        "criticalPath": False
    },
    {
        "taskId": 3,
        "earliestStart": 4,
        "earliestFinish": 8,
        "latestStart": 6,
        "latestFinish": 9,
        "slackTime": 0,
        "criticalPath": True
    },
    {
        "taskId": 6,
        "earliestStart": 7,
        "earliestFinish": 10,
        "latestStart": 8,
        "latestFinish": 11,
        "slackTime": 2,
        "criticalPath": False
    }
]

# Function to compute mean and standard deviation for task duration
def compute_distribution(task):
    duration_mean = (task["earliestFinish"] + task["earliestStart"]) / 2
    duration_std = (task["earliestFinish"] - task["earliestStart"]) / 6  # Approx. std dev using range / 6 rule
    return duration_mean, duration_std

# Function to return task and overall distributions
def get_task_distributions(tasks):
    task_distributions = []
    overall_mean = 0
    overall_variance = 0

    for task in tasks:
        mean, std_dev = compute_distribution(task)
        task_distributions.append({
            "taskId": task["taskId"],
            "mean": mean,
            "std_dev": std_dev
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

# Function to return critical path distributions
def get_critical_path_distributions(tasks):
    critical_path_distributions = []
    critical_path_mean = 0
    critical_path_variance = 0

    for task in tasks:
        if task["criticalPath"]:
            mean, std_dev = compute_distribution(task)
            critical_path_distributions.append({
                "taskId": task["taskId"],
                "mean": mean,
                "std_dev": std_dev
            })

            # Sum means and variances for critical path
            critical_path_mean += mean
            critical_path_variance += std_dev ** 2

    # Overall critical path mean and standard deviation
    critical_path_std = np.sqrt(critical_path_variance)
    critical_path_distributions.append({
        "taskId": "criticalPathOverall",
        "mean": critical_path_mean,
        "std_dev": critical_path_std
    })

    return critical_path_distributions

# Get task distributions
task_distributions = get_task_distributions(tasks)

# Get critical path distributions
critical_path_distributions = get_critical_path_distributions(tasks)

# Output results
print("Task Distributions:", task_distributions)
print("Critical Path Distributions:", critical_path_distributions)
    