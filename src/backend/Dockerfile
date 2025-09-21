FROM python:3.10

WORKDIR /app

COPY requirements.txt .

RUN pip install -r requirements.txt

COPY backend.py .
COPY utils/ utils/

CMD ["python", "-u", "backend.py"]