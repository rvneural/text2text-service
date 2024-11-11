FROM debian:latest
LABEL maintainer="gafarov@realnoevremya.ru"
RUN apt-get update && apt-get upgrade
EXPOSE 80
COPY . .
WORKDIR /build/linux
CMD [ "./text2text-service" ]

