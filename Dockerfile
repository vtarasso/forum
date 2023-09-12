FROM golang:alpine AS builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum main.go
FROM alpine
LABEL project-name="FORUM"
LABEL git-repo="git@git.01.alem.school:smustafi/forum.git"
LABEL authors="smustafi, vtarasso, dkazgozh"
LABEL release-date="15/02/2023"
WORKDIR /app
COPY --from=builder /app .
CMD ["./forum"]
EXPOSE 4000
