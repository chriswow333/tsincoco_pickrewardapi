#Stage 1 - Install dependencies and build
FROM --platform=linux/amd64 golang:1.20.4-alpine as builder

WORKDIR /app

ENV DOCKER_DEFAULT_PLATFORM=linux/amd64

# COPY . ./
COPY . .

RUN go mod download

ARG application

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apiv1 cmd/${application}/main.go

# Stage 2 - Create the run-time image
FROM --platform=linux/amd64 scratch

ENV DOCKER_DEFAULT_PLATFORM=linux/amd64

ENV GIN_MODE=release

WORKDIR /server

COPY --from=builder /app/apiv1 ./

ARG ENV
COPY .env.${ENV} ./
ENV ENV_FILE=.env.${ENV}

# COPY your SSL certificate and key to a known location
COPY ./script/assets/cert/fullchain.pem ./
COPY ./script/assets/cert/privkey.pem ./

EXPOSE 50055

CMD ["./apiv1"]


# docker build -t apiv1 .

# docker build --build-arg ENV=prod -t apiv1 .

# docker build --build-arg ENV=prod --build-arg application=app  -t apiv1 .

# docker tag apiv1 asia-east1-docker.pkg.dev/pickreward/pickrewardapi:apiv1 

# docker push asia-east1-docker.pkg.dev/pickreward/pickrewardapi/apiv1

# gcloud auth login


###### IN GCP ######
# sudo gcloud auth configure-docker asia-east1-docker.pkg.dev

