FROM golang:1.15
WORKDIR /app
COPY . ./
ENV CGO_ENABLED=0 
ENV GOOS=linux
RUN go build ./cmd/cli && go build ./cmd/server

FROM alpine
WORKDIR /
COPY --from=0 /app/cli ./cli
COPY --from=0 /app/server ./server
COPY api /api
COPY build/entrypoint.sh /entrypoint.sh
RUN mkdir /data

ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["server"]
