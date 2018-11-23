# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.2]

### Added
- Responsive URLs are now written to `aquatone_urls.txt`. Thanks [eur0pa](https://github.com/eur0pa)!
- A warning is printed when older versions of Chromium is detected which has known problems with screenshotting HTTPS URLs

### Fixed
- Aquatone had trouble processing a single or very few targets. A small delay has been added to give agents time to emit all their events

## [1.4.1]

### Changed
- List of User-Agents have been updated with most recent list of common User-Agents

### Fixed
- Random User-Agent and other spoofing request headers were not set correctly when requesting URLs. Thanks to [eur0pa](https://github.com/eur0pa) for pointing it out!

## [1.4.0]

### Added
- Passive fingerprinting of web technology in use on websites with Wappalyzer fingerprints
- Detection of domain takeover vulnerabilities across 20 different services

## [1.3.2]

Complete rewrite and simplification of Aquatone. Now written in Go and focused on reporting and screenshotting.

### Added
- Extraction of hosts, IPs and URLs from arbitrary data piped to Aquatone
- Parsing of Nmap/Masscan XML files
- Clustering of websites with similar structure in HTML report

### Removed
- Domain discovery (`aquatone-discover`)
- Domain takeover discovery (`aquatone-takeover`)

[Unreleased]: https://github.com/michenriksen/aquatone/compare/v1.4.2...HEAD
[1.4.2]: https://github.com/michenriksen/aquatone/compare/v1.4.1...1.4.2
[1.4.1]: https://github.com/michenriksen/aquatone/compare/v1.4.0...1.4.1
[1.4.0]: https://github.com/michenriksen/aquatone/compare/v1.3.2...1.4.0
[1.3.2]: https://github.com/michenriksen/aquatone/compare/v1.3.2
