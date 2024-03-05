import requests
import time
from flask import Flask, jsonify
import os
from kubernetes import client, config
import logging

# app = Flask(__name__)

# Set up logging
logging.basicConfig(level=logging.INFO)

try:
    config.load_incluster_config()
except config.config_exception.ConfigException:
    try:
        config.load_kube_config()
    except config.config_exception.ConfigException:
        raise

v1 = client.CoreV1Api()
service = v1.read_namespaced_service(name="kissaki-svc", namespace="katana")
cluster_ip = service.spec.cluster_ip
ports = service.spec.ports
port = ports[0].port


# @app.route("/")
def hello():
    return "Hello, world!"


# @app.route("/test")
def test_challenge_checker():
    res = (
        "making request to "
        + "http://"
        + str(cluster_ip)
        + ":"
        + str(port)
        + "/register "
    )
    return res


# @app.route("/register")
def register_challenge_checker():
    logging.info(
        "making request to "
        + "http://"
        + str(cluster_ip)
        + ":"
        + str(port)
        + "/register "
    )

    # Register with kissaki
    checker_info = {
        "name": "knock-challenge-checker",
        "challenge": "knock",
    }  # Example info

    response = requests.post(
        "http://" + str(cluster_ip) + ":" + str(port) + "/register",
        json=checker_info,
    )
    message = response.json().get("message")

    logging.info(f"Received message from kissaki: {message}")

    return "challenge_checker registered in kissaki"


# @app.route("/check")
def check_challenge():
    for i in range(10):
        # TODO: Implement challenge checking logic
        challenge_status = {"status": "OK"}  # Example status

        # Send status to kissaki service
        response = requests.post(
            "http://" + str(cluster_ip) + ":" + str(port) + "/status",
            json=challenge_status,
        )
        message = response.json().get("message")
        logging.info(f"Received message from kissaki: {message}")

        time.sleep(10)  # Check every 10 seconds

    return jsonify(challenge_status)


# if __name__ == "__main__":
#     app.run(host="0.0.0.0", port=8080)
