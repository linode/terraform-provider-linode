# Change Log

<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/chiefy/linodego/compare/v0.0.1...v0.1.0) (2018-07-30)

Adds more Domain handling

### Fixed

- go-resty doesnt pass errors when content-type is not set
- Domain, DomainRecords, tests and fixtures

### Added

- add CreateDomainRecord, UpdateDomainRecord, and DeleteDomainRecord

<a name="v0.1.0"></a>
## [v0.1.0](https://github.com/chiefy/linodego/compare/v0.0.1...v0.1.0) (2018-07-23)

Deals with NewClient and context for all http requests

### Breaking Changes

- changed `NewClient(token, *http.RoundTripper)` to `NewClient(*http.Client)`
- changed all `Client` `Get`, `List`, `Create`, `Update`, `Delete`, and `Wait` calls to take context as the first parameter

### Fixed

- fixed docs should now show Examples for more functions

### Added

- added `Client.SetBaseURL(url string)`

<a name="v0.0.1"></a>
## v0.0.1 (2018-07-20)

### Changed

* Initial tagged release
