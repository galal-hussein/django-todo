FROM python:3.11-slim       

WORKDIR /app                 

COPY . .    

RUN pip install -r requirements.txt  

#CMD ["sh", "+e", "-c", "python manage.py makemigrations && python manage.py migrate && python manage.py runserver 0.0.0.0:8000"]

ENTRYPOINT ["/app/entry.sh"]
