# Use an official base image
FROM ubuntu:latest

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY scheduler.txt /app

# Make port 80 available to the world outside this container
EXPOSE 80

# # Define environment variable
# ENV NAME World

# Run command when the container launches
CMD ["echo", "container created successfully"]
