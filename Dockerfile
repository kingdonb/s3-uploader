FROM golang:1.6-alpine

MAINTAINER Lachlan Evenson <lachlan.evenson@gmail.com>

RUN apk --update add ca-certificates

ARG VCS_REF
ARG BUILD_DATE

# Metadata
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/lachie83/s3-uploader" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.docker.dockerfile="/Dockerfile"

COPY . /go/src/github.com/lachie83/s3-uploader
ADD test.jpg /upload/test.jpg

ENV GOPATH /go
RUN cd $GOPATH/src/github.com/lachie83/s3-uploader && go install -v .

CMD ["s3-uploader"]