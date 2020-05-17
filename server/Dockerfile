FROM ubuntu

RUN apt update && apt install -y curl

WORKDIR /usr/src/app
COPY . .

EXPOSE 80
CMD [ "./server" ]