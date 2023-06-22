{
  inputs = {
    nixpkgs.url = "flake:nixpkgs/master";

    utils.url = "flake:flake-utils";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.utils.follows = "utils";
    };
  };

  outputs =
    { self
    , nixpkgs
    , utils
    , gomod2nix
    , ...
    }: utils.lib.eachSystem [
      "aarch64-darwin"
      "x86_64-darwin"
      "aarch64-linux"
      "x86_64-linux"
    ]
      (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            gomod2nix.overlays.default
          ];
        };

        lib = pkgs.lib;

        stdenv = pkgs.stdenv;

        hopper = pkgs.buildGoApplication rec {
          name = "hopper";

          src = lib.cleanSource ./.;

          modules = ./gomod2nix.toml;
        };
      in
      {
        packages.default = hopper;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gomod2nix.packages.${system}.default

            llvmPackages.clangUseLLVM
            llvmPackages.bintools-unwrapped
          ];
        };

        formatter = pkgs.nixpkgs-fmt;
      }
      );
}
