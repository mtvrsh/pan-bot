FROM gcr.io/distroless/static-debian11
COPY pan-bot /bin
ENTRYPOINT ["pan-bot"]
