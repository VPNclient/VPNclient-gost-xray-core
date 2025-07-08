#!/bin/bash

# GOST Certificate Generator Script


# https://www.dmosk.ru/miniinstruktions.php?mini=gost-openssl-ubuntu
# https://www.altlinux.org/%D0%93%D0%9E%D0%A1%D0%A2_%D0%B2_OpenSSL

apt install libengine-gost-openssl1.1
openssl version

find / -name gost.so
mv /gost.so /usr/lib/x86_64-linux-gnu/

vi /etc/ssl/openssl.cnf


openssl_conf = openssl_def

[openssl_def]
engines = engine_section

[engine_section]
gost = gost_section

[gost_section]
engine_id = gost
dynamic_path = /usr/lib/x86_64-linux-gnu/gost.so
default_algorithms = ALL
CRYPT_PARAMS = id-Gost28147-89-CryptoPro-A-ParamSet

openssl ciphers | tr ':' '\n' | grep GOST


LEGACY-GOST2012-GOST8912-GOST8912
IANA-GOST2012-GOST8912-GOST8912
GOST2001-GOST89-GOST89


openssl req -x509 -newkey gost2012_256 -pkeyopt paramset:A -nodes -keyout cert.key -out cert.pem

