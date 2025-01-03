FROM golang:1.23-alpine3.19 AS build

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags="-s -w" -o ./userli-synapse-user-retention


FROM scratch AS runtime

COPY --from=build /app/userli-synapse-user-retention /userli-synapse-user-retention

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER appuser:appuser

ENTRYPOINT ["/userli-synapse-user-retention"]
