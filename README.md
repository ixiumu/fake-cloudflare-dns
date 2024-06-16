# Fake Cloudflare DNS

Fake Cloudflare DNS is a local server that acts as a fake DNS resolver. Regardless of the requested domain name, it always returns an IP address pointing to the Cloudflare front-end CDN. It automatically looks for the fastest IP address in the background. You can provide a list of IP addresses or domain names in the ip.txt file, with one entry per line.

Usage:

```
./fake-cloudflare-dns -p 853 -f ip.txt
```

## How it works

The Fake Cloudflare DNS server intercepts DNS requests and responds with a predefined IP address associated with the Cloudflare front-end CDN. It continuously tests the connection speed to the provided IP addresses or domain names in the ip.txt file. The server selects the fastest IP address and responds with that address for all DNS queries.

The server listens on port 853 (as specified by the -p option) for incoming DNS requests. It reads the list of IP addresses or domain names from the ip.txt file (specified by the -f option). The server automatically determines the IP address with the fastest connection speed and responds with that address for all DNS queries.

Note: Make sure to adjust the paths and file names in the example code to match your setup.

## Requirements

Go programming language (for building the server)
A text file (ip.txt) containing IP addresses or domain names
Building the Server
To build the Fake Cloudflare DNS server, follow these steps:

Install Go on your system.
Create a directory and navigate to it.
Create a file named fake-cloudflare-dns.go and paste the server code into it.
Open a terminal and navigate to the directory where fake-cloudflare-dns.go is located.
Run the following command to build the server:

```
go build -o fake-dns fake-cloudflare-dns.go
```

After a successful build, you will have an executable file named fake-dns.
Running the Server
To run the Fake Cloudflare DNS server, use the following command:

```
./fake-dns -p 853 -f ip.txt
```

Make sure to adjust the port number (-p) and the path to the ip.txt file (-f) as per your requirements.

## Autostart Service (OpenWrt)

To automatically start the Fake Cloudflare DNS server on OpenWrt, follow these steps:

Create a new file, for example, `fake-dns`, and paste the autostart service code into it.

```
#!/bin/sh /etc/rc.common

START=99

start() {
    /usr/bin/fake-dns -p 55 2>/dev/null 1>/dev/null &
}

stop() {
    killall fake-dns
}
```

Move the file to the `/etc/init.d/` directory on your OpenWrt device and make the file executable:

```
chmod +x /etc/init.d/fake-dns
```

Enable the autostart service:

```
/etc/init.d/fake-dns enable
```

Now, the Fake Cloudflare DNS server will start automatically on system boot.

## Disclaimer

Please note that using a fake DNS server for deceptive purposes may violate network policies or laws in your jurisdiction. Use this software responsibly and ensure compliance with applicable regulations.