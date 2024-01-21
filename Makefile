GO_OUT=.
DOCKER_IMAGE=namely/protoc-all
PROTO_FILE=balance.proto
GENERATED_FILES=$(patsubst %.proto,$(GO_OUT)/%.pb.go,$(PROTO_FILE))
# Правило для генерации кода Go из .proto файлов
$(GENERATED_FILES): 
	docker run -v $(PWD):/defs --rm $(DOCKER_IMAGE) -l go -d /defs -i /defs -o $(GO_OUT) --go-source-relative $(PROTO_FILE)

# Правило для генерации кода Java из .proto файлов
JAVA_OUT=src/main/java
GENERATED_JAVA_FILES=$(patsubst %.proto,$(JAVA_OUT)/%.java,$(PROTO_FILE))
all: $(GENERATED_FILES) $(GENERATED_JAVA_FILES) generate-transfer

# Правило для генерации кода Java из .proto файлов
$(GENERATED_JAVA_FILES): 
	docker run -v $(PWD):/defs --rm $(DOCKER_IMAGE) -l java -d /defs -i /defs -o $(JAVA_OUT) $(PROTO_FILE)

# Правило для генерации гошных обёрток transfer.proto
TRANSFER_GO_OUT=cmd/testing
TRANSFER_PROTO_FILE=transfer.proto
GENERATED_TRANSFER_FILES=$(patsubst %.proto,$(TRANSFER_GO_OUT)/%.pb.go,$(TRANSFER_PROTO_FILE))
generate-transfer: $(GENERATED_TRANSFER_FILES)

$(GENERATED_TRANSFER_FILES): 
	docker run -v $(PWD):/defs --rm $(DOCKER_IMAGE) -l go -d /defs -i /defs -o $(TRANSFER_GO_OUT) --go-source-relative $(TRANSFER_PROTO_FILE)
	rm -f $(TRANSFER_GO_OUT)/balance*