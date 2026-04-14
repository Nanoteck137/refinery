{
  description = "Devshell for refinery";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, devtools, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        backend = pkgs.buildGoModule {
          pname = "refinery";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/refinery" "cmd/refinery-cli"];

          ldflags = [
            "-X github.com/nanoteck137/refinery.Version=${version}"
            "-X github.com/nanoteck137/refinery.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "sha256-0bNFL1iIrtG+siWbHu3ipTX3U9essxR7xnaR7BpAUyw=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/refinery --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
            wrapProgram $out/bin/refinery-cli --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
          '';
        };

        web = pkgs.buildNpmPackage {
          name = "refinery-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "sha256-xIrTtTfnssHyyOF8ruqfeEdcXrzKHENy5fL06B2vZD8=";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.dirtyRev or self.rev or "no-commit";
          PUBLIC_API_ADDRESS="";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            runHook postInstall
          '';
        };

        docker = pkgs.dockerTools.buildLayeredImage {
          name = "refinery";
          tag  = fullVersion;

          contents = [
            pkgs.dockerTools.caCertificates
            web
            backend
          ];

          extraCommands = ''
            mkdir -p tmp
            chmod 1777 tmp

            mkdir -p data
            mkdir -p media
          '';

          config = {
            Entrypoint   = [ "/bin/refinery" ];
            Cmd = [ "serve" ];
            ExposedPorts = { "3000/tcp" = {}; };
            Env = [
              "DWEBBLE_WEB=${web}"
            ];
          };
        };

        tools = devtools.packages.${system};
      in
      {
        packages = {
          default = backend;
          inherit backend web docker;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs
            imagemagick
            ffmpeg

            tools.publishVersion
          ];
        };
      }
    ) // {
      nixosModules.backend = import ./nix/backend.nix { inherit self; };
      nixosModules.default = { ... }: {
        imports = [
          self.nixosModules.backend
        ];
      };
    };
}
