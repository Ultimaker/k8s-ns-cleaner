FROM alpine:3.7

COPY bin/ns-cleaner /usr/bin/ns-cleaner
RUN chmod +x /usr/bin/ns-cleaner
