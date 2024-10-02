from flask import Flask, jsonify, request
from cpm import get_task_distributions, get_critical_path_distributions
from flask_cors import CORS

app = Flask(__name__)

cors = CORS(app, resources={r"/api/*": {"origins": "http://localhost:4000"}})

@app.route("/api/pert", methods=['POST'])  # Allow POST requests
def handle_pert():
    data = request.json  # Expect a JSON payload
    if not data:
        return jsonify({'error': 'Invalid data'}), 400  # Return an error if no data is received
    print("data in python: ",data)
    task_distribution = get_task_distributions(data)
    critical_path = get_critical_path_distributions(data)
    print(task_distribution, critical_path)
    
    res = {
        'mean': task_distribution,
        'path': critical_path
    }
    print(res)
    return jsonify(res)

if __name__ == "__main__":
    app.run()
