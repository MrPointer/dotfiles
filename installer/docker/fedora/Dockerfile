# Fedora-based Docker image with all prerequisites for installer testing
FROM fedora:latest

# Set environment variables to avoid interactive prompts during package installation
ENV TZ=UTC

# Update package list and install all required prerequisites
# Note: Fedora uses dnf and different package names than Debian/Ubuntu
RUN dnf update -y && dnf install -y \
    @development-tools \
    procps-ng \
    curl \
    file \
    git \
    ca-certificates \
    sudo \
    && dnf clean all

# Create a non-root user for testing
RUN useradd -m -s /bin/bash testuser && \
    usermod -aG wheel testuser && \
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
