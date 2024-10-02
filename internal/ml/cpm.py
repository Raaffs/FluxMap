import numpy as np
from scipy.stats import norm

# Define the CPM data for tasks

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

            critical_path_mean += mean
            critical_path_variance += std_dev ** 2

    critical_path_std = np.sqrt(critical_path_variance)
    critical_path_distributions.append({
        "taskId": "criticalPathOverall",
        "mean": critical_path_mean,
        "std_dev": critical_path_std
    })

    return critical_path_distributions
