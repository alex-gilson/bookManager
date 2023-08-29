# Specify the base image for the go app.
FROM golang:1.21
# Specify that we now need to execute any commands in this directory.
WORKDIR /app
# Copy everything from this project into the filesystem of the container.
COPY . .
RUN go mod init github.com/go-mysql
# Obtain the package needed to run redis commands. Alternatively use GO Modules.
RUN go get -u github.com/go-sql-driver/mysql
RUN go get github.com/gin-gonic/gin
# RUN go get .
# Compile the binary exe for our app.
RUN go build -o main .
# Start the application.
CMD ["./main"]