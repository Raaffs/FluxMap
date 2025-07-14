from flask import Flask, jsonify, request
from cpm import calculateCpm
from pert import find_critical_path,get_pert_task_distributions,build_graph
from flask_cors import CORS
import json
app = Flask(__name__)

cors = CORS(app, resources={r"/api/*": {"origins": "http://localhost:4000"}})

@app.route("/api/cpm", methods=['POST'])  # Allow POST requests
def handle_cpm():
    data = request.json  # Expect a JSON payload
    if not data:
        return jsonify({'error': 'Invalid data'}), 400

    print("Data received in Python: ", data)

    try:
        result=calculateCpm(data)
        resp={
            'Result':result
        }
        print("Result is: ", resp)

        return jsonify(resp)  # Return the result as JSON
    except Exception as e:
        print("Error processing request: ", e)
        return jsonify({'error': 'An error occurred during processing'}), 500

@app.route("/api/pert",methods=['POST'])
def handle_pert():
    data=request.json
    print("Data received in Python: ", data)

    if not data:
        return jsonify({'error':'invalid data'}),400
    try:
        task_map = {}
        task_distributions = get_pert_task_distributions(data)
        for task_info in task_distributions:
            task_map[task_info['taskId']] = task_info

        # Build graph and calculate critical path
        graph, in_degree = build_graph(task_map)
        critical_path = find_critical_path(graph, in_degree, task_map)
        task_results = [
        {   
            'taskId': info['taskId'],
            'mean': info['mean'],
            'variance': info['variance'],
            'stddev': info['stddev'],
            'predecessorId': info['predecessorId']
        }
        for info in task_distributions
        ]

        res={
            'taskResults': task_results,
            'criticalPath': critical_path
        }
        print("pert result ",res)
        return jsonify(res)
    except Exception as e:
        print("Error calculating pert: ",e.with_traceback)
        return jsonify({'error':'an error occurred during processing'}),500

if __name__ == "__main__":
    app.run()

