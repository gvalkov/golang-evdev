all: ecodes.go

ecodes.go: ecodes.go.template
	./bin/generate-ecodes.sh /usr/include/linux/input.h $< > $@
.PHONY: ecodes.go
