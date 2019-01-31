FROM iron/go

WORKDIR /app

# Now just add the binary
ADD ./buid/ix-ad-service /app/
ADD ./config.yml /app/

ENTRYPOINT ["./ix-ad-service"]
