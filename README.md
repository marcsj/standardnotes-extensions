### Standard Notes Extensions server

- Auto-updating extensions server for your self-hosted Standard Notes server
- Simple to run, written in Go, Dockerfile included

### Usage

Use activation code `https://extensions.your.domain/index.json` in the activate tab on the Standard Notes application.

### Docker Compose example

```yaml
services:
  extensions:
    build: path/to/this/repo
    environment:
      - BASE_URL=https://extensions.your.domain
    expose:
      - 80
    volumes:
      - /var/notes/extensions:/repos
```
