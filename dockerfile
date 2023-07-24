# Use an official MongoDB base image
FROM mongo:latest

# Expose the default MongoDB port
EXPOSE 27017

# Start MongoDB when the container launches
CMD ["mongod"]
