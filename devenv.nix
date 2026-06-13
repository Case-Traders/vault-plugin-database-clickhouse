{ pkgs, config, lib, ... }:
let
  hasDocker = builtins.pathExists "/var/run/docker.sock";
  containerCli = if hasDocker then pkgs.docker else pkgs.podman;

  opamProofInputs = with pkgs; [
    opam
    pkg-config
    gmp
    zlib
    ncurses
    linuxHeaders
    findutils
    which
    patch
    m4
    gcc
    git
    python3
  ];

  configureRuntime = ''
    if [[ -z "''${DOCKER_HOST:-}" ]]; then
      if [[ -S /var/run/docker.sock ]]; then
        export DOCKER_HOST=unix:///var/run/docker.sock
      elif [[ -S "''${XDG_RUNTIME_DIR:-/run/user/$(id -u)}/podman/podman.sock" ]]; then
        export DOCKER_HOST=unix://''${XDG_RUNTIME_DIR}/podman/podman.sock
        export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE="''${XDG_RUNTIME_DIR}/podman/podman.sock"
        export TESTCONTAINERS_RYUK_DISABLED=true
      fi
    fi
  '';

  opamEnvHook = ''
    if [[ -d "$OPAMROOT/$OPAMSWITCH" ]]; then
      eval "$(opam env --switch="$OPAMSWITCH" 2>/dev/null)" || true
      unset out dev 2>/dev/null || true
    fi
  '';

  proofSetup = ''
    set -euo pipefail
    if command -v rocq >/dev/null \
      && test -f "$PERENNIAL_ROOT/new/proof/proof_prelude.vo"; then
      echo "proof-setup: already done ($PERENNIAL_ROOT)"
      exit 0
    fi
    ${configureRuntime}
    ${opamEnvHook}

    opam_install=(opam install --confirm-level=unsafe-yes --assume-depexts -y)
    jobs="''${OPAMJOBS:-$(nproc 2>/dev/null || echo 4)}"

    if [[ ! -d "$OPAMROOT" ]]; then
      opam init --bare --disable-sandboxing --no-setup
    fi

    if ! opam switch list --short 2>/dev/null | grep -qx "$OPAMSWITCH"; then
      opam switch create "$OPAMSWITCH" ocaml-base-compiler
    fi
    eval "$(opam env --switch="$OPAMSWITCH")"
    opam option depext=false 2>/dev/null || true
    unset out dev 2>/dev/null || true
    export PATH="$(opam var bin):$PATH"

    if ! opam repo list -s | grep -qx rocq-released; then
      opam repo add rocq-released --all-switches --set-default https://rocq-prover.org/opam/released
    fi

    if ! command -v rocq >/dev/null; then
      echo ">>> Installing Rocq (rocq-prover)..."
      "''${opam_install[@]}" conf-gmp conf-pkg-config conf-linux-libc-dev
      "''${opam_install[@]}" rocq-prover
    fi

    mkdir -p "$(dirname "$PERENNIAL_ROOT")"
    if [[ ! -d "$PERENNIAL_ROOT/.git" ]]; then
      echo ">>> Cloning Perennial..."
      git clone --depth 1 https://github.com/mit-pdos/perennial.git "$PERENNIAL_ROOT"
    fi
    echo ">>> Checking out Perennial $PERENNIAL_PIN..."
    git -C "$PERENNIAL_ROOT" fetch --depth 1 origin "$PERENNIAL_PIN"
    git -C "$PERENNIAL_ROOT" checkout -f "$PERENNIAL_PIN"

    cd "$PERENNIAL_ROOT"
    opam pin add -n . --no-action 2>/dev/null || true
    echo ">>> Installing Perennial opam deps (rocq-iris builds from source — expect 30–60 min, CPU busy, little output)..."
    export OPAMJOBS="''${jobs}"
    "''${opam_install[@]}" --verbose --deps-only ./perennial.opam
    echo ">>> Building Perennial proof_prelude.vo..."
    make -j"''${jobs}" new/proof/proof_prelude.vo

    mkdir -p "$DEVENV_ROOT/.cache"
    date -Iseconds >"$DEVENV_ROOT/.cache/proof-setup.stamp"
    echo "proof-setup: PERENNIAL_ROOT=$PERENNIAL_ROOT switch=$OPAMSWITCH"
  '';
