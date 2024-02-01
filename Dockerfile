FROM golang:1.19

#Create directory for funwith channels
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
COPY *.go ./

# Build
RUN GOOS=linux go build -o /finnhubb

# Run
CMD ["/finnhubb"]