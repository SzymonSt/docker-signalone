FROM golang:1.20-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /agent
COPY agent/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY agent/. .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

FROM --platform=$BUILDPLATFORM node:18.12-alpine3.16 AS client-builder
WORKDIR /ui
# cache packages in layer
COPY ui/package.json /ui/package.json
COPY ui/package-lock.json /ui/package-lock.json
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY ui /ui
RUN npm cache clean --force && \
npm run build

FROM alpine
LABEL org.opencontainers.image.title="Signal0ne" \
    org.opencontainers.image.description="Painless AI powered debugging for containerized apps." \
    org.opencontainers.image.vendor="signalone" \
    com.docker.desktop.extension.api.version="0.3.4" \
    com.docker.extension.screenshots='[{"alt":"Main Screen","url":"https://signaloneai.com/online-assets/screen2.PNG"},{"alt":"Login Screen", "url": "https://signaloneai.com/online-assets/screen3.png"}]' \
    com.docker.desktop.extension.icon="https://signaloneai.com/online-assets/Signal0ne_min.jpg" \
    com.docker.extension.publisher-url="https://signaloneai.com/" \
    com.docker.extension.detailed-description="The magic of the Signal0ne Docker Extension is its very short feedback loop for working with and developing containerized applications. Register new account and forget endless log diving – the extension offers a clear and concise summary of what's happening inside your containers, pinpointing potential issues and even suggesting solutions." \
    com.docker.extension.additional-urls='[{"title":"Demo Video","url":"https://youtu.be/LWOUSr9s_Zw"},{"title":"Discord","url":"https://discord.gg/vAZrxKs5f6"},{"title":"Linkedin","url":"https://www.linkedin.com/company/signal0ne/"}]' \
    com.docker.extension.categories="utility-tools" \
    com.docker.extension.changelog="Added new signup type - signal0ne account" \
COPY --from=builder /agent/bin/service /
COPY docker-compose.yaml .
COPY metadata.json .
COPY metadata.json /metadata.json
COPY s0.svg .
COPY agent/.default.env .
COPY --from=client-builder /ui/build ui
CMD /service
