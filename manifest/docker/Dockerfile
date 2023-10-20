FROM alpine:3.18

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
ENV CONFIG                  $WORKDIR/config.yaml
ADD config.yaml             $WORKDIR/config.yaml
ADD checker $WORKDIR/checker
RUN chmod +x $WORKDIR/checker && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./checker
