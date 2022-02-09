FROM golang:alpine
 
RUN mkdir /app
 
COPY . /app
 
WORKDIR /app

RUN go mod init main

RUN go mod tidy

RUN go build -o backend *.go 

# RUn go run . 
 
CMD [ "/app/backend" ]