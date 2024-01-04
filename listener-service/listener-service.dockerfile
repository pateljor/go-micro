# build a tiny docker image
FROM --platform=linux/amd64 alpine:latest

RUN mkdir /app

COPY listenerApp /app

CMD [ "/app/listenerApp" ]