from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"

if __name__ == '__main__':
    # Run on port 5000, accessible from outside the container
    app.run(host='0.0.0.0', port=5000)