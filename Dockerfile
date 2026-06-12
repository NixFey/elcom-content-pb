FROM golang:1.26 AS build
WORKDIR /src

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .
RUN go build -ldflags="-extldflags=-static" -o /out/cms .

FROM alpine AS runtime
WORKDIR /app
EXPOSE 8090
COPY --from=build /out/cms ./cms
ENTRYPOINT ["./cms", "serve", "--http=0.0.0.0:8090"]
