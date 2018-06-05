FROM iron/base
WORKDIR /app
COPY patientapp /app/
ENTRYPOINT ["./patientapp"]
