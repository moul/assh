require "language/go"

class Assh < Formula
  desc "assh: Advanced SSH config - A transparent wrapper that adds regex, aliases, gateways, includes, dynamic hostnames to SSH"
  homepage "https://github.com/moul/advanced-ssh-config"
  url "https://github.com/moul/advanced-ssh-config/archive/v2.4.1.tar.gz"
  sha256 "8867df447e654dae384cf598a81eb6be57a49082449ef4387f33725bb216853c"

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

  def caveats
    <<-EOS.undent
    To activate advanced pattern matching, add the following at the end of your .bashrc or .zshrc:

      alias ssh="assh wrapper ssh"
    EOS
  end

  test do
    output = shell_output(bin/"assh --version")
    assert output.include? "assh version 2"
  end
end
