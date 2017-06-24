# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]
### Added

### Changed
 - Capture potential `NameError` exception in `asked_for_progress?` method,
   related to known JRuby bug (issue #4)
 - Capture potential `Errno::EBADF` exception in `asked_for_progress?` method (issue #15)


## [0.1.1]
### Added

### Changed
- Capture `Errno::ENETUNREACH` exception in aquatone-scan to prevent it from
  erroring out when networks are unreachable.

## 0.1.0
### Added
- Initial release

### Changed

[Unreleased]: https://github.com/michenriksen/aquatone/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/michenriksen/aquatone/compare/v0.1.0...v0.1.1
