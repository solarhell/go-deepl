# Changelog

## 1.0.0 (2026-03-12)


### Features

* add `Client.Translation` method that returns a `Translation` ([4f89a6b](https://github.com/solarhell/go-deepl/commit/4f89a6b5684e5a4a3b63bc64b06886800ad3d1fa))
* add `HTTPClient()` option ([56907ac](https://github.com/solarhell/go-deepl/commit/56907ac72d403c691a44123d4a8b9601638adcfa))
* add new language codes ([34298ab](https://github.com/solarhell/go-deepl/commit/34298ab2ffd869214bbec862de98bad772d002a3))
* add support for html tag handling (closes [#5](https://github.com/solarhell/go-deepl/issues/5)) ([7f39d03](https://github.com/solarhell/go-deepl/commit/7f39d03e86967d4ad569f374170219f97e55898e))
* expose tag_handling and ignore_tags options ([1bab3f9](https://github.com/solarhell/go-deepl/commit/1bab3f900aec5c358c1fd1f582a33bb61698b98d))
* implement client ([ca09c59](https://github.com/solarhell/go-deepl/commit/ca09c59db45fee22243ff6af8aa25dce80cdf425))
* provide configuration through getters ([01fa0a1](https://github.com/solarhell/go-deepl/commit/01fa0a19b69c481f1ef339d50a9a4dd79a051a13))


### Bug Fixes

* czech language code (closes [#6](https://github.com/solarhell/go-deepl/issues/6)) ([a3a1742](https://github.com/solarhell/go-deepl/commit/a3a17423b924d4577b8040af0e755bff205ece40))
* ensure baseurl consistency when call after New ([ad6c5f7](https://github.com/solarhell/go-deepl/commit/ad6c5f7724aa65fb9612ed96413ccf2db43bef72))
* resolve errcheck lint issues in source and test files ([8c3cb52](https://github.com/solarhell/go-deepl/commit/8c3cb5254541bec46f72fdbc8882f077347b94df))
* typo in DeleteGlossary doc URL, update doc examples to errors.AsType ([59f1e09](https://github.com/solarhell/go-deepl/commit/59f1e094b90e2f4fd3f97ea0a57aede37a826b77))
