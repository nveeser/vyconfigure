FROM gcr.io/distroless/static-debian11

COPY vysync /

CMD ["/vysync"]
