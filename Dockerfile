FROM golang:1.11-alpine

MAINTAINER Kingdon Barrett <kingdon@teamhephy.com>

RUN apk --update add ca-certificates

ARG VCS_REF
ARG BUILD_DATE

# Metadata
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/teamhephy/s3-uploader" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.docker.dockerfile="/Dockerfile"

COPY . /go/src/github.com/teamhephy/s3-uploader
ADD test.jpg /upload/test.jpg

ENV GOPATH /go
RUN cd $GOPATH/src/github.com/teamhephy/s3-uploader && go install -v .

CMD ["s3-uploader"]
