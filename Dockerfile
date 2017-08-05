FROM alpine:3.4
MAINTAINER Fabien Foerster <fabienfoerster@gmail.com>
RUN apk add --no-cache ca-certificates
ADD autosub /usr/bin/autosub
CMD autosub /autosub/data