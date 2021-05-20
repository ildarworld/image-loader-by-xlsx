# Mass file downloader by http requests from XLSX file

The app is easy, getting file hardcoded (especially by customer requests).
Then, app go through dedicated column (column name hard coded too) and downloads files in links whose are there.
links can be separated by semicolon.

## Build for old versions of MacOS:

```bash
CGO_CFLAGS="-mmacosx-version-min=10.12" CGO_LDFLAGS="-mmacosx-version-min=10.12"  go build -o photo.app/Contents/MacOS/photo
```
