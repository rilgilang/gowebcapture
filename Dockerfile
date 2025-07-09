# STEP 1: Build the application
FROM --platform=linux/amd64 golang:1.23-bullseye AS builder

ARG GO_BUILD_COMMAND="go build -tags static_all -o webcapture"

# Install Chrome dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libc6 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libexpat1 \
    libfontconfig1 \
    libgbm1 \
    libgcc1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libstdc++6 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    lsb-release \
    xdg-utils

# Install Brave Browser (AMD64)
RUN apt-get install -y curl && \
    curl -fsSLo /usr/share/keyrings/brave-browser-archive-keyring.gpg https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/brave-browser-archive-keyring.gpg] https://brave-browser-apt-release.s3.brave.com/ stable main" > /etc/apt/sources.list.d/brave-browser-release.list && \
    apt-get update && apt-get install -y brave-browser

# Set environment variables
ENV DISPLAY=:99

# Set Go project path
WORKDIR /go/src/github.com/rilgilang/gowebcapture
COPY . .

# Download and build the Go app
RUN go mod tidy && go mod download
RUN eval $GO_BUILD_COMMAND

# STEP 2: Runtime image
FROM --platform=linux/amd64 debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    procps \
    xvfb \
    imagemagick \
    x11-apps \
    fonts-freefont-ttf \
    fonts-noto \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libc6 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libexpat1 \
    libfontconfig1 \
    libgbm1 \
    libgcc1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libstdc++6 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    xdg-utils \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Timezone configuration
RUN ln -snf /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone

# To make sur xvfb keep run when image is booting/reload
RUN Xvfb -ac :99 -screen 0 360x640x24 &

# Set env vars
ENV DISPLAY=:99
ENV PROJECT_DIR=/go/src/github.com/rilgilang/gowebcapture

# Create app directory
WORKDIR $PROJECT_DIR

# Copy app binary from builder
COPY --from=builder /go/src/github.com/rilgilang/gowebcapture/webcapture .

# Copy Brave from builder stage
COPY --from=builder /usr/bin/brave-browser /usr/bin/brave-browser
COPY --from=builder /opt/brave.com /opt/brave.com

# Make binary executable
RUN chmod +x webcapture

# Launch Xvfb and run your app
CMD ["sh", "-c", "Xvfb :99 -screen 0 360x640x24 -nocursor & ./webcapture"]