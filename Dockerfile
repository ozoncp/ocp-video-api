# Compile stage
FROM golang:1.16.5 AS build-env
ADD . /dockerdev
WORKDIR /dockerdev
RUN make all

# Final stage
FROM debian:buster
EXPOSE 7000 7002 8192
WORKDIR /
COPY --from=build-env /server /
CMD ["/server"]
