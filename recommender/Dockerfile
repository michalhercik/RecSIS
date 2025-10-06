FROM python:3.10

WORKDIR /app

COPY requirements.txt /app
RUN pip install --no-cache-dir --upgrade -r requirements.txt

CMD ["fastapi", "dev", "main.py", "--host", "0.0.0.0", "--port", "8002"]


