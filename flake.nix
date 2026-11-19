{
  description = "Brotchaschte Nix Flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];

      perSystem = {pkgs, ...}: let
        pnpm = pkgs.pnpm_12;

        dependencies = with pkgs; [go go-task air sqlc zig pnpm nodejs];

        brotchaschte = pkgs.stdenv.mkDerivation (finallAttrs: {
          pname = "brotchaschte";
          version = "1.0";
          src = ./.;
          nativeBuildInputs = dependencies ++ [pkgs.pnpmConfigHook];

          pnpmRoot = "frontend";
          pnpmDeps = pkgs.fetchPnpmDeps {
            inherit (finallAttrs) pname version;
            inherit pnpm;

            src = ./frontend;

            fetcherVersion = 5;
            hash = "sha256-DZwi2vZ8sm1FWl0RX13xaO78IP5Vd9x6V/K/vmL3LyU=";
          };

          buildPhase = ''
            runHook preBuild

            task build

            runHook postBuild
          '';

          installPhase = ''
            mkdir -p $out/bin
            install -m755 ./main $out/bin/brotchaschte
          '';
        });
      in rec {
        packages.default = brotchaschte;
        packages.brotchaschte = brotchaschte;

        apps.default = {
          type = "app";
          program = "${brotchaschte}/bin/brotchaschte";
        };
        apps.brotchaschte = apps.default;

        devShells.default = pkgs.mkShell {
          nativeBuildInputs = dependencies;
        };
      };
    };
}
