#! /bin/sh
set -e

case "$1" in
	"cli") 
		shift
		exec /cli "$@" ;;
	"server") 
		shift
		exec /server "$@" ;;
	*) exec "$@"
esac
