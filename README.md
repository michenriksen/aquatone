# AQUATONE

AQUATONE is a set of tools for performing reconnaissance on domain names. It can
discover subdomains on a given domain by using open sources as well as the more
common subdomain dictionary brute force approach. After subdomain discovery,
AQUATONE can then scan the hosts for common web ports and HTTP headers, HTML
bodies and screenshots can be gathered and consolidated into a report for easy
analysis of the attack surface.

## Installation

### Dependencies

AQUATONE depends on [Node.js] and [NPM] package manager for its web page screenshotting capabilities. Follow [this guide] for Installation instructions.

You will also need a newer version of Ruby installed. If you plan to use AQUATONE in [Kali] Linux, you are already set up with this. If not, it is recommended to install Ruby with [RVM].

Finally, the tool itself can be installed with the following command in a terminal:

    $ gem install aquatone

**IMPORTANT:** AQUATONE's screenshotting capabilities depend on being run on a system with a graphical desktop environment. It is strongly recommended to install and run AQUATONE in a [Kali linux](https://www.kali.org/) virtual machine.
**I will not provide support or bug fixing for other systems than Kali Linux.**

## Usage

### Discovery

The first stage of an AQUATONE assessment is the discovery stage where subdomains are discovered on the target domain using open sources, services and the more common dictionary brute force approach:

    $ aquatone-discover --domain example.com

aquatone-discover will find the target's nameservers and shuffle DNS lookups between them. Should a lookup fail on the target domain's nameservers, aquatone-discover will fall back to using Google's public DNS servers to maximize discovery. The fallback DNS servers can be changed with the `--fallback-nameservers` option:

    $ aquatone-discover --domain example.com --fallback-nameservers 87.98.175.85,5.9.49.12

#### Tuning

aquatone-discover will use 5 threads as default for concurrently performing DNS lookups. This provides reasonable performance but can be tuned to be more or less aggressive with the `--threads` option:

    $ aquatone-discover --domain example.com --threads 25

Hammering a DNS server with failing lookups can potentially be picked up by intrusion detection systems, so if that is a concern for you, you can make aquatone-discover a bit more stealthy with the `--sleep` and `--jitter` options. `--sleep` accepts a number of seconds to sleep between each DNS lookup while `--jitter` accepts a percentage of the `--sleep` value to randomly add or subtract to or from the sleep interval in order to break the sleep pattern and make it less predictable.

    $ aquatone-discover --domain example.com --sleep 5 --jitter 30

Please note that setting the `--sleep` option will force the thread count to one. The `--jitter` option will only be considered if the `--sleep` option has also been set.

#### API keys

Some of the passive collectors will require API keys or similar credentials in order to work. Setting these values can be done with the `--set-key` option:

    $ aquatone-discover --set-key shodan o1hyw8pv59vSVjrZU3Qaz6ZQqgM91ihQ

All keys will be saved in  `~/aquatone/.keys.yml`.

#### Results

When aquatone-discover is finished, it will create a `hosts.txt` file in the `~/aquatone/<domain>` folder, so for a scan of example.com it would be located at `~/aquatone/example.com/hosts.txt`. The format will be a comma-separated list of hostnames and their IP, for example:

    example.com,93.184.216.34
    www.example.com,93.184.216.34
    secret.example.com,93.184.216.36
    cdn.example.com,192.0.2.42
    ...

In addition to the `hosts.txt` file, it will also generate a `hosts.json` which includes the same information but in JSON format. This format might be preferable if you want to use the information in custom scripts and tools. `hosts.json` will also be used by the aquatone-scan and aquatone-gather tools.

See `aquatone-discover --help` for more options.

### Scanning

The scanning stage is where AQUATONE will enumerate the discovered hosts for open TCP ports that are commonly used for web services:

    $ aquatone-scan --domain example.com

The `--domain` option will look for `hosts.json` in the domain's AQUATONE assessment directory, so in the example above it would look for `~/aquatone/example.com/hosts.json`. This file should be present if `aquatone-discover --domain example.com` has been run previously.

#### Ports

By default, aquatone-scan will scan the following TCP ports: 80, 443, 8000, 8080 and 8443. These are very common ports for web services and will provide a reasonable coverage. Should you want to specifiy your own list of ports, you can use the `--ports` option:

    $ aquatone-scan --domain example.com --ports 80,443,3000,8080

Instead of a comma-separated list of ports, you can also specify one of the built-in list aliases:

 * **small**: 80, 443
 * **medium**: 80, 443, 8000, 8080, 8443 (same as default)
 * **large**: 80, 81, 443, 591, 2082, 2095, 2096, 3000, 8000, 8001, 8008, 8080, 8083, 8443, 8834, 8888, 55672
 * **huge**: 80, 81, 300, 443, 591, 593, 832, 981, 1010, 1311, 2082, 2095, 2096, 2480, 3000, 3128, 3333, 4243, 4567, 4711, 4712, 4993, 5000, 5104, 5108, 5280, 5281, 5800, 6543, 7000, 7396, 7474, 8000, 8001, 8008, 8014, 8042, 8069, 8080, 8081, 8083, 8088, 8090, 8091, 8118, 8123, 8172, 8222, 8243, 8280, 8281, 8333, 8337, 8443, 8500, 8834, 8880, 8888, 8983, 9000, 9043, 9060, 9080, 9090, 9091, 9200, 9443, 9800, 9981, 11371, 12443, 16080, 18091, 18092, 20720, 55672

