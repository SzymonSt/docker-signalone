FROM python:3.9-slim-buster
WORKDIR /app
COPY . .
RUN pip install --default-timeout=150 -r requirements.txt
EXPOSE 8000
CMD ["python", "-u", "-m", "api.main"]