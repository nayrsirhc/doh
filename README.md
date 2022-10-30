# doh
The purpose of this command and project in general, is to solve the problem of being terminal first, working behind a corporate firewall can be challenging, you might want to see what a DNS record resolves to externally, but port 53 is blocked for security reasons, DNS over HTTPS obviously uses port 443 to do so hence the birth of this command.

Here are a few features of this command:
- Written using the very powerful and convenient Cobra CLI Library
- Has reusable packages you can import into your project to use some of these capabilities, rough so far, I will work on refining these
- Robust and resolves against 4 different DNS providers
- Runs DNS queries concurrently and always returns the fastest provider first

`Usage`

```sh
DOH can be used to resolve DNS over HTTPS

This will be useful in the case you need to resolve external DNS for a domain but external resolution over port 53 is blocked

Usage:
  doh [command] [domain] [flags]
  doh [command]

Available Commands:
  a           Resolves A records for domain
  aaaa        Resolves AAAA records for domain
  all         Resolves all records for a domain
  any         Resolve whatever record type you find for this domain
  cname       Resolves CNAME records for domain
  completion  Generate the autocompletion script for the specified shell
  extensive   Resolves an extensive list of records for a domain
  help        Help about any command
  mx          Resolves MX/MAIL records for domain
  ns          Resolves NS records for domain
  soa         Resolves SOA for domain
  txt         Resolves TXT records for domain

Flags:
  -h, --help      help for doh
  -v, --version   version for doh

Use "doh [command] --help" for more information about a command.
```