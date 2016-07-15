require "language/go"

class Assh < Formula
  desc "assh: Advanced SSH config - A transparent wrapper that adds regex, aliases, gateways, includes, dynamic hostnames to SSH"
  homepage "https://github.com/moul/advanced-ssh-config"
  url "https://github.com/moul/advanced-ssh-config/archive/v2.4.0.tar.gz"
  sha256 "4929e2c3947b383fb86550f23528590e45b9d5166b1d001c9ef043d9e5c5fbf3"

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
