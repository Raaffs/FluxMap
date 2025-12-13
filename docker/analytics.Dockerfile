FROM python:3.13-slim

WORKDIR /app

COPY analytics/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY analytics/ .

EXPOSE 5000
CMD ["python", "main.py"]
