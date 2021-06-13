# Compile stage
FROM golang:1.16.5 AS build-env
# Build Delve
#RUN go get github.com/go-delve/delve/cmd/dlv
ADD . /dockerdev
WORKDIR /dockerdev
RUN ls -l
#RUN go build -gcflags="all=-N -l" -o /server cmd/ocp-video-api/main.go
RUN make all

# Final stage
FROM debian:buster
# EXPOSE 40000
EXPOSE 7000 7002
WORKDIR /
#COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /server /
CMD ["/server"]
#CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]
