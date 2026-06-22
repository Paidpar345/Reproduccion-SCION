FROM scion-base:1.0

# Copy endhost topology
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

RUN apt-get install -y apt-transport-https && \
    echo "deb [trusted=yes] https://packages.netsec.inf.ethz.ch/debian all main" \
    > /etc/apt/sources.list.d/scionlab.list && \
    apt-get update && \
    apt-get install -y scion-apps-bwtester

# --- CyberFlood IP forwarding and NAT (OPT-IN: flood/attack scenarios only) ---
# WARNING: this path uses plain IP NAT/MASQUERADE onto the Docker underlay. It does
# NOT traverse the SCION fabric and CONFLICTS with the SIG: it flushes iptables and
# its MASQUERADE rewrites the source IP, which breaks the SIG return mapping. It is
# therefore NOT enabled by default. Enable it manually only for one-way flood/attack
# tests that do not need a return path:
#   docker exec endhost-as15 systemctl enable --now cyberflood-forwarding.service
# Rationale: docs/adr/0001-inyeccion-cyberflood-sig.md (Option B).
RUN apt-get install -y iptables iproute2 procps
COPY cyberflood-forwarding.sh /usr/local/bin/cyberflood-forwarding.sh
RUN chmod +x /usr/local/bin/cyberflood-forwarding.sh
COPY cyberflood-forwarding.service /etc/systemd/system/cyberflood-forwarding.service
# NOTE: intentionally NOT 'systemctl enable'd — opt-in only (see warning above).

# --- SCION-IP Gateway (SIG) — CyberFlood traffic through SCION paths ---
# SIG binary is installed from the SCION Debian packages in scion-base. The systemd
# unit runs /usr/bin/scion-ip-gateway (matches base/systemd/scion-ip-gateway.service).
# Verify inside the container with:  ls /usr/bin | grep -i gat
RUN mkdir -p /etc/scion/sig
COPY sig.toml /etc/scion/sig/sig.toml
COPY traffic.json /etc/scion/sig/traffic.json
COPY routing.policy /etc/scion/sig/routing.policy
COPY sig-routes.sh /usr/local/bin/sig-routes.sh
RUN chmod +x /usr/local/bin/sig-routes.sh
COPY scion-ip-gateway.service /etc/systemd/system/scion-ip-gateway.service
RUN systemctl enable scion-ip-gateway.service

CMD ["/lib/systemd/systemd", "--unit=multi-user.target"]
