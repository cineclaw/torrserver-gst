#!/bin/sh
set -e

FLAGS=""
[ -n "$TS_PORT" ] && FLAGS="${FLAGS} --port ${TS_PORT}"
[ -n "$TS_PATH" ] && FLAGS="${FLAGS} --path ${TS_PATH}"
[ -n "$TS_LOGFILE" ] && FLAGS="${FLAGS} --logpath ${TS_LOGPATHDIR}/${TS_LOGFILE}"
[ -n "$TS_WEBLOGFILE" ] && FLAGS="${FLAGS} --weblogpath ${TS_LOGPATHDIR}/${TS_WEBLOGFILE}"
[ -n "$TS_RDB" ] && FLAGS="${FLAGS} --rdb"
[ -n "$TS_HTTPAUTH" ] && FLAGS="${FLAGS} --httpauth"
[ -n "$TS_DONTKILL" ] && FLAGS="${FLAGS} --dontkill"
[ -n "$TS_TORRENTSDIR" ] && FLAGS="${FLAGS} --torrentsdir ${TS_TORRENTSDIR}"
[ -n "$TS_TORRENTADDR" ] && FLAGS="${FLAGS} --torrentaddr ${TS_TORRENTADDR}"
[ -n "$TS_PUBIPV4" ] && FLAGS="${FLAGS} --pubipv4 ${TS_PUBIPV4}"
[ -n "$TS_PUBIPV6" ] && FLAGS="${FLAGS} --pubipv6 ${TS_PUBIPV6}"
[ -n "$TS_SEARCHWA" ] && FLAGS="${FLAGS} --searchwa"

[ -n "$TS_PATH" ] && [ ! -d "$TS_PATH" ] && mkdir -p "$TS_PATH"
[ -n "$TS_LOGPATHDIR" ] && [ ! -d "$TS_LOGPATHDIR" ] && mkdir -p "$TS_LOGPATHDIR"
[ -n "$TS_TORRENTSDIR" ] && [ ! -d "$TS_TORRENTSDIR" ] && mkdir -p "$TS_TORRENTSDIR"

echo "=== CineClaw TorrServer GStreamer ==="
echo "Port: ${TS_PORT:-8090}"
echo "Flags: ${FLAGS}"
gst-inspect-1.0 --version | head -n 1

export GODEBUG=madvdontneed=1

exec /usr/bin/torrserver ${FLAGS}