**Example:**

    $ aquatone-scan --domain example.com --ports large

#### Tuning

Like aquatone-discover, you can make the scanning more or less aggressive with the `--threads` option which accepts a number of threads for concurrent port scans. The default number of threads is 5.

    $ aquatone-scan --domain example.com --threads 25

As aquatone-scan is performing port scanning, it can obviously be picked up by intrusion detection systems. While it will attempt to lessen the risk of detection by randomising hosts and ports, you can tune the stealthiness more with the `--sleep` and `--jitter` options which work just like the similarly named options for aquatone-discover. Keep in mind that setting the `--sleep` option will force the number of threads to one.

#### Results

When aquatone-scan is finished, it will create a `urls.txt` file in the `~/aquatone/<domain>` directory, so for a scan of example.com it would be located at `~/aquatone/example.com/urls.txt`. The format will be a list of URLs, for example:

    http://example.com/
    https://example.com/
    http://www.example.com/
    https://www.example.com/
    http://secret.example.com:8001/
    https://secret.example.com:8443/
    http://cdn.example.com/
    https://cdn.example.com/
    ...

This file can be loaded into other tools such as [EyeWitness].

aquatone-scan will also generate a `open_ports.txt` file, which is a comma-separated list of hosts and their open ports, for example:

    93.184.216.34,80,443
    93.184.216.34,80
    93.184.216.36,80,443,8443
    192.0.2.42,80,8080
    ...

See `aquatone-scan --help` for more options.

### Gathering

The final stage is the gathering part where the results of the discovery and scanning stages are used to query the discovered web services in order to retrieve and save HTTP response headers and HTML bodies, as well as taking screenshots of how the web pages look like in a web browser to make analysis easier. The screenshotting is done with the [Nightmare.js] Node.js library. This library will be installed automatically if it's not present in the system.

    $ aquatone-gather --domain example.com

aquatone-gather will look for `hosts.json` and `open_ports.txt` in the given domain's AQUATONE assessment directory and request and screenshot every IP address for each domain name for maximum coverage.

#### Tuning

Like aquatone-discover and aquatone-scan, you can make the gathering more or less aggressive with the `--threads` option which accepts a number of threads for concurrent requests. The default number of threads is 5.

    $ aquatone-gather --domain example.com --threads 25

As aquatone-gather is interacting with web services, it can be picked up by intrusion detection systems. While it will attempt to lessen the risk of detection by randomising hosts and ports, you can tune the stealthiness more with the `--sleep` and `--jitter` options which work just like the similarly named options for aquatone-discover. Keep in mind that setting the `--sleep` option will force the number of threads to one.

#### Results

When aquatone-gather is finished, it will have created several directories in the domain's AQUATONE assessment directory:

 * `headers/`: Contains text files with HTTP response headers from each web page
 * `html/`: Contains text files with HTML response bodies from each web page
 * `screenshots/`: Contains PNG images of how each web page looks like in a browser
 * `report/` Contains report files in HTML displaying the gathered information for easy analysis

### Subdomain Takeover

Subdomain takeover is a very prevalent and potentially critical security issue which commonly occurs when an organization assigns a subdomain to a third-party service provider and then later discontinues use, but forgets to remove the DNS configuration. This leaves the subdomain vulnerable to complete takover by attackers by signing up to the same service provider and claiming the dangling subdomain.

aquatone-takeover can be used to check hosts uncovered by aquatone-discover for potential domain takeover vulnerabilities:

    $ aquatone-takeover --domain example.com

aquatone-takeover can detect potential subdomain takeover situations from 25 different service providers, including GitHub Pages, Heroku, Amazon S3, Desk and WPEngine.

#### Results

aquatone-takeover will create a `takeovers.json` file in the domain's assessment directory which will contain information in JSON format about any potential subdomain takeover vulnerabilities:

```
{
  "shop.example.com": {
    "service": "Shopify",
    "service_website": "https://www.shopify.com/",
    "description": "Ecommerce platform",
    "resource": {
      "type": "CNAME",
      "value": "shops.myshopify.com"
    }
  },
  "help.example.com": {
    "service": "Desk",
    "service_website": "https://www.desk.com/",
    "description": "Customer service and helpdesk ticket software",
    "resource": {
      "type": "CNAME",
      "value": "example.desk.com"
    }
  },
  ...
}
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/michenriksen/aquatone.


## License

AQUATONE is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).

[Node.js]: https://nodejs.org/
[NPM]: https://www.npmjs.com/
[this guide]: https://docs.npmjs.com/getting-started/installing-node
[Kali]: https://kali.org/
[RVM]: https://rvm.io/
[Nightmare.js]: http://www.nightmarejs.org/
[EyeWitness]: https://www.christophertruncer.com/eyewitness-usage-guide/
