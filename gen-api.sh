docker run -v `pwd`:/defs namely/protoc-all -f api/*/*.proto -l go -o internal/ --with-gateway
