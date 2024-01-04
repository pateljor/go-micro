# build a tiny docker image
FROM --platform=linux/amd64 alpine:latest

RUN mkdir /app

COPY frontEndApp /app

CMD [ "/app/frontEndApp" ]