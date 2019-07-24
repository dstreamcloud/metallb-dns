FROM alpine:latest
WORKDIR /
ADD ./bin/controller /controller
RUN chmod a+x /controller
EXPOSE 53
ENTRYPOINT ["/controller"]