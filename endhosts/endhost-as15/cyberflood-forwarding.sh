#!/bin/bash
# =============================================================================
# CyberFlood NAT/Forwarding Setup
# Endhost: endhost-as15 (ISD1, AS5)
#
# Detects interfaces by subnet and configures:
#   - IP forwarding (sysctl)
#   - MASQUERADE on SCION-facing interfaces (transit_net, as_net)
#   - FORWARD rules from CyberFlood interface to SCION interfaces
# =============================================================================

set -euo pipefail

CYBERFLOOD_SUBNET="192.168.0."
TRANSIT_SUBNET="10.100.0."
AS_SUBNET="10.10.5."

log() { echo "[cyberflood-fwd] $*"; }

# --- Discover interfaces by IP address ---
find_iface_by_subnet() {
    local subnet="$1"
    ip -4 -o addr show | awk -v s="$subnet" '$4 ~ s {print $2}' | head -1
}

CF_IFACE=$(find_iface_by_subnet "$CYBERFLOOD_SUBNET")
TRANSIT_IFACE=$(find_iface_by_subnet "$TRANSIT_SUBNET")
AS_IFACE=$(find_iface_by_subnet "$AS_SUBNET")

if [ -z "$CF_IFACE" ]; then
    log "ERROR: No interface found on CyberFlood subnet ($CYBERFLOOD_SUBNET). Exiting."
    exit 1
fi

log "CyberFlood interface: $CF_IFACE"
log "Transit interface:    ${TRANSIT_IFACE:-not found}"
log "AS-local interface:   ${AS_IFACE:-not found}"

# --- Enable IP forwarding ---
sysctl -w net.ipv4.ip_forward=1
log "IP forwarding enabled"

# --- Flush existing NAT/FORWARD rules (idempotent) ---
iptables -t nat -F
iptables -F FORWARD
log "Flushed existing rules"

# --- MASQUERADE on outgoing SCION interfaces ---
for IFACE in $TRANSIT_IFACE $AS_IFACE; do
    if [ -n "$IFACE" ]; then
        iptables -t nat -A POSTROUTING -o "$IFACE" -j MASQUERADE
        log "MASQUERADE on $IFACE"
    fi
done

# --- FORWARD: CyberFlood -> SCION interfaces ---
for IFACE in $TRANSIT_IFACE $AS_IFACE; do
    if [ -n "$IFACE" ]; then
        iptables -A FORWARD -i "$CF_IFACE" -o "$IFACE" -j ACCEPT
        log "FORWARD $CF_IFACE -> $IFACE"
    fi
done

# --- FORWARD: Allow established/related return traffic ---
iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT
log "FORWARD RELATED,ESTABLISHED enabled"

log "CyberFlood forwarding setup complete"
