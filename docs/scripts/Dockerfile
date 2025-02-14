FROM registry.access.redhat.com/ubi9/ubi@sha256:10fab05b54dfbef5ca8335b150b0a4633abc1c3cfeb28dee5a9f8c6c9792b94e
ARG APPS_REPO
COPY quick_start.sh /quick_start.sh
RUN chmod +x /quick_start.sh
ENV APPS_REPO=${APPS_REPO}
RUN /quick_start.sh

# Optionally, set the container to continue running after execution
CMD ["bash"]

