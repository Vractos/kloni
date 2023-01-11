# Build
FROM golang:1.18 as builder

WORKDIR /usr/src/app/

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o ./bin/dolly

# Final Image
FROM alpine

WORKDIR /usr/src/app/

ENV APP_ENV=development

# Auth0 #
ARG auth0_domain
ENV AUTH0_DOMAIN ${auth0_domain}

ARG auth0_audience 
ENV AUTH0_AUDIENCE ${auth0_audience}

# Postgres #
ARG postgres_user
ENV POSTGRES_USER ${postgres_user}

ARG postgres_password
ENV POSTGRES_PASSWORD ${postgres_password}

ARG postgres_db_name
ENV POSTGRES_DB_NAME ${postgres_db_name}

ARG postgres_host
ENV POSTGRES_HOST ${postgres_host}

# Redis #
ARG redis_host
ENV REDIS_HOST ${redis_host}

# Meli #
ARG meli_app_id
ENV MELI_APP_ID ${meli_app_id}

ARG meli_secret_key
ENV MELI_SECRET_KEY ${meli_secret_key}

ARG meli_redirect_url
ENV MELI_REDIRECT_URL ${meli_redirect_url}

ARG meli_endpoint
ENV MELI_ENDPOINT ${meli_endpoint}

# SQS #
ARG ORDER_QUEUE_URL
ENV ORDER_QUEUE_URL ${meli_order_url}

COPY --from=builder /usr/src/app/bin/dolly ./

EXPOSE 8080
ENTRYPOINT [ "./dolly" ]

