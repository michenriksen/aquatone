# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Added

### Changed


## [0.2.0]
### Added
 - New Collector: riddler.io (Thanks, [@jolle](https://github.com/jolle)!)
 - New Collector: crt.sh (Thanks, [@jolle](https://github.com/jolle)!)
 - New Collector: censys.io (Thanks, [@vortexau](https://github.com/vortexau)!)
 - New Collector: passivetotal.org

### Changed
 - Capture potential `NameError` exception in `asked_for_progress?` method,
   related to known JRuby bug (issue #4)
 - Capture potential `Errno::EBADF` exception in `asked_for_progress?` method (issue #15)
 - Improve handling of error when aquatone-gather is run on a system without a graphical desktop session (X11)
 - Exclude hosts resolving to broadcast addresses in aquatone-discover (issue #11)


## [0.1.1]
### Added

### Changed
- Capture `Errno::ENETUNREACH` exception in aquatone-scan to prevent it from
  erroring out when networks are unreachable.

## 0.1.0
### Added
- Initial release

### Changed

[Unreleased]: https://github.com/michenriksen/aquatone/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/michenriksen/aquatone/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/michenriksen/aquatone/compare/v0.1.0...v0.1.1
