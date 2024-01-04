# build a tiny docker image
FROM --platform=linux/amd64 alpine:latest

RUN mkdir /app

COPY mailServiceApp /app
COPY templates /templates

CMD [ "/app/mailServiceApp" ]