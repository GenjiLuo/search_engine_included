FROM golang:1.12 as builder

ARG GOPROXY
ENV GORPOXY ${GOPROXY}

ADD . /builder

WORKDIR /builder

RUN go build main.go && go build api.go

FROM golang:1.12

COPY --from=builder /builder/main /app/search-engine-include

COPY --from=builder /builder/api /app/search-engine-include-api

WORKDIR /app
