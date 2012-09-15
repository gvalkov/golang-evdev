all: install

ecodes.go: ecodes.tmpl.go
	./ecodes.sh /usr/include/linux/input.h $< > $@

install: ecodes.go
	go install .

test: ecodes.go
	go test .

getioctl: getioctl.c
	gcc $< -o $@
