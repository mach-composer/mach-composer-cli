{
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
        };
        mach-composer = import ./. { inherit pkgs; };
      in {
        packages = {
          default = mach-composer;
        };
      }) // {
        overlays.default = final: prev: {
          mach-composer = import ./. { pkgs = final; };
        };
      };
}