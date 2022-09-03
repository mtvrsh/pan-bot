FROM gcr.io/distroless/static-debian11
COPY pan-bot /bin
USER 1001
ENTRYPOINT ["pan-bot"]
