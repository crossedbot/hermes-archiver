#!/bin/bash

REPLAYER=$(command -v replayer)

log()
{
    echo "$(date +"%F %T"): $*"
}

usage()
{
    echo -e "$(basename "$0") [-h] [-p <port>] [-d <database address>] [-k <key>] [-s <salt>] [i <ipfs address>] -- program to start the Hermes replayer

where:
    -h  show this help text
    -c  configuration file location; default is '${HOME}/.hermes/replayer.conf'
    -p  set listening port of HTTP API; default is 8888
    -d  set indexed record database address; default is '127.0.0.1:6379'
    -k  set encryption key; default is ''
    -s  set encryption salt; default is ''
    -i  set the IPFS node multi-address; default is '/ip4/127.0.0.1/tcp/5001'"
    exit
}

# START #

conf="${HOME}/.hermes/replayer.conf"
port="8989"
warc="${HOME}/.hermes/warc"
db_addr="127.0.0.1:6379"
key=""
salt=""
ipfs_addr="/ip4/127.0.0.1/tcp/5001"

while getopts "hc:p:d:k:s:i:" opt; do
    case "$opt" in
    [h?]) usage
        ;;
    c)  conf="${OPTARG}"
        ;;
    p)  port="${OPTARG}"
        ;;
    d)  db_addr="${OPTARG}"
        ;;
    k)  key="${OPTARG}"
        ;;
    s)  salt="${OPTARG}"
        ;;
    i)  ipfs_addr="${OPTARG}"
        ;;
    esac
done

conf_dir=$(dirname ${conf})
if [ ! -d ${conf_dir} ]; then
    log $(mkdir -vp ${conf_dir})
fi

if [ ! -d ${warc} ]; then
    log $(mkdir -vp ${warc})
fi

cat <<EOF > ${conf}
host="0.0.0.0"
port=${port}
read_timeout=30
write_timeout=30
database_addr="${db_addr}"
encryption_key="${key}"
encryption_salt="${salt}"
ipfs_address="${ipfs_addr}"
EOF
log "created '${conf}':
$(cat ${conf}| sed 's/^/\t/')"

$REPLAYER --config-file="${conf}"
