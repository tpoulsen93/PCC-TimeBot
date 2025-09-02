# Use the official Go image as base
FROM mcr.microsoft.com/vscode/devcontainers/go:1.24-bullseye

# Install additional OS packages for development
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    postgresql-client \
    curl \
    wget \
    git \
    vim \
    nano \
    htop \
    tree \
    jq \
    make \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install Go tools for development
RUN go install github.com/air-verse/air@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest

# Set the working directory
WORKDIR /workspace

# Create a non-root user for development
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Configure Go environment for the user
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Create directories with proper permissions
RUN mkdir -p $GOPATH/src $GOPATH/bin && chmod -R 777 $GOPATH

# Copy go mod files to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Set up shell for development
RUN echo 'alias ll="ls -la"' >> /home/$USERNAME/.bashrc \
    && echo 'alias la="ls -la"' >> /home/$USERNAME/.bashrc \
    && echo 'alias ..="cd .."' >> /home/$USERNAME/.bashrc \
    && echo 'export PS1="\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ "' >> /home/$USERNAME/.bashrc

# Switch to non-root user
USER $USERNAME
