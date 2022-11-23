# Start from golang base image
FROM golang:alpine as builder

# Add Maintainer info
LABEL maintainer="Afteni Daniela"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN mkdir /app 

ARG config config
ARG portone  
ARG porttwo

COPY . /app  

COPY . /app  
# Replacing the configurations folder files with needed configurations 
COPY ${config} /app/config

WORKDIR /app


RUN export GO111MODULE=on

#RUN go mod tidy 

EXPOSE ${portone} ${porttwo}

CMD go run .