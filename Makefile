tag:
	@git tag `grep -P '^\tversion = ' dapi.go|cut -f2 -d'"'`
	@git tag|grep -v ^v
