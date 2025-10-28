# GitHub Actions Workflows

This directory contains CI/CD workflows for the Kubiya Control Plane Terraform Provider.

## Workflows

### CI Workflow (`ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**

1. **Lint** - Code quality checks
   - Run `go fmt` and verify formatting
   - Run `go vet` for static analysis
   - Verify `go mod tidy` doesn't change anything

2. **Test** - Run tests
   - Test on multiple Go versions (1.22, 1.23)
   - Run tests with race detector
   - Generate coverage reports
   - Upload coverage to Codecov

3. **Build** - Build the provider
   - Compile the provider binary
   - Verify the binary builds successfully

4. **Validate** - Validate configurations
   - Install the provider
   - Validate example Terraform configurations

### Release Workflow (`release.yml`)

**Triggers:**
- Push of version tags (e.g., `v0.1.0`)

**Jobs:**

1. **GoReleaser** - Build and publish release
   - Build binaries for all platforms (Linux, macOS, Windows)
   - Sign binaries with GPG
   - Create GitHub release with artifacts
   - Generate checksums and manifest

**Required Secrets:**
- `GPG_PRIVATE_KEY` - GPG private key for signing
- `PASSPHRASE` - GPG key passphrase
- `GITHUB_TOKEN` - Automatically provided by GitHub

## Local Development

You can run CI checks locally using the Makefile:

```bash
# Run all CI checks (fmt, vet, test)
make ci

# Run individual checks
make fmt-check    # Check formatting
make vet          # Run go vet
make test         # Run tests
make validate     # Validate Terraform examples

# Run with coverage
make test-coverage
```

## Creating a Release

1. Update the CHANGELOG.md with release notes
2. Create a release using the Makefile:

```bash
make release VERSION=0.1.0
```

This will:
- Update the VERSION file
- Create a git commit with the version
- Create a git tag
- Prompt you to push the changes

3. Push the changes and tag:

```bash
git push && git push --tags
```

4. The release workflow will automatically:
   - Build binaries for all platforms
   - Sign the binaries
   - Create a GitHub release
   - Upload all artifacts

## Manual Release Steps

If you prefer to create releases manually:

```bash
# 1. Update version
echo "0.1.0" > VERSION

# 2. Update CHANGELOG.md with release notes

# 3. Commit changes
git add VERSION CHANGELOG.md
git commit -m "Release 0.1.0"

# 4. Create tag
git tag -a v0.1.0 -m "Release 0.1.0"

# 5. Push
git push && git push --tags
```

## Setting up GPG Signing

To set up GPG signing for releases:

1. Generate a GPG key:
```bash
gpg --full-generate-key
```

2. Export the private key:
```bash
gpg --export-secret-keys YOUR_KEY_ID | base64 > gpg_private_key.txt
```

3. Add secrets to your GitHub repository:
   - `GPG_PRIVATE_KEY`: Contents of `gpg_private_key.txt`
   - `PASSPHRASE`: Your GPG key passphrase

4. Update `.goreleaser.yml` with your GPG fingerprint or use the environment variable

## GoReleaser Configuration

The release process is configured in `.goreleaser.yml` at the project root. This file defines:
- Build targets (OS/architecture combinations)
- Binary naming
- Archive formats
- Checksum generation
- Signing configuration

For more information, see the [GoReleaser documentation](https://goreleaser.com/intro/).
