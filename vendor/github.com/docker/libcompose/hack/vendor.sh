#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
rm -rf vendor/
source 'hack/.vendor-helpers.sh'

clone git github.com/Sirupsen/logrus v0.10.0
clone git github.com/urfave/cli v1.18.0
clone git github.com/docker/distribution 77b9d2997abcded79a5314970fe69a44c93c25fb
clone git github.com/vbatts/tar-split v0.10.1
clone git github.com/docker/docker 8658748ef716e43a5f6d834825d818012ed6e2c4
clone git github.com/docker/go-units f2145db703495b2e525c59662db69a7344b00bb8
clone git github.com/docker/go-connections 988efe982fdecb46f01d53465878ff1f2ff411ce
clone git github.com/flynn/go-shlex 3f9db97f856818214da2e1057f8ad84803971cff
clone git github.com/gorilla/context v1.1
clone git github.com/gorilla/mux v1.1
clone git github.com/opencontainers/runc cc29e3dded8e27ba8f65738f40d251c885030a28
clone git github.com/stretchr/testify a1f97990ddc16022ec7610326dd9bce31332c116
clone git github.com/davecgh/go-spew 5215b55f46b2b919f50a1df0eaa5886afe4e3b3d
clone git github.com/pmezard/go-difflib d8ed2627bdf02c080bf22230dbb337003b7aba2d
clone git golang.org/x/crypto 3fbbcd23f1cb824e69491a5930cfeff09b12f4d2 https://github.com/golang/crypto.git
clone git golang.org/x/net 2beffdc2e92c8a3027590f898fe88f69af48a3f8 https://github.com/tonistiigi/net.git
clone git golang.org/x/sys eb2c74142fd19a79b3f237334c7384d5167b1b46 https://github.com/golang/sys.git
clone git golang.org/x/time a4bde12657593d5e90d0533a3e4fd95e635124cb https://github.com/golang/time.git
clone git gopkg.in/check.v1 11d3bc7aa68e238947792f30573146a3231fc0f1
# clone git github.com/go-check/check 4ed411733c5785b40214c70bce814c3a3a689609 https://github.com/cpuguy83/check.git
clone git gopkg.in/yaml.v2 e4d366fc3c7938e2958e662b4258c7a89e1f0e3e
clone git github.com/Azure/go-ansiterm 388960b655244e76e24c75f48631564eaefade62
clone git github.com/Microsoft/go-winio v0.3.4
clone git github.com/xeipuuv/gojsonpointer e0fe6f68307607d540ed8eac07a342c33fa1b54a
clone git github.com/xeipuuv/gojsonreference e02fc20de94c78484cd5ffb007f8af96be030a45
clone git github.com/xeipuuv/gojsonschema ac452913faa25c08bb78810d3e6f88b8a39f8f25
clone git github.com/kr/pty 5cf931ef8f
clone git github.com/pkg/errors 01fa4104b9c248c8945d14d9f128454d5b28d595

clone git github.com/spf13/pflag cb88ea77998c3f024757528e3305022ab50b43be

clean && mv vendor/src/* vendor
