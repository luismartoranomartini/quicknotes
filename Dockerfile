FROM golang:1.27-alpine 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/http/

FROM scratch
WORKDIR /bin
COPY --from=0 /app/server server
# COPY --from=0 /app/views /bin/views
COPY .env .env
CMD ["/bin/server"]
