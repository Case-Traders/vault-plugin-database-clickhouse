{ pkgs, config, lib, ... }:
let
  coq = pkgs.coq_8_18;
  coqLsp = pkgs.coqPackages_8_18.coq-lsp;
  hasDocker = builtins.pathExists "/var/run/docker.sock";

  # Prefer host Docker; otherwise expose Podman (and its API socket process).
  containerCli = if hasDocker then pkgs.docker else pkgs.podman;

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
in
{
  env = {
    PLUGIN_DIR = config.env.DEVENV_ROOT;
    GOPATH = lib.mkDefault "${config.env.DEVENV_ROOT}/.go";
    GOBIN = lib.mkDefault "${config.env.DEVENV_ROOT}/.go/bin";
    GOMODCACHE = lib.mkDefault "${config.env.DEVENV_ROOT}/.go/pkg/mod";
    GOTOOLCHAIN = lib.mkForce "go1.25.11";
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
    settings = {
      "coq-lsp.path" = "${coqLsp}/bin/coq-lsp";
    };
  };

  files.".vscode/extensions.json".json = {
    recommendations = [ "ejgallego.coq-lsp" ];
  };

  files.".vscode/settings.json".text = ''
    {
      "coq-lsp.path": "${coqLsp}/bin/coq-lsp",
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
    coq
    coqLsp
    containerCli
    gnumake
    python3
    python3Packages.mkdocs
    python3Packages.mkdocs-material
    goreleaser
    cosign
    syft
    gh
  ];

  # When Docker is absent, start a Podman API socket for testcontainers.
  processes = lib.mkIf (!hasDocker) {
    podman = {
      exec = "podman system service --time=0 unix://$XDG_RUNTIME_DIR/podman/podman.sock";
      ready.exec = "test -S \"$XDG_RUNTIME_DIR/podman/podman.sock\"";
    };
  };

  enterShell = ''
    ${configureRuntime}
    export PATH=${coq}/bin:${coqLsp}/bin:$PATH:$GOBIN
    mkdir -p $GOPATH $GOBIN $GOMODCACHE
    if ! command -v gomarkdoc >/dev/null 2>&1; then
      go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
    fi
    if [[ -n "''${DOCKER_HOST:-}" ]]; then
      if [[ "''${DOCKER_HOST}" == *podman* ]]; then
        echo "container runtime: podman (''${DOCKER_HOST})"
      else
        echo "container runtime: docker (''${DOCKER_HOST})"
      fi
    else
      echo "container runtime: none (start docker or: devenv up)"
    fi
    echo "coq: $(coqc --version | head -1)"
    echo "coq-lsp: $(coq-lsp --version 2>/dev/null || echo ${coqLsp}/bin/coq-lsp)"
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
