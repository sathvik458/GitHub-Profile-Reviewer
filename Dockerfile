# ---------- build ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Manifests first, on their own layer, so dependency resolution stays cached
# when only source files change.
COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 links nothing from libc, producing a fully static binary that
# runs on a base image with no shared libraries at all.
# -s -w strip the symbol table and DWARF data: ~9.0MB down to ~6.1MB.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./server

# ---------- runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot

# distroless/static ships CA certificates, which this service needs for its
# outbound TLS calls to api.github.com. On a `FROM scratch` base they have to be
# copied across from the builder or every GitHub request fails verification with
# "x509: certificate signed by unknown authority".
COPY --from=builder /out/server /server

# Containers run as root by default. uid 65532, already present in the base.
USER nonroot:nonroot

# Documentation only; publishing the port is a `docker run -p` concern.
EXPOSE 8080

ENTRYPOINT ["/server"]
