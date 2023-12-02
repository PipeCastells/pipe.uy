# Use the official Golang image based on Alpine
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application files into the container
COPY main.go .
COPY go.mod .
COPY go.sum .

# Create directories for the application
RUN mkdir public html projects

# Copy the contents of the local public, html, and projects directories into the respective directories in the container
COPY public/ ./public
COPY html/ ./html
COPY projects/ ./projects

# Build the Go application
RUN go build -o main .

# Expose the port that the application will run on (adjust as needed)
EXPOSE 8080

# Command to run the application
CMD ["./main"]