FROM ubuntu:20.04

WORKDIR /app
ADD . /app

RUN apt-get update && apt-get install -y python3 python3-pip curl
RUN pip3 install flask

EXPOSE 8080
CMD ["python3","app.py"]