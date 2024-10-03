import numpy as np

# Function to compute PERT distribution for a single task
def compute_pert_distribution(task):
    optimistic = task["optimistic"]
    most_likely = task["mostLikely"]
    pessimistic = task["pessimistic"]

    # Calculate the expected duration (mean) using the PERT formula
    mean = (optimistic + 4 * most_likely + pessimistic) / 6

    # Calculate the standard deviation (using the PERT formula)
    std_dev = (pessimistic - optimistic) / 6

    return mean, std_dev

# Function to compute Z-value for a task
def compute_z_value(observed_duration, mean, std_dev):
    if std_dev == 0:  # To avoid division by zero
        return 0
    return (observed_duration - mean) / std_dev

# Function to process PERT tasks, including predecessor relationships and Z-values
def get_pert_distributions(tasks):
    task_map = {task["parentTaskId"]: task for task in tasks}
    task_distributions = []
    overall_mean = 0
    overall_variance = 0

    for task in tasks:
        mean, std_dev = compute_pert_distribution(task)

        # If there is a predecessor task, add its mean to the current task
        if task["predecessorTaskId"] in task_map:
            predecessor_task = task_map[task["predecessorTaskId"]]
            predecessor_mean, _ = compute_pert_distribution(predecessor_task)
            mean += predecessor_mean  # Add predecessor's duration to the current task

        # Compute Z-value (you can use the most likely value as observed duration)
        observed_duration = task["mostLikely"]
        z_value = compute_z_value(observed_duration, mean, std_dev)

        task_distributions.append({
            "taskId": task["parentTaskId"],
            "mean": mean,
            "std_dev": std_dev,
            "z_value": z_value
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

# Function to return critical path distributions (if needed for PERT analysis)
def get_critical_path_distributions(tasks):
    task_map = {task["parentTaskId"]: task for task in tasks}
    critical_path_distributions = []
    critical_path_mean = 0
    critical_path_variance = 0

    for task in tasks:
        if task.get("criticalPath", False):
            mean, std_dev = compute_pert_distribution(task)

            # Add predecessor task's mean to critical path task if exists
            if task["predecessorTaskId"] in task_map:
                predecessor_task = task_map[task["predecessorTaskId"]]
                predecessor_mean, _ = compute_pert_distribution(predecessor_task)
                mean += predecessor_mean

            # Compute Z-value for critical path task
            observed_duration = task["mostLikely"]
            z_value = compute_z_value(observed_duration, mean, std_dev)

            critical_path_distributions.append({
                "taskId": task["parentTaskId"],
                "mean": mean,
                "std_dev": std_dev,
                "z_value": z_value
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

# Example usage
if __name__ == "__main__":
    tasks = [
        {"parentTaskId": 1, "predecessorTaskId": None, "optimistic": 2, "mostLikely": 4, "pessimistic": 6},
        {"parentTaskId": 2, "predecessorTaskId": 1, "optimistic": 3, "mostLikely": 5, "pessimistic": 8},
        {"parentTaskId": 3, "predecessorTaskId": 2, "optimistic": 1, "mostLikely": 2, "pessimistic": 5, "criticalPath": True},
        {"parentTaskId": 4, "predecessorTaskId": 2, "optimistic": 2, "mostLikely": 4, "pessimistic": 7, "criticalPath": True},
        {"parentTaskId": 5, "predecessorTaskId": 3, "optimistic": 2, "mostLikely": 4, "pessimistic": 7, "criticalPath": False}
    ]

    pert_distributions = get_pert_distributions(tasks)
    critical_path_distributions = get_critical_path_distributions(tasks)

    print("PERT Task Distributions:", pert_distributions)
    print("Critical Path Distributions:", critical_path_distributions)
