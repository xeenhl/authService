FROM golang:1.8

WORKDIR /go/src/github.com/xeenhl/spendlog/backend/authService
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure --vendor-only
RUN go build
RUN chmod 777 server 

EXPOSE 8081

CMD ./authService
# ENTRYPOINT  ["chmod", "+x", "./go/src/github.com/xeenhl/spendlog/backend/authService/authService"]