in
{
  env = {
    PLUGIN_DIR = config.env.DEVENV_ROOT;
    GOPATH = lib.mkDefault "${config.env.DEVENV_ROOT}/.go";
    GOBIN = lib.mkDefault "${config.env.DEVENV_ROOT}/.go/bin";
    GOMODCACHE = lib.mkDefault "${config.env.DEVENV_ROOT}/.go/pkg/mod";
    GOTOOLCHAIN = lib.mkForce "go1.25.11";
    OPAMROOT = lib.mkDefault "${config.env.DEVENV_ROOT}/.opam";
    OPAMSWITCH = lib.mkDefault "perennial-proof";
    OPAMDISABLESANDBOXING = "1";
    OPAMASSUMESDEPEXTNODEPS = "1";
    OPAMCONFIRMLEVEL = "unsafe-yes";
    OPAMDEPEXT = "disable";
    OPAMYES = "1";
    PERENNIAL_ROOT = lib.mkDefault "${config.env.DEVENV_ROOT}/.cache/perennial";
    PERENNIAL_PIN = lib.mkDefault "c15c19774d4394959ae1e9ee85e5852df00046e7";
    PKG_CONFIG_PATH = lib.makeSearchPath "lib/pkgconfig" [
      pkgs.gmp
      pkgs.zlib
    ];
    CPATH = lib.concatStringsSep ":" [
      "${pkgs.gmp.dev}/include"
      "${pkgs.zlib.dev}/include"
      "${pkgs.ncurses.dev}/include"
      "${pkgs.linuxHeaders}/include"
    ];
    LIBRARY_PATH = lib.makeLibraryPath [
      pkgs.gmp
      pkgs.zlib
      pkgs.ncurses
    ];
    LD_LIBRARY_PATH = lib.makeLibraryPath [
      pkgs.gmp
      pkgs.zlib
      pkgs.ncurses
    ];
  };

  languages.go = {
    enable = true;
    package = pkgs.go_1_25;
  };

  cachix.enable = false;

  devcontainer.enable = true;
  devcontainer.settings.customizations.vscode = {
    extensions = [
      "datakurre.devenv"
      "ejgallego.coq-lsp"
    ];
  };

  files.".vscode/extensions.json".json = {
    recommendations = [ "ejgallego.coq-lsp" ];
  };

  files.".vscode/settings.json".text = ''
    {
      "go.goroot": "${pkgs.go_1_25}/share/go",
      "go.alternateTools": {
        "go": "${pkgs.go_1_25}/bin/go",
        "gopls": "${pkgs.gopls}/bin/gopls",
        "staticcheck": "${pkgs.go-tools}/bin/staticcheck"
      },
      "go.toolsEnvVars": {
        "GOPATH": "${config.env.DEVENV_ROOT}/.go",
        "GOMODCACHE": "${config.env.DEVENV_ROOT}/.go/pkg/mod"
      },
      "go.toolsManagement.autoUpdate": false,
      "go.lintTool": "staticcheck"
    }
  '';

  packages = with pkgs; [
    go_1_25
    go-tools
    gopls
    gofumpt
    golangci-lint
    containerCli
    gnumake
    python3Packages.mkdocs
    python3Packages.mkdocs-material
    goreleaser
    cosign
    syft
    gh
  ] ++ opamProofInputs;

  processes = lib.mkIf (!hasDocker) {
    podman = {
      exec = "podman system service --time=0 unix://$XDG_RUNTIME_DIR/podman/podman.sock";
      ready.exec = "test -S \"$XDG_RUNTIME_DIR/podman/podman.sock\"";
    };
  };

  enterShell = ''
    ${configureRuntime}
    ${opamEnvHook}
    export PATH=$GOBIN:$PATH
    mkdir -p $GOPATH $GOBIN $GOMODCACHE
    if ! command -v gomarkdoc >/dev/null 2>&1; then
      go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
    fi
    make -C "''${DEVENV_ROOT}" goose-tools 2>/dev/null || true
    if [[ -n "''${DOCKER_HOST:-}" ]]; then
      if [[ "''${DOCKER_HOST}" == *podman* ]]; then
        echo "container runtime: podman (''${DOCKER_HOST})"
      else
        echo "container runtime: docker (''${DOCKER_HOST})"
      fi
    else
      echo "container runtime: none (start docker or: devenv up)"
    fi
    echo "opam: $(opam --version 2>/dev/null || echo missing)"
    if command -v rocq >/dev/null 2>&1; then
      echo "rocq: $(rocq --version 2>/dev/null | head -1)"
    else
      echo "rocq: not installed (run: ch-proof-setup)"
    fi
    echo "perennial: $PERENNIAL_ROOT"
  '';

  scripts = {
    ch-build.exec = ''
      cd "$PLUGIN_DIR" && make build-linux-amd64
    '';
    ch-test.exec = ''
      cd "$PLUGIN_DIR" && make test
    '';
    ch-proof.exec = ''
      cd "$PLUGIN_DIR" && make proof
    '';
    ch-proof-setup.exec = proofSetup;
    ch-ci.exec = ''
      cd "$PLUGIN_DIR" && make ci
    '';
    ch-test-integration.exec = ''
      ${configureRuntime}
      cd "$PLUGIN_DIR" && make test-integration
    '';
    ch-integration.exec = ''
      ${configureRuntime}
      cd "$PLUGIN_DIR" && make ci-integration
    '';
    ch-docs.exec = ''
      cd "$PLUGIN_DIR" && make docs
    '';
    ch-docs-serve.exec = ''
      cd "$PLUGIN_DIR" && make docs-serve
    '';
    ch-release-snapshot.exec = ''
      cd "$PLUGIN_DIR" && goreleaser release --snapshot --clean --skip=sign
    '';
  };
}
