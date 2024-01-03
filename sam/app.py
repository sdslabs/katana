from flask import Flask

app = Flask(__name__)


@app.route("/")
def hello_world():
    return "Hello, World!!"


@app.route("/path1")
def path1():
    return "You have reached path1"


@app.route("/path2")
def path2():
    return "You have reached path2"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
