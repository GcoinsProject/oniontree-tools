include Build.def

.PHONY: all rebuild generate clean

all: generate
	go build -v -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)" -o "$(TARGET)"

rebuild: generate
	go build -a -v -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -asmflags="$(ASMFLAGS)" -o "$(TARGET)"

generate:
	go generate

clean:
	$(RM) $(TARGET)

