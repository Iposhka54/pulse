FROM golang:1.26.10 AS build
LABEL authors="Iposhka54"
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/pulse \
    .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/pulse /pulse

EXPOSE 9100

ENTRYPOINT ["/pulse"]
CMD ["serve"]