#!/bin/bash

set -euo pipefail

mkdir -p /tmp/tutorial-scion-certs-isd04 && cd /tmp/tutorial-scion-certs-isd04
mkdir AS{16..20}

# Create voting and root keys and (self-signed) certificates for core ASes
pushd AS16
scion-pki certificate create --profile=sensitive-voting <(echo '{"isd_as": "19-ffaa:1:41", "common_name": "19-ffaa:1:41 sensitive voting cert"}') sensitive-voting.pem sensitive-voting.key
scion-pki certificate create --profile=regular-voting <(echo '{"isd_as": "19-ffaa:1:41", "common_name": "19-ffaa:1:41 regular voting cert"}') regular-voting.pem regular-voting.key
scion-pki certificate create --profile=cp-root <(echo '{"isd_as": "19-ffaa:1:41", "common_name": "19-ffaa:1:41 cp root cert"}') cp-root.pem cp-root.key
popd

pushd AS17
scion-pki certificate create --profile=cp-root <(echo '{"isd_as": "19-ffaa:1:42", "common_name": "19-ffaa:1:42 cp root cert"}') cp-root.pem cp-root.key
popd

pushd AS18
scion-pki certificate create --profile=sensitive-voting <(echo '{"isd_as": "19-ffaa:1:43", "common_name": "19-ffaa:1:43 sensitive voting cert"}') sensitive-voting.pem sensitive-voting.key
scion-pki certificate create --profile=regular-voting <(echo '{"isd_as": "19-ffaa:1:43", "common_name": "19-ffaa:1:43 regular voting cert"}') regular-voting.pem regular-voting.key
popd

# Create the TRC
mkdir -p tmp
cat <<EOF > trc-B1-S1-pld.tmpl
isd = 19
description = "ISD 19"
serial_version = 1
base_version = 1
voting_quorum = 2

core_ases = ["ffaa:1:41", "ffaa:1:42", "ffaa:1:43"]
authoritative_ases = ["ffaa:1:41", "ffaa:1:42", "ffaa:1:43"]
cert_files = ["AS16/sensitive-voting.pem", "AS16/regular-voting.pem", "AS16/cp-root.pem", "AS17/cp-root.pem", "AS18/sensitive-voting.pem", "AS18/regular-voting.pem"]

[validity]
not_before = __CURRENT_TIMESTAMP__
# 364d (no 365): la TRC se genera ~1-2 s despues que los certs raiz (validez 365d) y
# debe quedar DENTRO de su ventana; con 365d se sale por el final y scion-pki la
# rechaza ("TRC validity not covered by certificate"). 1 dia de margen.
validity = "364d"
EOF

sed -i "s/__CURRENT_TIMESTAMP__/$(date +%s)/" trc-B1-S1-pld.tmpl

scion-pki trc payload --out=tmp/ISD19-B1-S1.pld.der --template trc-B1-S1-pld.tmpl
rm trc-B1-S1-pld.tmpl

# Sign and bundle the TRC
scion-pki trc sign tmp/ISD19-B1-S1.pld.der AS16/sensitive-voting.{pem,key} --out tmp/ISD19-B1-S1.AS16-sensitive.trc
scion-pki trc sign tmp/ISD19-B1-S1.pld.der AS16/regular-voting.{pem,key} --out tmp/ISD19-B1-S1.AS16-regular.trc
scion-pki trc sign tmp/ISD19-B1-S1.pld.der AS18/sensitive-voting.{pem,key} --out tmp/ISD19-B1-S1.AS18-sensitive.trc
scion-pki trc sign tmp/ISD19-B1-S1.pld.der AS18/regular-voting.{pem,key} --out tmp/ISD19-B1-S1.AS18-regular.trc

scion-pki trc combine tmp/ISD19-B1-S1.AS{16,18}-{sensitive,regular}.trc --payload tmp/ISD19-B1-S1.pld.der --out ISD19-B1-S1.trc
rm tmp -r

# Create CA key and certificate for issuing ASes
pushd AS16
scion-pki certificate create --profile=cp-ca <(echo '{"isd_as": "19-ffaa:1:41", "common_name": "19-ffaa:1:41 CA cert"}') cp-ca.pem cp-ca.key --ca cp-root.pem --ca-key cp-root.key --not-after 2160h
popd
pushd AS17
scion-pki certificate create --profile=cp-ca <(echo '{"isd_as": "19-ffaa:1:42", "common_name": "19-ffaa:1:42 CA cert"}') cp-ca.pem cp-ca.key --ca cp-root.pem --ca-key cp-root.key --not-after 2160h
popd

# Create AS key and certificate chains
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "19-ffaa:1:41", "common_name": "19-ffaa:1:41 AS cert"}') AS16/cp-as.pem AS16/cp-as.key --ca AS16/cp-ca.pem --ca-key AS16/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "19-ffaa:1:42", "common_name": "19-ffaa:1:42 AS cert"}') AS17/cp-as.pem AS17/cp-as.key --ca AS17/cp-ca.pem --ca-key AS17/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "19-ffaa:1:43", "common_name": "19-ffaa:1:43 AS cert"}') AS18/cp-as.pem AS18/cp-as.key --ca AS16/cp-ca.pem --ca-key AS16/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "19-ffaa:1:44", "common_name": "19-ffaa:1:44 AS cert"}') AS19/cp-as.pem AS19/cp-as.key --ca AS16/cp-ca.pem --ca-key AS16/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "19-ffaa:1:45", "common_name": "19-ffaa:1:45 AS cert"}') AS20/cp-as.pem AS20/cp-as.key --ca AS17/cp-ca.pem --ca-key AS17/cp-ca.key --not-after 720h --bundle

echo 'copying to shared folder'
cp -r /tmp/tutorial-scion-certs-isd04 /shared/