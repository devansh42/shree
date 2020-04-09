FROM golang as base
RUN git clone https://github.com/devansh42/shree.git
WORKDIR shree/ca
RUN go build .

FROM alpine
RUN apk update && apk add libc6-compat
RUN mkdir -p /ca/keys 
COPY --from=base /go/shree/ca/ca /ca/ca 
WORKDIR /ca
CMD ["/ca/ca"]