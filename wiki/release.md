
# Release

- Export `$VERSION` env variable to the version you want
- Run `make release`

## Release without script

let's say the next tag is `v1.0.0`, release will be:

### SDK

- Run the command `make pulumi go_sdk nodejs_sdk python_sdk`
- Create a tag `git tag -a sdk/v1.0.0 -m "release sdk/v1.0.0"`
- Push the tag `git push origin sdk/v1.0.0`

### Main release

- Check `goreleaser check`
- Create a tag `git tag -a v1.0.0 -m "release v1.0.0"`
- Push the tag `git push origin v1.0.0`
- the release workflow will release the tag automatically.
- [Pulumi registry](https://github.com/pulumi/registry) twice-daily scans for new releases of Pulumi packages hence it will pick up the new release and create a new PR for that to be published.

## Tags Convention

The following convention should be followed for tagging in this project:

Release Tags: For release names and GitHub tags, the tag format should be prefixed with v0.0.0. For example, a release tag could be v1.2.3, where 1.2.3 represents the version number of the release.

Following this convention will help maintain consistency and clarity in tagging across all the grid components.
