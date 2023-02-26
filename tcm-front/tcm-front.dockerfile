FROM alpine:latest

RUN mkdir /app

COPY ./templates /templates

COPY tcm-frontApp /app

CMD [ "/app/tcm-frontApp"]