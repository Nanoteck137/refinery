{
  description = "Devshell for refinery";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        app = pkgs.buildGoModule {
          pname = "refinery";
          version = fullVersion;
          src = ./.;

          # ldflags = [
          #   "-X github.com/nanoteck137/refinery.Version=${version}"
          #   "-X github.com/nanoteck137/refinery.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          # ];

          vendorHash = "sha256-GBLUL8rs4rRysi7LKmlCVaceBkq11KnHwg4vwprzWTg=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/refinery --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.nix pkgs.git pkgs.skopeo ]}
          '';
        };
      in
      {
        packages = {
          default = app;
          inherit app;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            skopeo
            git
          ];
        };
      }
    );
}
