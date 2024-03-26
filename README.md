# Init

```bash
cd notes
ln -s $(pwd) $HOME/.tellme
```

# Build 

```bash
GOOS=darwin GOARCH=amd64 go build -o tellme .
chmod +x tellme
```
