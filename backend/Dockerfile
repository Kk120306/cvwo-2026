FROM golang:1.25.5-alpine

# Sets the working directory inside the container
WORKDIR /app

# Copies go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Runs all dependancies in go.mod and go.sum
RUN go mod tidy 

# Copies the rest of the application code to the working directory
COPY . .

# Build go app
RUN go build -o main ./main.go

# Kill binary permission
RUN chmod +x main

# Command to run the executable
EXPOSE 4040

# Final command to run 
CMD ["./main"]