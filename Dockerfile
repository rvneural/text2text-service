FROM ubuntu:latest
EXPOSE 80
RUN ap-get update && apt-get upgrade
COPY . .
WORKDIR /build/linux
CMD [ "./text2text-serivce" ]

