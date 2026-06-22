#!/bin/sh
# Forzar que el trafico a la subred B (servidor, 192.168.1.0/24) salga por el SIG
# ingress = endhost-as15 (192.168.0.10), en vez de por el gateway por defecto del
# macvlan (192.168.0.1, que NO existe como router real). Este unico salto es lo que
# mete el trafico en el tunel SIG -> SCION.
set -e

ip route replace 192.168.1.0/24 via 192.168.0.10

echo "[loadgen] ruta a la subred B instalada:"
ip route get 192.168.1.11 || true
echo "[loadgen] listo. Ejemplos:"
echo "  docker exec loadgen h2load -n 200 -c 20 https://192.168.1.11/obj_100M.bin"
echo "  docker exec loadgen iperf3 -c 192.168.1.11 -t 30 -P 8 -R"

# Mantener el contenedor vivo para lanzar pruebas con 'docker exec'.
exec sleep infinity
