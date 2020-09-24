FROM alpine
ADD ./service /
ENV X_TOKEN ""
ENV NATS_HOSTS ""
ENV NATS_SUBJECT ""
ENTRYPOINT ["/service"]
