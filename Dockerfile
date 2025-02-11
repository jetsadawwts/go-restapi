FROM golang:1.20-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /bin/app


FROM gcr.io/distroless/static-debian11

COPY --from=build /bin/app /bin
COPY .env.prod /bin/.env.prod


EXPOSE 3000

ENTRYPOINT [ "/bin/app", "/bin/.env.prod" ]

