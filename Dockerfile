FROM golang:1.15.2-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /usr/local/ps-tag-onboarding-go/
COPY . .

RUN go build -o /go/bin/ps-tag-onboarding-go

FROM scratch

COPY --from=builder /go/bin/my-project /go/bin/ps-tag-onboarding-go

EXPOSE ${HTTP_PORT}

FROM builder

ENTRYPOINT ["/go/bin/my-project"]