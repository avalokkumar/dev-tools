# Homebrew formula for DevForge.
#
# Usage (after publishing a release):
#
#   brew tap devforge/tap
#   brew install devforge
#
# Or directly from this file:
#
#   brew install --build-from-source ./packaging/homebrew/devforge.rb
#
# When cutting a release, replace the version + SHA-256 placeholders. If you
# use goreleaser's `brews` block (see .goreleaser.yaml), goreleaser maintains
# this file automatically in your tap repository — keep this copy as the
# canonical template.

class Devforge < Formula
  desc     "Unified, MCP-native developer toolkit (CLI + Web UI + MCP server)"
  homepage "https://github.com/devforge/devforge"
  version  "0.1.0"
  license  "MIT"

  on_macos do
    on_arm do
      url     "https://github.com/devforge/devforge/releases/download/v#{version}/devforge_#{version}_darwin_arm64.tar.gz"
      sha256  "REPLACE_WITH_SHA256_OF_devforge_#{version}_darwin_arm64.tar.gz"
    end
    on_intel do
      url     "https://github.com/devforge/devforge/releases/download/v#{version}/devforge_#{version}_darwin_amd64.tar.gz"
      sha256  "REPLACE_WITH_SHA256_OF_devforge_#{version}_darwin_amd64.tar.gz"
    end
  end

  on_linux do
    on_arm do
      url     "https://github.com/devforge/devforge/releases/download/v#{version}/devforge_#{version}_linux_arm64.tar.gz"
      sha256  "REPLACE_WITH_SHA256_OF_devforge_#{version}_linux_arm64.tar.gz"
    end
    on_intel do
      url     "https://github.com/devforge/devforge/releases/download/v#{version}/devforge_#{version}_linux_amd64.tar.gz"
      sha256  "REPLACE_WITH_SHA256_OF_devforge_#{version}_linux_amd64.tar.gz"
    end
  end

  def install
    bin.install "devforge"
    # Bundle the docs that ship inside every release archive.
    %w[README.md LICENSE].each { |f| doc.install f if File.exist?(f) }
    (doc/"docs").install Dir["docs/*.md"] if Dir.exist?("docs")

    # Generate shell completions on demand. `devforge` ships a hidden
    # `completion <shell>` command via cobra by default; if absent, this
    # block is a no-op.
    generate_completions_from_executable(bin/"devforge", "completion") if respond_to?(:generate_completions_from_executable)
  end

  test do
    # `version --json` is contractually stable.
    output = shell_output("#{bin}/devforge version --json")
    assert_match(/"version"\s*:/, output)

    # `run --list` must print every registered operation.
    list = shell_output("#{bin}/devforge run --list")
    assert_match(/uuid_generate/, list)
    assert_match(/json_format/,   list)
  end
end
