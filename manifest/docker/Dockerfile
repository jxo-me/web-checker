FROM loads/alpine:3.8

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
ADD config.yaml             $WORKDIR/config.yaml
ADD checker $WORKDIR/checker
RUN chmod +x $WORKDIR/checker

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./checker
