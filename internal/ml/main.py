from flask import Flask,jsonify,request
from cpm import get_task_distributions, critical_path_distributions
from flask_cors import CORS
app=Flask(__name__)

cors = CORS(app, resources={r"/api/*": {"origins": "http://localhost:4000"}})
@app.route("/api/pert",methods=['POST'])
def handle_pert():
    data=request.json
    var,mean,total_mean,total_var=get_task_distributions(data['optimisticTimes'],data['pessimisticTimes'],data['mostLikelyTimes'])
    res={
        'var':var,
        'mean':mean,
        'total_mean':total_mean,
        'total_var':total_var
    }
    return jsonify(res)

if __name__=="__main__":
    app.run()
