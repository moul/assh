require "language/go"

class Assh < Formula
  desc "assh: Advanced SSH config - A transparent wrapper that adds regex, aliases, gateways, includes, dynamic hostnames to SSH"
  homepage "https://github.com/moul/advanced-ssh-config"
  url "https://github.com/moul/advanced-ssh-config/archive/v2.1.0.tar.gz"
  sha256 "c6566457f1734afe1cc83740d53705a7aa5ef7e775bc429561901b7c9d48d6e1"

  head "https://github.com/moul/advanced-ssh-config.git"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["CGO_ENABLED"] = "0"
    ENV.prepend_create_path "PATH", buildpath/"bin"

    mkdir_p buildpath/"src/github.com/moul"
    ln_s buildpath, buildpath/"src/github.com/moul/advanced-ssh-config"
    Language::Go.stage_deps resources, buildpath/"src"

    # FIXME: update version
    system "go", "build", "-o", "assh", "./cmd/assh"
    bin.install "assh"

    # FIXME: add autocompletion
  end

  test do
    output = shell_output(bin/"assh --version")
    assert output.include? "assh version 2"
  end
end
