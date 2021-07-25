# cloudflare-auto-dns
quick ddns polling hack for changing cloudflare DNS if your ip changes

## installing and usage
```bash
git clone https://github.com/cauefcr/cloudflare-auto-dns
cd cloudflare-auto-dns
go build .
export CF_API_KEY=yourapikey
export CF_API_EMAIL=your@email.com
export CF_ZONE_NAME=your.site
./cloudflare-auto-dns 
```
