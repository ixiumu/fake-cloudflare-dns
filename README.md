# Fake Cloudflare DNS

Fake Cloudflare DNS is a local server that acts as a fake DNS resolver. Regardless of the requested domain name, it always returns an IP address pointing to the Cloudflare front-end CDN. It automatically looks for the fastest IP address in the background. You can provide a list of IP addresses or domain names in the ip.txt file, with one entry per line.

Usage:

```
  -dns string
        Upstream DNS Server (default "8.8.8.8:53")
  -domain string
        Default domain (default "creativecommons.org")
  -f string
        File name (default "ip.txt")
  -i int
        Each speed measurement interval (default 360)
  -log string
        Log level: none | err | info (default "info")
  -p int
        Port number (default 53)
  -t int
        Ping times (default 3)
```

## Disclaimer

Please note that using a fake DNS server for deceptive purposes may violate network policies or laws in your jurisdiction. Use this software responsibly and ensure compliance with applicable regulations.