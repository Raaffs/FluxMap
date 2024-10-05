import json
import numpy as np
from collections import defaultdict, deque

def calculate_pert_values(task):
    optimistic = task['optimistic']
    pessimistic = task['pessimistic']
    most_likely = task['mostLikely']
    
    expected_time = (optimistic + 4 * most_likely + pessimistic) / 6
    variance = ((pessimistic - optimistic) / 6) ** 2
    stddev = np.sqrt(variance)
    
    return expected_time, variance, stddev

def get_pert_task_distributions(tasks):
    results = []
    for task in tasks:
        expected_time, variance, stddev = calculate_pert_values(task)
        task_info = {
            'taskId': task['parentTaskId'],
            'mean': expected_time,
            'variance': variance,
            'stddev': stddev,
            'predecessorId': task['predecessorTaskId']
        }
        results.append(task_info)
    return results

def build_graph(task_map):
    graph = defaultdict(list)
    in_degree = defaultdict(int)

    for task in task_map.values():
        if task['predecessorId'] is not None:
            graph[task['predecessorId']].append(task['taskId'])
            in_degree[task['taskId']] += 1
            
    return graph, in_degree

def find_critical_path(graph, in_degree, task_map):
    queue = deque()
    earliest_finish = {task['taskId']: 0 for task in task_map.values()}

    for task_id in task_map.keys():
        if in_degree[task_id] == 0:
            queue.append(task_id)

    while queue:
        current = queue.popleft()
        current_finish = earliest_finish[current] + task_map[current]['mean']
        
        for neighbor in graph[current]:
            in_degree[neighbor] -= 1
            earliest_finish[neighbor] = max(earliest_finish[neighbor], current_finish)
            if in_degree[neighbor] == 0:
                queue.append(neighbor)

    max_finish_time = max(earliest_finish.values())
    
    critical_path = []
    for task_id in task_map.keys():
        if earliest_finish[task_id] == max_finish_time:
            critical_path.append(task_id)
            max_finish_time -= task_map[task_id]['mean']

    if critical_path:
        critical_path.reverse()
        parent_map = {task['taskId']: task['predecessorId'] for task in task_map.values() if task['predecessorId'] is not None}

        complete_path = []
        current_task = critical_path[0]
        while current_task is not None:
            complete_path.append(current_task)
            current_task = parent_map.get(current_task)

        return complete_path[::-1]

    return []

def main(json_data):
    tasks = json.loads(json_data)
    
    # Get task distributions
    task_map = {}
    task_distributions = get_pert_task_distributions(tasks)
    for task_info in task_distributions:
        task_map[task_info['taskId']] = task_info
    
    # Build graph and calculate critical path
    graph, in_degree = build_graph(task_map)
    critical_path = find_critical_path(graph, in_degree, task_map)

    # Print results
    print("Task Results:", [
        {
            'taskId': info['taskId'],
            'mean': info['mean'],
            'variance': info['variance'],
            'stddev': info['stddev'],
            'predecessorId': info['predecessorId']
        }
        for info in task_distributions
    ])
    print("Critical Path:", critical_path)
