FROM --platform=linux/amd64 golang:1.23

# Set noninteractive mode
ENV DEBIAN_FRONTEND=noninteractive

# Install base dependencies
RUN apt-get update && apt-get install -y \
    wget \
    xvfb \
    fluxbox \
    x11vnc \
    dbus-x11 \
    ffmpeg \
    libx11-dev \
    libxtst6 \
    libxrandr2 \
    libgtk-3-0 \
    libnss3 \
    libasound2 \
    ca-certificates \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

# Install Brave Browser (AMD64)
RUN apt-get install -y curl && \
    curl -fsSLo /usr/share/keyrings/brave-browser-archive-keyring.gpg https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/brave-browser-archive-keyring.gpg] https://brave-browser-apt-release.s3.brave.com/ stable main" > /etc/apt/sources.list.d/brave-browser-release.list && \
    apt-get update && apt-get install -y brave-browser

WORKDIR /app
COPY . .

RUN go build -o app

# Environment variables
ENV DISPLAY=:99
ENV DBUS_SESSION_BUS_ADDRESS=/dev/null

EXPOSE 5900

CMD bash -c "\
    rm -f /tmp/.X99-lock /tmp/.X11-unix/X99 && \
    dbus-daemon --system --fork && \
    Xvfb :99 -screen 0 1024x768x16 -ac & \
    sleep 2 && \
    fluxbox & \
    x11vnc -display :99 -nopw -listen 0.0.0.0 -forever & \
    sleep 2 && \
    ./app"