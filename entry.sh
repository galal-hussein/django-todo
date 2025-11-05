#!/bin/sh

set -x
set +e

python manage.py makemigrations

python manage.py migrate

set -e

python manage.py runserver 0.0.0.0:8000
