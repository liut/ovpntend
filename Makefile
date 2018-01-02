.SILENT :

WITH_ENV = env `cat .env 2>/dev/null | xargs`

test:
	mkdir -p tests
	@$(WITH_ENV) go test -v -cover -coverprofile tests/openVPNstatus.out ./ovpn/status
	@$(WITH_ENV) go tool cover -html=tests/openVPNstatus.out -o tests/openVPNstatus.out.html
