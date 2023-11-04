FROM python:3.9-slim-buster
WORKDIR /app
COPY . .
RUN pip install --default-timeout=150 -r requirements.txt
CMD ["python", "-u" ,"-m", "mlworker.main"]