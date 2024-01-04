# build a tiny docker image
FROM --platform=linux/amd64 alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp" ]