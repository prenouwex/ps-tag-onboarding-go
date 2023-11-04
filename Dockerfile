FROM golang:1.21.3-bullseye AS builder

ENV GOARCH=arm64
ENV CGO_ENABLED=1

WORKDIR /usr/local/ps-tag-onboarding-go/
COPY . .

RUN go build -o /go/bin/ps-tag-onboarding-go ./cmd/ps-tag-onboarding-go

FROM scratch

COPY  --from=builder /go/bin/ps-tag-onboarding-go /go/bin/ps-tag-onboarding-go

EXPOSE ${HTTP_PORT}

FROM builder

ENTRYPOINT ["/go/bin/ps-tag-onboarding-go"]