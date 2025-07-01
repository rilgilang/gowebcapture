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

# Install Google Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list \
    && apt-get update \
    && apt-get install -y google-chrome-stable \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables
ENV DISPLAY=:99
ENV CHROME_PATH=/usr/bin/google-chrome

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
    xvfb \
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

# Set env vars
ENV DISPLAY=:99
ENV CHROME_PATH=/usr/bin/google-chrome
ENV PROJECT_DIR=/go/src/github.com/rilgilang/gowebcapture

# Create app directory
WORKDIR $PROJECT_DIR

# Copy app binary from builder
COPY --from=builder /go/src/github.com/rilgilang/gowebcapture/gowebcapture .

# Copy Chrome from builder stage
COPY --from=builder /usr/bin/google-chrome /usr/bin/google-chrome
COPY --from=builder /usr/share/man/man1/google-chrome.1.gz /usr/share/man/man1/
COPY --from=builder /opt/google/chrome /opt/google/chrome

# Make binary executable
RUN chmod +x gowebcapture

# Launch Xvfb and run your app
CMD ["sh", "-c", "Xvfb :99 -screen 0 1280x720x24 & ./gowebcapture"]