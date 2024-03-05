#!/usr/bin/python3
from flask import Flask, request
import logging

app = Flask(__name__)

# Set up logging
logging.basicConfig(level=logging.INFO)


@app.route("/", methods=["POST"])
def handle_post():
    data = request.get_json()  # Get JSON data from the request
    logging.info(f"Received POST request with data: {data}")  # Log the data
    return "POST request received!", 200


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
