FROM golang:1.15
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/cli

FROM alpine
WORKDIR /
COPY --from=0 /app/cli ./cli
COPY build/entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["server"]
