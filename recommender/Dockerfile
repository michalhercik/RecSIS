FROM python:3.10 AS builder

WORKDIR /app

COPY requirements.txt /app
RUN pip3 install -r requirements.txt

COPY *.py /app

ENTRYPOINT ["python3"]
CMD ["server.py"]
