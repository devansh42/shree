FROM golang AS base
RUN  git clone https://github.com/devansh42/shree.git
WORKDIR /go/shree/back
RUN go build .

FROM alpine
RUN mkdir -p /back/keys
COPY --from=base /go/shree/back/back /back/back 
WORKDIR /back
CMD ["./back"]

