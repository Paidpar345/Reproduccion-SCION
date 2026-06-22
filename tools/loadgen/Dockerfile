# =============================================================================
# Generador de carga HTTP *cliente-solo* para medir ancho de banda a traves del
# SIG/SCION (Via C). Sustituye a CyberFlood, que sin licencia de servidor-externo
# se queda en closed-loop (emula cliente Y servidor) y se contesta a si mismo sin
# cruzar el SIG.
#
# Vive en la SUBRED A (192.168.0.0/24) y pega al Caddy de as35 (192.168.1.11,
# subred B). Como son /24 distintos, un cliente REAL trata el destino como
# off-subnet y enruta por su gateway (as15) -> SIG -> SCION. Cruza el escenario.
#
# Herramientas:
#   - h2load (nghttp2-client): carga HTTP/2 y HTTP/1.1 con req/s, throughput y latencias.
#   - iperf3: ancho de banda TCP puro (sanity-check).
#   - curl:   peticion suelta para confirmar el camino y la latencia.
# =============================================================================
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        nghttp2-client \
        iperf3 \
        curl \
        ca-certificates \
        iproute2 \
    && rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
