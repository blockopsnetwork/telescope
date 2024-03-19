# Use a base image 
FROM alpine:latest

RUN apk add --no-cache curl  bash 

# Copy the script into the container
COPY install.sh /

# Set the script as executable
RUN chmod +x /install.sh

# Define the entrypoint for the container
ENTRYPOINT ["/install.sh"]