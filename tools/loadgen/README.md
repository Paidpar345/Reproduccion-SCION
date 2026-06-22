# loadgen — medir ancho de banda a través de SCION (Vía C)

Generador de carga HTTP **cliente-solo** que **sí cruza el SIG**, a diferencia de
CyberFlood en closed-loop (que sin licencia de servidor-externo se contesta a sí mismo).

- Vive en la **subred A** (`192.168.0.50`), gateway = **`192.168.0.10`** (as15 = SIG ingress).
- Pega al **Caddy de as35** (`192.168.1.11`, subred B). Al ser `/24` distintos, enruta
  por as15 → SIG → SCION → as35. Cruza el escenario.

## Construir y arrancar

```bash
docker compose build loadgen
docker compose up -d loadgen

# Confirmar la ruta (debe salir 'via 192.168.0.10'):
docker exec loadgen ip route get 192.168.1.11
```

## Tests de ancho de banda

Objetos disponibles en Caddy: `obj_1M.bin`, `obj_10M.bin`, `obj_100M.bin`, `obj_1G.bin`
(datos aleatorios incompresibles; compresión OFF → throughput real).

```bash
# HTTP/2 sobre TLS (el endpoint principal; cert self-signed, h2load no lo valida)
docker exec loadgen h2load -n 200 -c 20 -m 1 https://192.168.1.11/obj_100M.bin

# HTTP/1.1 en claro (línea base, sin TLS)
docker exec loadgen h2load --h1 -n 200 -c 20 http://192.168.1.11/obj_100M.bin

# Ancho de banda TCP puro con iperf3 (download = -R, server -> client)
docker exec -d endhost-as35 iperf3 -s
docker exec loadgen iperf3 -c 192.168.1.11 -t 30 -P 8 -R

# Petición suelta: confirma camino + latencia + velocidad
docker exec loadgen curl -kso /dev/null \
  -w 'http=%{http_version} t=%{time_total}s speed=%{speed_download}B/s\n' \
  https://192.168.1.11/obj_10M.bin
```

`h2load` imprime al final `finished in Xs, N req/s, Y MB/s` y latencias (min/max/mean) →
esos son tus números de ancho de banda y latencia **a través de SCION**.

## Verificar que cruza el SIG (en otra terminal, durante el test)

```bash
# Contador del SIG en as15: debe subir mucho durante la carga
docker exec endhost-as15 sh -c 'curl -s 127.0.0.1:30456/metrics | grep gateway_ippkts_sent_total'

# El tráfico dentro del túnel (debe verse 192.168.0.50 > 192.168.1.11)
docker exec endhost-as15 tcpdump -ni sig 'host 192.168.1.11'
```

**El tell:** la latencia de `curl`/`h2load` será **~1 ms** (no ~0.06 ms como el bucle interno
de CyberFlood). Ese salto es el coste real de SCION = tu resultado para el TFG.
