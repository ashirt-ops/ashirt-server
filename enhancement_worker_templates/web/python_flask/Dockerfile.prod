FROM python:3.10-alpine AS builder

WORKDIR /build
RUN pip install --user pipenv

COPY Pipfile.lock Pipfile ./

RUN /root/.local/bin/pipenv requirements > requirements.txt
###

FROM python:3.10-alpine AS runner

WORKDIR /app
COPY --from=builder /build/requirements.txt .

###################################
# Install other dependencies here #
###################################

RUN pip install -r  requirements.txt

# COPY bin/docker_start.sh ./start.sh
COPY src .

EXPOSE 80

# Run as Alpine's guest user
USER 405

# some guidance on using gunicorn in containers:
# https://pythonspeed.com/articles/gunicorn-in-docker/
CMD ["gunicorn", "--worker-tmp-dir", "/dev/shm", \
     "--workers=2", "--threads=4", "--worker-class=gthread", \
     "--log-file=-", \
     "-b", "0.0.0.0:80", "wsgi:app"]

