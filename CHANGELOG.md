# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.7.0]

### Added
- New `session:start` and `session:end` events have been introduced in the event bus to allow agents to perform bootstrap
and cleanup tasks if needed
- A temporary user directory is now created for the Chrome/Chromium process and additional command line flags have been
added to increase compartmentalization

### Changed
- Production versions of Vue.js and Vue Router are now used in the HTML report for increased performance
- List of user agents have been updated to current list of most common user agents

## [1.7.0-beta.2]

### Fixed
- The pagination logic in the new HTML report would skip the page or cluster at index 0 as the `v-for` function on an integer value in Vue.js starts from 1 and not 0

## [1.7.0-beta]

### Added
- Session data will now be written to output directory as `aquatone_session.json`
- New `url_hostname_resolver` agent that resolves page's hostnames to IP addresses
- New `url_page_title_extractor` that extracts HTML page titles from responsive pages
- New command line flag `-template-path` to specify a custom template to use for the HTML report
- New command line flag `-session` to load a previous Aquatone session file and generate a report on its data
- Aquatone is now compiled for ARM64 in `build.sh`

### Changed
- Bigger refactoring of session and pages
- **New [Vue.js](https://vuejs.org/) powered HTML report with lots of new cool stuff:**
   - New look and feel
   - Pages can now be viewed in different modes:
      - **By Similarity**: Pages are displayed in clusters by their HTML structure similarity
      - **By Hostname:** Pages are displayed in clusters by their hostname
      - **Single Pages:** Pages are shown one-by-one with bigger screenshots and response headers (oldschool Aquatone style)
   - **[Vis.js](http://visjs.org/) powered network graph view** to see relations between pages, IP addresses and technologies
   - Page clusters are now rendered in a paginated carousel view instead of horizontally scrollable lanes
   - Clusters and pages are paginated to improve performance on large reports
   - Page titles are now shown for pages

### Removed
- `url_logger` agent (no longer needed)

## [1.6.0]

### Fixed
- The Nmap/Masscan XML report parser did not ignore closed/filtered ports. It now only works on ports with state `open`.

### Added
- Support for processing of multiple URLs on the same host by appending hash of URL path and fragment to file names
- Support for defining default output directory in `AQUATONE_OUT_PATH` environment variable

## [1.5.0]

### Added
- Automatic SSL/TLS detection on non-standard ports
- URL Screenshotter agent now takes extra steps to ensure that the browser process is killed after use
- Version flag to output current version (woah!!!)

### Changed
- Packages and other dependencies have been updated to latest versions
- User-Agent list has been updated to current most common agents
- Wappalyzer technology fingerprints have been updated

## [1.4.3]

### Fixed
- The Sub Resource Integrity attribute on the external CSS resource in the HTML report caused it to not load as the file had changed. Removed SRI on external CSS resource

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

[Unreleased]: https://github.com/michenriksen/aquatone/compare/v1.7.0...HEAD
[1.7.0]: https://github.com/michenriksen/aquatone/compare/v1.7.0-beta.2...v1.7.0
[1.7.0-beta.2]: https://github.com/michenriksen/aquatone/compare/v1.7.0-beta...v1.7.0-beta.2
[1.7.0-beta]: https://github.com/michenriksen/aquatone/compare/v1.6.0...v1.7.0-beta
[1.6.0]: https://github.com/michenriksen/aquatone/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/michenriksen/aquatone/compare/v1.4.3...v1.5.0
[1.4.3]: https://github.com/michenriksen/aquatone/compare/v1.4.2...v1.4.3
[1.4.2]: https://github.com/michenriksen/aquatone/compare/v1.4.1...v1.4.2
[1.4.1]: https://github.com/michenriksen/aquatone/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/michenriksen/aquatone/compare/v1.3.2...v1.4.0
[1.3.2]: https://github.com/michenriksen/aquatone/compare/v1.3.2
