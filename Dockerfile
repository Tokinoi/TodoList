# Étape 1 : build avec Go sans go.mod
FROM golang:1.21-alpine AS builder

    WORKDIR /app

    # Copier le fichier source
    COPY server.go .

    # Construire sans modules
    RUN go env -w GO111MODULE=off && go build -o app


# Étape 2 : image minimale
FROM scratch

    COPY --from=builder /app/app /app
    ENTRYPOINT ["/app"]
