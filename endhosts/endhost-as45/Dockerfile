FROM scion-base:1.0

# Copy endhost topology (AS 19-ffaa:1:45; relies on scion45's AS services)
COPY topology.json /etc/scion/topology.json

# Generate master keys
RUN head -c 16 /dev/urandom | base64 > /etc/scion/keys/master0.key && \
    head -c 16 /dev/urandom | base64 > /etc/scion/keys/master1.key

# Disable services, not needed on endhosts
RUN systemctl disable scion-control@cs.service && \
    systemctl disable scion-router@br.service && \
    systemctl disable scion-node-manager.service

# Only run daemon and dispatcher
RUN systemctl enable scion-daemon.service && \
    systemctl enable scion-dispatcher.service

# SCION application tools (bandwidth tester) for the experiments
RUN apt-get install -y apt-transport-https && \
    echo "deb [trusted=yes] https://packages.netsec.inf.ethz.ch/debian all main" \
    > /etc/apt/sources.list.d/scionlab.list && \
    apt-get update && \
    apt-get install -y scion-apps-bwtester

CMD ["/lib/systemd/systemd", "--unit=multi-user.target"]
