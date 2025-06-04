# Release Process for ComplyTime

The release process values simplicity and automation in order to provide better predictability and low cost for maintainers.

## Process Description

Release artifacts are orchestrated by [GoReleaser](https://goreleaser.com/), which is configured in [.goreleaser.yaml](https://github.com/complytime/complytime/blob/main/.goreleaser.yaml)

There is a [Workflow](https://github.com/complytime/complytime/blob/main/.github/workflows/release.yml) created specifically for releases. This workflow is triggered manually by a project maintainer when a new release is ready to be published.

Once the automation is finished without issues, the release is available in [releases page](https://github.com/complytime/complytime/releases)

## Tests

Tests relevant for releases are incorporated in CI tests for every PR.

## Cadence

Releases are currently expected every three weeks. Project maintainers always discuss and agree on releases. Therefore, some releases may be triggered a bit earlier or later when necessary.
