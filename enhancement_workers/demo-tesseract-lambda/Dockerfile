FROM public.ecr.aws/lambda/nodejs:16

# RUN yum update -y \
#     && yum install -y tesseract-ocr

RUN rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm \
    && yum -y update \
    && yum -y install tesseract

COPY app ${LAMBDA_TASK_ROOT}
COPY package.json ${LAMBDA_TASK_ROOT}

RUN npm install

CMD ["app.handler"]
