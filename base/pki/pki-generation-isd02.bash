#!/bin/bash

set -euo pipefail

mkdir -p /tmp/tutorial-scion-certs-isd02 && cd /tmp/tutorial-scion-certs-isd02
mkdir AS{6..10}

# Create voting and root keys and (self-signed) certificates for core ASes
pushd AS6
scion-pki certificate create --profile=sensitive-voting <(echo '{"isd_as": "17-ffaa:1:21", "common_name": "17-ffaa:1:21 sensitive voting cert"}') sensitive-voting.pem sensitive-voting.key
scion-pki certificate create --profile=regular-voting <(echo '{"isd_as": "17-ffaa:1:21", "common_name": "17-ffaa:1:21 regular voting cert"}') regular-voting.pem regular-voting.key
scion-pki certificate create --profile=cp-root <(echo '{"isd_as": "17-ffaa:1:21", "common_name": "17-ffaa:1:21 cp root cert"}') cp-root.pem cp-root.key
popd

pushd AS7
scion-pki certificate create --profile=cp-root <(echo '{"isd_as": "17-ffaa:1:22", "common_name": "17-ffaa:1:22 cp root cert"}') cp-root.pem cp-root.key
popd

pushd AS8
scion-pki certificate create --profile=sensitive-voting <(echo '{"isd_as": "17-ffaa:1:23", "common_name": "17-ffaa:1:23 sensitive voting cert"}') sensitive-voting.pem sensitive-voting.key
scion-pki certificate create --profile=regular-voting <(echo '{"isd_as": "17-ffaa:1:23", "common_name": "17-ffaa:1:23 regular voting cert"}') regular-voting.pem regular-voting.key
popd

# Create the TRC
mkdir -p tmp
cat <<EOF > trc-B1-S1-pld.tmpl
isd = 17
description = "ISD 17"
serial_version = 1
base_version = 1
voting_quorum = 2

core_ases = ["ffaa:1:21", "ffaa:1:22", "ffaa:1:23"]
authoritative_ases = ["ffaa:1:21", "ffaa:1:22", "ffaa:1:23"]
cert_files = ["AS6/sensitive-voting.pem", "AS6/regular-voting.pem", "AS6/cp-root.pem", "AS7/cp-root.pem", "AS8/sensitive-voting.pem", "AS8/regular-voting.pem"]

[validity]
not_before = __CURRENT_TIMESTAMP__
# 364d (no 365): la TRC se genera ~1-2 s despues que los certs raiz (validez 365d) y
# debe quedar DENTRO de su ventana; con 365d se sale por el final y scion-pki la
# rechaza ("TRC validity not covered by certificate"). 1 dia de margen.
validity = "364d"
EOF

sed -i "s/__CURRENT_TIMESTAMP__/$(date +%s)/" trc-B1-S1-pld.tmpl


scion-pki trc payload --out=tmp/ISD17-B1-S1.pld.der --template trc-B1-S1-pld.tmpl
rm trc-B1-S1-pld.tmpl

# Sign and bundle the TRC
scion-pki trc sign tmp/ISD17-B1-S1.pld.der AS6/sensitive-voting.{pem,key} --out tmp/ISD17-B1-S1.AS6-sensitive.trc
scion-pki trc sign tmp/ISD17-B1-S1.pld.der AS6/regular-voting.{pem,key} --out tmp/ISD17-B1-S1.AS6-regular.trc
scion-pki trc sign tmp/ISD17-B1-S1.pld.der AS8/sensitive-voting.{pem,key} --out tmp/ISD17-B1-S1.AS8-sensitive.trc
scion-pki trc sign tmp/ISD17-B1-S1.pld.der AS8/regular-voting.{pem,key} --out tmp/ISD17-B1-S1.AS8-regular.trc

scion-pki trc combine tmp/ISD17-B1-S1.AS{6,8}-{sensitive,regular}.trc --payload tmp/ISD17-B1-S1.pld.der --out ISD17-B1-S1.trc
rm tmp -r

# Create CA key and certificate for issuing ASes
pushd AS6
scion-pki certificate create --profile=cp-ca <(echo '{"isd_as": "17-ffaa:1:21", "common_name": "17-ffaa:1:21 CA cert"}') cp-ca.pem cp-ca.key --ca cp-root.pem --ca-key cp-root.key --not-after 2160h
popd
pushd AS7
scion-pki certificate create --profile=cp-ca <(echo '{"isd_as": "17-ffaa:1:22", "common_name": "17-ffaa:1:22 CA cert"}') cp-ca.pem cp-ca.key --ca cp-root.pem --ca-key cp-root.key --not-after 2160h
popd

# Create AS key and certificate chains
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "17-ffaa:1:21", "common_name": "17-ffaa:1:21 AS cert"}') AS6/cp-as.pem AS6/cp-as.key --ca AS6/cp-ca.pem --ca-key AS6/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "17-ffaa:1:22", "common_name": "17-ffaa:1:22 AS cert"}') AS7/cp-as.pem AS7/cp-as.key --ca AS7/cp-ca.pem --ca-key AS7/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "17-ffaa:1:23", "common_name": "17-ffaa:1:23 AS cert"}') AS8/cp-as.pem AS8/cp-as.key --ca AS6/cp-ca.pem --ca-key AS6/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "17-ffaa:1:24", "common_name": "17-ffaa:1:24 AS cert"}') AS9/cp-as.pem AS9/cp-as.key --ca AS6/cp-ca.pem --ca-key AS6/cp-ca.key --not-after 720h --bundle
scion-pki certificate create --profile=cp-as <(echo '{"isd_as": "17-ffaa:1:25", "common_name": "17-ffaa:1:25 AS cert"}') AS10/cp-as.pem AS10/cp-as.key --ca AS7/cp-ca.pem --ca-key AS7/cp-ca.key --not-after 720h --bundle

echo 'copying to shared folder'
cp -r /tmp/tutorial-scion-certs-isd02 /shared/