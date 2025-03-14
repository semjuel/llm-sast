FROM openjdk:17-jdk-slim

# Set environment variables for Android SDK
ENV ANDROID_SDK_ROOT=/opt/android-sdk
ENV PATH=$PATH:$ANDROID_SDK_ROOT/cmdline-tools/latest/bin:$ANDROID_SDK_ROOT/build-tools/latest

# Set environment variables for APKEditor
ENV APKEDITOR_HOME=/opt/APKEditor
ENV PATH=$PATH:$APKEDITOR_HOME/bin

# Install necessary packages and enable i386 architecture
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y wget unzip \
    libstdc++6-i386-cross  && \
    rm -rf /var/lib/apt/lists/*

# Create directories for Android SDK and APKEditor
RUN mkdir -p $ANDROID_SDK_ROOT/cmdline-tools && \
    mkdir -p $ANDROID_SDK_ROOT/build-tools && \
    mkdir -p $APKEDITOR_HOME/bin

# Download and install Android SDK Command-line Tools
RUN wget https://dl.google.com/android/repository/commandlinetools-linux-8512546_latest.zip -O /tmp/cmdline-tools.zip && \
    unzip /tmp/cmdline-tools.zip -d /opt/android-sdk/cmdline-tools && \
    mv /opt/android-sdk/cmdline-tools/cmdline-tools /opt/android-sdk/cmdline-tools/latest && \
    rm /tmp/cmdline-tools.zip

# Accept licenses and install build-tools and platform-tools
RUN yes | sdkmanager --licenses && \
    sdkmanager "build-tools;33.0.0" "platform-tools"

# Create a symlink for build-tools
RUN ln -s $ANDROID_SDK_ROOT/build-tools/33.0.0 $ANDROID_SDK_ROOT/build-tools/latest

# Download APKEditor.jar file and prepare bin file
RUN wget https://github.com/REAndroid/APKEditor/releases/download/V1.4.2/APKEditor-1.4.2.jar  \
     -O $APKEDITOR_HOME/bin/apkeditor.jar && \
     echo '#!/bin/bash' > $APKEDITOR_HOME/bin/apkeditor && \
     echo 'java -jar /opt/APKEditor/bin/apkeditor.jar "$@"' >> $APKEDITOR_HOME/bin/apkeditor && \
     chmod +x $APKEDITOR_HOME/bin/apkeditor

ARG JADX_VERSION=1.4.7
RUN wget https://github.com/skylot/jadx/releases/download/v${JADX_VERSION}/jadx-${JADX_VERSION}.zip -O /tmp/jadx.zip && \
    unzip /tmp/jadx.zip -d /opt/jadx && \
    ln -s /opt/jadx/bin/jadx /usr/local/bin/jadx && \
    ln -s /opt/jadx/bin/jadx-gui /usr/local/bin/jadx-gui && \
    rm /tmp/jadx.zip

# Install dependencies and build tools (including a C compiler)
RUN apt-get update && \
    apt-get install -y wget build-essential && \
    rm -rf /var/lib/apt/lists/*

# Install a recent version of Go (e.g., 1.23.6)
ENV GO_VERSION=1.23.6
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"

# Enable CGO to allow compilation of C source files
ENV CGO_ENABLED=0

# Set working directory
WORKDIR /app

COPY ./keystore /keystore

# Copy all files from the current directory into the container
COPY . .

# On container start, run the Go application
CMD ["go", "run", "cmd/api/main.go"]
