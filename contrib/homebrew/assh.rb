require "language/go"

class Assh < Formula
  desc "assh: Advanced SSH config - A transparent wrapper that adds regex, aliases, gateways, includes, dynamic hostnames to SSH"
  homepage "https://github.com/moul/advanced-ssh-config"
  url "https://github.com/moul/advanced-ssh-config/archive/v2.8.0.tar.gz"
  sha256 "e04de57ab048f1abee75e9e739514c4f47e6cbb8acacb9d58a6e2892df30dc42"

  head "https://github.com/moul/advanced-ssh-config.git"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["GOBIN"] = buildpath
    ENV["GO15VENDOREXPERIMENT"] = "1"
    (buildpath/"src/github.com/moul/advanced-ssh-config").install Dir["*"]

    system "go", "build", "-o", "#{bin}/assh", "-v", "github.com/moul/advanced-ssh-config/cmd/assh/"

    bash_completion.install "src/github.com/moul/advanced-ssh-config/contrib/completion/bash_autocomplete"
    zsh_completion.install "src/github.com/moul/advanced-ssh-config/contrib/completion/zsh_autocomplete"
  end

  test do
    output = shell_output(bin/"assh --version")
    assert output.include? "assh version 2"
  end
end
