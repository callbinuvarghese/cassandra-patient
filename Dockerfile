FROM iron/base
WORKDIR /app
COPY patientapp /app/
EXPOSE 8080
ENTRYPOINT ["./patientapp"]
