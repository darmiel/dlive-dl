[![Build](https://github.com/darmiel/dlive-dl/actions/workflows/build.yml/badge.svg)](https://github.com/darmiel/dlive-dl/actions/workflows/build.yml)
# dlive-dl
Simple DLive Downloader

## Usage
```bash
$ dlive-dl dl [-u <url>] [-f video.ts]
```
**Example:*
```bash
$ dlive-dl dl -u https://dlive.tv/p/abcdefg+YVjyuzwGL
```

## Build
```bash
$ go build ./cmd/dlive-dl
```

## Install
```bash
$ go install ./cmd/dlive-dl
```