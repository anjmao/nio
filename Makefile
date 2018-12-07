tag:
	@git tag `grep -P '^\tversion = ' nio.go|cut -f2 -d'"'`
	@git tag|grep -v ^v
