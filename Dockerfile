from golang:1.12-alpine AS builder

# Need git
RUN apk update && apk add --no-cache git tzdata

# Used for security in image
RUN adduser -D -g '' appuser

WORKDIR /usr/src/myapp
ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/ambrosia .

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/ambrosia /ambrosia

# USER appuser
ENTRYPOINT ["/ambrosia"]
