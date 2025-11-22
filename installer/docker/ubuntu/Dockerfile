# Ubuntu-based Docker image with all prerequisites for installer testing
FROM ubuntu:latest

# Set environment variables to avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

# Update package list and install all required prerequisites
RUN apt-get update && apt-get install -y \
    build-essential \
    procps \
    curl \
    file \
    git \
    ca-certificates \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user for testing
RUN useradd -m -s /bin/bash testuser && \
    usermod -aG sudo testuser && \
    echo "testuser ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/testuser && \
    chmod 0440 /etc/sudoers.d/testuser

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set working directory
WORKDIR /workspace

# Set default shell to bash
SHELL ["/bin/bash", "-c"]

# Add helpful aliases and environment setup
RUN echo 'alias ll="ls -la"' >> ~/.bashrc && \
    echo 'alias la="ls -A"' >> ~/.bashrc && \
    echo 'export PS1="\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ "' >> ~/.bashrc

# Use entrypoint script to handle permission fixes at runtime
ENTRYPOINT ["/entrypoint.sh"]

# Default command
CMD ["/bin/bash"]
