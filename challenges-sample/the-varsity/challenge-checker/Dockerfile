FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y python3 python3-pip curl

WORKDIR /opt/kissaki/
COPY . .

RUN pip3 install -r requirements.txt

# this arg is not being passed need to fix that.------------------------
ARG chall_name
ENV CHALL_NAME=$chall_name

CMD ["/bin/bash", "-c", "python3 /opt/kissaki/app.py" ]
