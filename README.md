# 🐟 Flounder

A lightweight platform to help users build simple Gemini sites over http(s) and
serve those sites over http(s) and Gemini

For more information about Gemini: https://gemini.circumlunar.space/

Flounder is in ALPHA -- development and features are changing frequently,
especially as the Gemini spec and ecosystem remains relatively unstable.

See the flagship instance at https://flounder.online and
[gemini://flounder.online](gemini://flounder.online)

## Building and running locally Requirements:
* go >= 1.15
* sqlite dev libraries

To run locally, copy example-config.toml to flounder.toml, then run:

`mkdir files` `go run . serve`

Add the following to `/etc/hosts` (include any other users you want to create):

``` 
127.0.0.1 flounder.local admin.flounder.local proxy.flounder.local 
```

Then open `flounder.local:8165` in your browser. The default login/password is
admin/admin

## TLS Certs and Reverse Proxy

Gemini TLS certs are handled for you. For HTTP, when deploying your site,
you'll need a reverse proxy that does the following for you:

1. Cert for yourdomain.whatever
1. Wildcard cert for \*.yourdomain.whatever
2. On Demand cert for custom user domains

If you have a very small deployment of say, <100 users, for example, you can
use on demand certificates for all domains, and you can skip step 2.

However, for a larger deployment, you'll have to set up a wildcard cert.
Wildcard certs are a bit of a pain and difficult to do automatically, depending
on your DNS provider. For information on doing this via Caddy, you can follow
this guide:
https://caddy.community/t/how-to-use-dns-provider-modules-in-caddy-2/8148. 

For information on using certbot to manage wildcard certs, see this guide:
https://letsencrypt.org/docs/challenge-types/#dns-01-challenge

An example simple Caddyfile using on-demand certs is available in
Caddyfile.example

If you want to host using a server other than Caddy, there's no reason you
can't, but it may be more cumbersome to setup the http certs.

## Administration

Flounder is designed to be small, easy to host, and easy to administer. Signups
require manual admin approval.

## Development

Open a PR, or use one of the mailing lists on https://github.com/alexwennerberg

## Donate

If you'd like to support Flounder development, consider making a [Donation](https://www.buymeacoffee.com/alexwennerberg)
