# The base go-image
FROM golang:1.18-alpine

# Install Git
RUN apk add git

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
COPY . /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
RUN go build -o listener .

# Run the server executable
CMD [ "/app/listener" ]