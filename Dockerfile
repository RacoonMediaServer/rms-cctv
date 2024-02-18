FROM golang as builder
WORKDIR /src/service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-cctv -a -installsuffix cgo rms-cctv.go
RUN CGO_ENABLED=0 GOOS=linux go build -o registrator -a -installsuffix cgo ./app/registrator/registrator.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/service/rms-cctv .
COPY --from=builder /src/service/registrator .
COPY --from=builder /src/service/configs/rms-cctv.json /etc/rms/
CMD ["./rms-cctv"]