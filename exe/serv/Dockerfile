FROM golang:alpine as base
RUN git clone https://github.com/devansh42/shree.git
WORKDIR shree/exe/serv
RUN go build .

FROM alpine
RUN mkdir -p /serv/keys
COPY --from=base /go/shree/exe/serv/serv /serv/serv 
WORKDIR /serv
CMD ["./serv"]