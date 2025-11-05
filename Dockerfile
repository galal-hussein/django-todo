FROM python:3.11-slim       

ENV DJANGO_SUPERUSER_USERNAME root
ENV DJANGO_SUPERUSER_EMAIL root@gmail.com
ENV DJANGO_SUPERUSER_PASSWORD root

WORKDIR /app                 

COPY . .    

RUN pip install -r requirements.txt  

#CMD ["sh", "+e", "-c", "python manage.py makemigrations && python manage.py migrate && python manage.py createsuperuser --noinput || true && python manage.py runserver 0.0.0.0:8000"]

ENTRYPOINT ["/app/entry.sh"]
