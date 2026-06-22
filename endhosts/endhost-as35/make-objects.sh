#!/bin/bash
# =============================================================================
# Objetos de prueba para medir ancho de banda HTTP a traves del SIG/SCION.
# - Datos INCOMPRESIBLES (/dev/urandom): si se comprimieran (gzip/zstd o un DUT),
#   el throughput medido seria ficticio. Por eso urandom + compresion OFF en Caddy.
# - IDEMPOTENTE: no regenera un objeto si ya existe con el tamano correcto.
# - Se ejecuta en RUNTIME (ExecStartPre de bw-server.service), NO en build, para
#   no inflar la imagen Docker con varios GB.
# =============================================================================
set -euo pipefail

ROOT=/srv/objects
mkdir -p "$ROOT"

# nombre -> bytes
declare -A SIZES=( [1M]=1048576 [10M]=10485760 [100M]=104857600 [1G]=1073741824 )

for name in "${!SIZES[@]}"; do
	bytes=${SIZES[$name]}
	f="$ROOT/obj_${name}.bin"
	if [ -f "$f" ] && [ "$(stat -c%s "$f" 2>/dev/null || echo 0)" = "$bytes" ]; then
		continue
	fi
	echo "[make-objects] generando $f ($bytes bytes)"
	head -c "$bytes" /dev/urandom > "$f"
done

# Manifiesto para reproducibilidad (verificar integridad entre corridas del TFG)
( cd "$ROOT" && sha256sum obj_*.bin > MANIFEST.sha256 2>/dev/null || true )
echo "[make-objects] objetos listos en $ROOT"
