FROM golang AS base
RUN  git clone https://github.com/devansh42/shree.git
WORKDIR ~/shree/back
RUN go build .

FROM alpine
RUN mkdir -p /backend
COPY --from=base ~/shree/back/back /backend/back 
WORKDIR /backend
CMD ["./back"]

