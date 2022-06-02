FROM alpine:3

COPY  ./pveops /pveopst

RUN chmod +x /pveops

USER 101

ENTRYPOINT ["/pveops"]
