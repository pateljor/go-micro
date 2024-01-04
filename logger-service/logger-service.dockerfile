# build a tiny docker image
FROM --platform=linux/amd64 alpine:latest

RUN mkdir /app

COPY loggerServiceApp /app

CMD [ "/app/loggerServiceApp" ]