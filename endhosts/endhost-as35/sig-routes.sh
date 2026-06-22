#!/bin/bash
# =============================================================================
# SIG Route Setup — endhost-as35  (WATCHER PERSISTENTE)
# Re-instala las rutas del túnel SIG cada vez que aparece el dispositivo 'sig' o
# sus rutas se vacían (p.ej. cuando el gateway refresca el túnel al cambiar los
# caminos disponibles). Así el IP-sobre-SCION sobrevive a reinicios y refrescos
# del gateway SIN necesidad de rebuild. No sale nunca: corre en background vía el
# ExecStartPost de scion-ip-gateway.service.
# =============================================================================

log() { echo "[sig-routes] $*"; }

# Prefijos que deben encaminarse por el túnel SIG (vuelta: as35 -> as15 y retorno
# a la subred de CyberFlood, que llega encapsulada por el túnel desde as15).
ROUTES="10.10.5.0/24 192.168.0.0/24 10.101.0.0/16"

# Pool de SERVIDORES (10.201.0.0/16) detras del Virtual Router (VR) de CyberFlood en subred B.
# La IDA (desencapsulada del SIG) NO va por 'dev sig': se ENTREGA al VR servidor.
# CF_VR_SERVER = IP del VR servidor de CyberFlood (gateway de la subred B). Ajusta si cambia.
CF_VR_SERVER="192.168.1.1"
SERVER_POOL="10.201.0.0/16"

log "Starting SIG route watcher (persistent)..."
sysctl -w net.ipv4.ip_forward=1 >/dev/null 2>&1 || true

while true; do
    if [ -d /sys/class/net/sig ]; then
        for r in $ROUTES; do
            if ! ip route show "$r" 2>/dev/null | grep -q 'dev sig'; then
                ip route replace "$r" dev sig 2>/dev/null && log "Route (re)installed: $r dev sig"
            fi
        done

        # Entrega del pool LOCAL (servidores) hacia el VR de CyberFlood (no por el tunel)
        if ! ip route show "$SERVER_POOL" 2>/dev/null | grep -q "via $CF_VR_SERVER"; then
            ip route replace "$SERVER_POOL" via "$CF_VR_SERVER" 2>/dev/null && log "Pool route: $SERVER_POOL via $CF_VR_SERVER"
        fi

        # Permitir forwarding entre la interfaz CyberFlood y el túnel SIG (si existe)
        CF_IFACE=$(ip -4 -o addr show | awk '/192\.168\.1\./ {print $2}' | head -1)
        if [ -n "$CF_IFACE" ]; then
            iptables -C FORWARD -i "$CF_IFACE" -o sig -j ACCEPT 2>/dev/null || \
                iptables -A FORWARD -i "$CF_IFACE" -o sig -j ACCEPT
            iptables -C FORWARD -i sig -o "$CF_IFACE" -j ACCEPT 2>/dev/null || \
                iptables -A FORWARD -i sig -o "$CF_IFACE" -j ACCEPT
            iptables -C FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT 2>/dev/null || \
                iptables -A FORWARD -m state --state RELATED,ESTABLISHED -j ACCEPT
        fi
    fi
    sleep 5
done
