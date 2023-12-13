FROM --platform=linux/amd64 public.ecr.aws/lambda/python:3.9

RUN rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm \
    && yum update -y \
    && yum install -y \
    tesseract

COPY requirements.txt  .
RUN  pip3 install -r requirements.txt --target "${LAMBDA_TASK_ROOT}"

COPY app ${LAMBDA_TASK_ROOT}/

# RUN npm install

CMD ["app.__main__"]
