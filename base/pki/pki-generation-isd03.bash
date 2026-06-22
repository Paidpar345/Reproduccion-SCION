#!/bin/bash

set -euo pipefail

mkdir -p /tmp/tutorial-scion-certs-isd03 && cd /tmp/tutorial-scion-certs-isd03
mkdir AS{11..15}

# Create voting and root keys and (self-signed) certificates for core ASes
pushd AS11
scion-pki certificate create --profile=sensitive-voting <(echo '{"isd_as": "18-ffaa:1:31", "common_name": "18-ffaa:1:31 sensitive voting cert"}') sensitive-voting.pem sensitive-voting.key
scion-pki certificate create --profile=regular-voting <(echo '{"isd_as": "18-ffaa:1:31", "common_name": "18-ffaa:1:31 regular voting cert"}') regular-voting.pem regular-voting.key
scion-pki certificate create --profile=cp-root <(echo '{"isd_as": "18-ffaa:1:31", "common_name": "18-ffaa:1:31 cp root cert"}') cp-root.pem cp-root.key
popd

pushd AS12
scion-pki certificate create --profile=cp-root <(echo '{"isd_as": "18-ffaa:1:32", "common_name": "18-ffaa:1:32 cp root cert"}') cp-root.pem cp-root.key
popd

pushd AS13
scion-pki certificate create --profile=sensitive-voting <(echo '{"isd_as": "18-ffaa:1:33", "common_name": "18-ffaa:1:33 sensitive voting cert"}') sensitive-voting.pem sensitive-voting.key
scion-pki certificate create --profile=regular-voting <(echo '{"isd_as": "18-ffaa:1:33", "common_name": "18-ffaa:1:33 regular voting cert"}') regular-voting.pem regular-voting.key
popd

# Create the TRC
mkdir -p tmp
cat <<EOF > trc-B1-S1-pld.tmpl
isd = 18
description = "ISD 18"
serial_version = 1
base_version = 1
voting_quorum = 2

core_ases = ["ffaa:1:31", "ffaa:1:32", "ffaa:1:33"]
authoritative_ases = ["ffaa:1:31", "ffaa:1:32", "ffaa:1:33"]
cert_files = ["AS11/sensitive-voting.pem", "AS11/regular-voting.pem", "AS11/cp-root.pem", "AS12/cp-root.pem", "AS13/sensitive-voting.pem", "AS13/regular-voting.pem"]

[validity]
not_before = __CURRENT_TIMESTAMP__
# 364d (no 365): la TRC se genera ~1-2 s despues que los certs raiz (validez 365d) y
# debe quedar DENTRO de su ventana; con 365d se sale por el final y scion-pki la
# rechaza ("TRC validity not covered by certificate"). 1 dia de margen.
validity = "364d"
EOF

sed -i "s/__CURRENT_TIMESTAMP__/$(date +%s)/" trc-B1-S1-pld.tmpl

scion-pki trc payload --out=tmp/ISD18-B1-S1.pld.der --template trc-B1-S1-pld.tmpl
rm trc-B1-S1-pld.tmpl

# Sign and bundle the TRC
scion-pki trc sign tmp/ISD18-B1-S1.pld.der AS11/sensitive-voting.{pem,key} --out tmp/ISD18-B1-S1.AS11-sensitive.trc
scion-pki trc sign tmp/ISD18-B1-S1.pld.der AS11/regular-voting.{pem,key} --out tmp/ISD18-B1-S1.AS11-regular.trc
scion-pki trc sign tmp/ISD18-B1-S1.pld.der AS13/sensitive-voting.{pem,key} --out tmp/ISD18-B1-S1.AS13-sensitive.trc
scion-pki trc sign tmp/ISD18-B1-S1.pld.der AS13/regular-voting.{pem,key} --out tmp/ISD18-B1-S1.AS13-regular.trc

scion-pki trc combine tmp/ISD18-B1-S1.AS{11,13}-{sensitive,regular}.trc --payload tmp/ISD18-B1-S1.pld.der --out ISD18-B1-S1.trc
rm tmp -r

# Create CA key and certificate for issuing ASes
pushd AS11
scion-pki certificate create --profile=cp-ca <(echo '{"isd_as": "18-ffaa:1:31", "common_name": "18-ffaa:1:31 CA cert"}') cp-ca.pem cp-ca.key --ca cp-root.pem --ca-key cp-root.key --not-after 2160h
popd
pushd AS12
scion-pki certificate create --profile=cp-ca <(echo '{"isd_as": "18-ffaa:1:32", "common_name": "18-ffaa:1:32 CA cert"}') cp-ca.pem cp-ca.key --ca cp-root.pem --ca-key cp-root.key --not-after 2160h
popd

# Create AS key and certificate chains
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "18-ffaa:1:31", "common_name": "18-ffaa:1:31 AS cert"}') AS11/cp-as.pem AS11/cp-as.key --ca AS11/cp-ca.pem --ca-key AS11/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "18-ffaa:1:32", "common_name": "18-ffaa:1:32 AS cert"}') AS12/cp-as.pem AS12/cp-as.key --ca AS12/cp-ca.pem --ca-key AS12/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "18-ffaa:1:33", "common_name": "18-ffaa:1:33 AS cert"}') AS13/cp-as.pem AS13/cp-as.key --ca AS11/cp-ca.pem --ca-key AS11/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "18-ffaa:1:34", "common_name": "18-ffaa:1:34 AS cert"}') AS14/cp-as.pem AS14/cp-as.key --ca AS11/cp-ca.pem --ca-key AS11/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "18-ffaa:1:35", "common_name": "18-ffaa:1:35 AS cert"}') AS15/cp-as.pem AS15/cp-as.key --ca AS12/cp-ca.pem --ca-key AS12/cp-ca.key --not-after 720h --bundle

echo 'copying to shared folder'
cp -r /tmp/tutorial-scion-certs-isd03 /shared/