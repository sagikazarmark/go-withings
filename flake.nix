{
  description = "Withings Go API client";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        buildDeps = with pkgs; [ git go gnumake go-task ];
        devDeps = with pkgs; buildDeps ++ [ gotestsum golangci-lint ];

        goShell = go:
          pkgs.mkShell {
            buildInputs = (pkgs.lib.remove pkgs.go devDeps) ++ [ go ];
          };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = devDeps;
        };

        devShells.go1_15 = goShell pkgs.go_1_15;
        devShells.go1_16 = goShell pkgs.go_1_16;
        devShells.go1_17 = goShell pkgs.go_1_17;
      });
}
