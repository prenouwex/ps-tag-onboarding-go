FROM golang:1.20-alpine AS builder

#RUN apk update && apk add --no-cache git
#RUN apk update

ENV GOARCH=arm64

RUN apk add --no-cache gcc g++ git openssh-client

WORKDIR /usr/local/ps-tag-onboarding-go/
COPY . .

#RUN go build -o /go/bin/ps-tag-onboarding-go
#RUN ./build/compile.sh

RUN go env -w GOARCH=arm64

RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o /go/bin/ps-tag-onboarding-go ./cmd/ps-tag-onboarding-go

FROM scratch

#COPY --from=builder /go/bin/my-project /go/bin/ps-tag-onboarding-go
COPY  --from=builder /go/bin/ps-tag-onboarding-go /go/bin/ps-tag-onboarding-go

#RUN echo ${HTTP_PORT}

#EXPOSE ${HTTP_PORT}
EXPOSE 8089

FROM builder

ENTRYPOINT ["/go/bin/ps-tag-onboarding-go"]