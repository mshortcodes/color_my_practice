FROM debian:stable-slim

RUN curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.14.2/cloud-sql-proxy.linux.amd64 && \
    chmod +x cloud-sql-proxy

COPY color_my_practice /bin/color_my_practice

ENTRYPOINT [ "/bin/sh", "-c", "cloud-sql-proxy $INSTANCE_CONN_NAME & exec /bin/color_my_practice" ]