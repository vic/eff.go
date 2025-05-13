{
  inputs = {
    systems.url = "github:nix-systems/default";
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devshell.url = "github:numtide/devshell";
    devshell.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs: inputs.flake-parts.lib.mkFlake { inherit inputs; } {
    systems = import inputs.systems;
    imports = [
      inputs.devshell.flakeModule
    ];

    perSystem = {pkgs, ...}: {
      devshells.default = { 
        devshell.packages = [ pkgs.go pkgs.gopls ];
        commands = [ 
          { package = pkgs.mdbook; } 
        ];
      };
    };
  };
}