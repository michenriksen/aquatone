# AQUATONE

Aquatone is a tool for visual inspection of websites across a large amount of hosts and is convenient for quickly gaining an overview of HTTP-based attack surface.

## Installation

1. Install [Google Chrome](https://www.google.com/chrome/) or [Chromium](https://www.chromium.org/getting-involved/download-chromium) browser -- **Note:** Google Chrome is currently giving unreliable results when running in *headless* mode, so it is recommended to install Chromium for the best results.
2. Download the [latest release](https://github.com/michenriksen/aquatone/releases/latest) of Aquatone for your operating system.
3. Uncompress the zip file and move the `aquatone` binary to your desired location. You probably want to move it to a location in your `$PATH` for easier use.

### Compiling the source code

If you for some reason don't trust the pre-compiled binaries, you can also compile the code yourself. **You are on your own if you want to do this. I do not support compiling problems. Good luck with it!**

## Usage

### Command-line options:

```
  -chrome-path string
    	Full path to the Chrome/Chromium executable to use. By default, aquatone will search for Chrome or Chromium
  -debug
    	Print debugging information
  -http-timeout int
    	Timeout in miliseconds for HTTP requests (default 3000)
  -nmap
    	Parse input as Nmap/Masscan XML
  -out string
    	Directory to write files to (default ".")
  -ports string
    	Ports to scan on hosts. Supported list aliases: small, medium, large, xlarge (default "80,443,8000,8080,8443")
  -proxy string
    	Proxy to use for HTTP requests
  -resolution string
    	screenshot resolution (default "1440,900")
  -save-body
    	Save response bodies to files (default true)
  -scan-timeout int
    	Timeout in miliseconds for port scans (default 100)
  -screenshot-timeout int
    	Timeout in miliseconds for screenshots (default 30000)
  -session string
    	Load Aquatone session file and generate HTML report
  -silent
    	Suppress all output except for errors
  -template-path string
    	Path to HTML template to use for report
  -threads int
    	Number of concurrent threads (default number of logical CPUs)
  -version
    	Print current Aquatone version
```

### Giving Aquatone data

Aquatone is designed to be as easy to use as possible and to integrate with your existing toolset with no or minimal glue. Aquatone is started by piping output of a command into the tool. It doesn't really care how the piped data looks as URLs, domains, and IP addresses will be extracted with regular expression pattern matching. This means that you can pretty much give it output of any tool you use for host discovery.

IPs, hostnames and domain names in the data will undergo scanning for ports that are typically used for web services and transformed to URLs with correct scheme.  If the data contains URLs, they are assumed to be alive and do not undergo port scanning.

**Example:**

    $ cat targets.txt | aquatone

### Output

When Aquatone is done processing the target hosts, it has created a bunch of files and folders in the current directory:

 - **aquatone_report.html**: An HTML report to open in a browser that displays all the collected screenshots and response headers clustered by similarity.
 - **aquatone_urls.txt**: A file containing all responsive URLs. Useful for feeding into other tools.
 - **aquatone_session.json**: A file containing statistics and page data. Useful for automation.
 - **headers/**: A folder with files containing raw response headers from processed targets
 - **html/**: A folder with files containing the raw response bodies from processed targets. If you are processing a large amount of hosts, and don't need this for further analysis, you can disable this with the `-save-body=false` flag to save some disk space.
 - **screenshots/**: A folder with PNG screenshots of the processed targets

The output can easily be zipped up and shared with others or archived.

#### Changing the output destination

If you don't want Aquatone to create files in the current working directory, you can specify a different location with the `-out` flag:

    $ cat hosts.txt | aquatone -out ~/aquatone/example.com

It is also possible to set a permanent default output destination by defining an environment variable:

    export AQUATONE_OUT_PATH="~/aquatone"


### Specifying ports to scan

Be default, Aquatone will scan target hosts with a small list of commonly used HTTP ports: 80, 443, 8000, 8080 and 8443. You can change this to your own list of ports with the `-ports` flag:

    $ cat hosts.txt | aquatone -ports 80,443,3000,3001

Aquatone also supports aliases of built-in port lists to make it easier for you:

 - **small**: 80, 443
 - **medium**: 80, 443, 8000, 8080, 8443 (same as default)
 - **large**: 80, 81, 443, 591, 2082, 2087, 2095, 2096, 3000, 8000, 8001, 8008, 8080, 8083, 8443, 8834, 8888
 - **xlarge**: 80, 81, 300, 443, 591, 593, 832, 981, 1010, 1311, 2082, 2087, 2095, 2096, 2480, 3000, 3128, 3333, 4243, 4567, 4711, 4712, 4993, 5000, 5104, 5108, 5800, 6543, 7000, 7396, 7474, 8000, 8001, 8008, 8014, 8042, 8069, 8080, 8081, 8088, 8090, 8091, 8118, 8123, 8172, 8222, 8243, 8280, 8281, 8333, 8443, 8500, 8834, 8880, 8888, 8983, 9000, 9043, 9060, 9080, 9090, 9091, 9200, 9443, 9800, 9981, 12443, 16080, 18091, 18092, 20720, 28017

**Example:**

    $ cat hosts.txt | aquatone -ports large


### Usage examples

Aquatone is designed to play nicely with all kinds of tools. Here's some examples:

#### Amass DNS enumeration

[Amass](https://github.com/OWASP/Amass) is currently my preferred tool for enumerating DNS. It uses a bunch of OSINT sources as well as active brute-forcing and clever permutations to quickly identify hundreds, if not thousands, of subdomains on a  domain:

```bash
$ amass -active -brute -o hosts.txt -d yahoo.com
alerts.yahoo.com
ads.yahoo.com
am.yahoo.com
- - - SNIP - - -
prd-vipui-01.infra.corp.gq1.yahoo.com
cp103.mail.ir2.yahoo.com
prd-vipui-01.infra.corp.bf1.yahoo.com
$ cat hosts.txt | aquatone
```

There are plenty of other DNS enumeration tools out there and Aquatone should work just as well with any other tool:

- [Sublist3r](https://github.com/aboul3la/Sublist3r)
- [Subfinder](https://github.com/subfinder/subfinder)
- [Knock](https://github.com/guelfoweb/knock)
- [Fierce](https://www.aldeid.com/wiki/Fierce)
- [Gobuster](https://github.com/OJ/gobuster)

#### Nmap or Masscan

Aquatone can make a report on hosts scanned with the [Nmap](https://nmap.org/) or [Masscan](https://github.com/robertdavidgraham/masscan) portscanner. Simply feed Aquatone the XML output and give it the `-nmap` flag to tell it to parse the input as Nmap/Masscan XML:

    $ cat scan.xml | aquatone -nmap

### Credits

- Thanks to [EdOverflow](https://twitter.com/EdOverflow) for the [can-i-take-over-xyz](https://github.com/EdOverflow/can-i-take-over-xyz/) project which Aquatone's domain takeover capability is based on.
- Thanks to [Elbert Alias](https://github.com/AliasIO) for the [Wappalyzer](https://github.com/AliasIO/Wappalyzer) project's technology fingerprints which Aquatone's technology fingerprinting capability is based on.
