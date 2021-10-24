{
  description = "Withings Go API client";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;

          /* config.allowUnfree = true; */
          config.allowUnfreePredicate = pkg: builtins.elem (nixpkgs.lib.getName pkg) [
            "ngrok"
          ];
        };
        #pkgs = nixpkgs.legacyPackages.${system};
        buildDeps = with pkgs; [ git go gnumake ];
        devDeps = with pkgs; buildDeps ++ [ golangci-lint gotestsum ngrok ];

        generateGoEnv = go:
          pkgs.buildEnv {
            name = "go" + go.version;
            paths = (pkgs.lib.remove pkgs.go devDeps) ++ [ go ];
          };
      in {
        devShell = pkgs.mkShell {
          buildInputs = devDeps;
        };

        packages.go1_15 = generateGoEnv pkgs.go_1_15;
        packages.go1_16 = generateGoEnv pkgs.go_1_16;
        packages.go1_17 = generateGoEnv pkgs.go_1_17;
      });
}
