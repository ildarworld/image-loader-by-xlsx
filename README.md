# Build for old versions of MacOS:

```bash
CGO_CFLAGS="-mmacosx-version-min=10.12" CGO_LDFLAGS="-mmacosx-version-min=10.12"  go build -o photo.app/Contents/MacOS/photo
```