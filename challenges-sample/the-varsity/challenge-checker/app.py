import requests
import time
from flask import Flask, jsonify
import os
from kubernetes import client, config
import logging

import chall_checker  # import the module containing the script

app = Flask(__name__)


chall_name = os.environ.get("CHALL_NAME")
chall_name = chall_name[:-3]
chall_name = "the-varsity"  # fix this later
# kissaki
svc_name = "kissaki-svc"
namespace = "katana"
team_count = 2  # ----------harcoded for now, need to fix this
# Set up logging
logging.basicConfig(level=logging.INFO)

try:
    config.load_incluster_config()
except config.config_exception.ConfigException:
    try:
        config.load_kube_config()
    except config.config_exception.ConfigException:
        raise

# v1 = client.CoreV1Api()
# service = v1.read_namespaced_service(name=svc_name, namespace=namespace)
# # cluster_ip = service.spec.cluster_ip
# ports = service.spec.ports
# port = ports[0].port

v1 = client.CoreV1Api()


def service_port(svc, ns):
    service = v1.read_namespaced_service(name=svc, namespace=ns)
    # cluster_ip = service.spec.cluster_ip
    ports = service.spec.ports
    port = ports[0].port
    return port


@app.route("/")
def hello():
    logging.info(chall_name)
    return "Hello, world!"


@app.route("/test")
def test_challenge_checker():
    res = (
        "making request to "
        + "http://"
        + svc_name
        + "."
        + namespace
        + ".svc.cluster.local:"
        + str(service_port(svc_name, namespace))
        + "/register "
    )
    return res


# @app.route("/register")
# def register_challenge_checker():
#     logging.info(
#         "making request to "
#         + "http://"
#         + str(cluster_ip)
#         + ":"
#         + str(port)
#         + "/register "
#     )

#     # Register with kissaki
#     checker_info = {
#         "name": "knock-challenge-checker",
#         "challenge": "knock",
#     }  # Example info

#     response = requests.post(
#         "http://" + str(cluster_ip) + ":" + str(port) + "/register",
#         json=checker_info,
#     )
#     message = response.json().get("message")

#     logging.info(f"Received message from kissaki: {message}")

#     return "challenge_checker registered in kissaki"


@app.route("/register")
def register_challenge_checker():
    logging.info(
        "making request to "
        + "http://"
        + svc_name
        + "."
        + namespace
        + ".svc.cluster.local:"
        + str(service_port(svc_name, namespace))
        + "/register "
    )

    # Register with kissaki
    # keys in checker_info are harcoded if changed here then some change may be needed in katana-services/Kissaki/src/app.py
    checker_info = {"ccName": chall_name + "-cc"}

    response = requests.post(
        "http://"
        + svc_name
        + "."
        + namespace
        + ".svc.cluster.local:"
        + str(service_port(svc_name, namespace))
        + "/register",
        json=checker_info,
    )
    message = response.json().get("message")

    logging.info(f"Received message from kissaki: {message}")

    return "challenge_checker registered in kissaki"


# {service_port(chall_svc, chall_ns)}
@app.route("/check")
def check_challenge():
    i = 0
    chall_svc = f"{chall_name}-svc-{i}"
    chall_ns = f"katana-team-{i}-ns"
    url = f"http://{chall_svc}.{chall_ns}.svc.cluster.local:80/"
    return url
    try:
        status = chall_checker.check_challenge(url)
        return status
    except Exception as e:
        logging.error(f"Error checking challenge: {str(e)}")
        return str(e)


@app.route("/checker")
def check_route():
    results = {"challengeName": chall_name, "data": []}
    for i in range(team_count):
        team_name = f"katana-team-{i}"
        result = {"team-name": team_name}
        chall_svc = f"{chall_name}-svc-{i}"
        chall_ns = f"katana-team-{i}-ns"
        url = f"http://{chall_svc}.{chall_ns}.svc.cluster.local:80/"  # update this later ---port should not be hardcoded----
        try:
            status = chall_checker.check_challenge(url)
            result["status"] = status
        except Exception as e:
            logging.error(f"Error checking challenge: {str(e)}")
            result["error"] = str(e)
        results["data"].append(result)
    return jsonify(results)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
