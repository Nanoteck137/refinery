{ self }: 
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.refinery;

  refineryConfig = pkgs.writeText "config.toml" ''
    listen_addr = "${cfg.host}:${toString cfg.port}"
    data_dir = "${cfg.dataDir}"
    library_dir = "${cfg.libraryDir}"
    username = "${cfg.username}"
    initial_password = "${cfg.initialPassword}"
    jwt_secret = "${cfg.jwtSecret}"
  '';
in
{
  options.services.refinery = {
    enable = mkEnableOption "Enable the refinery service";

    port = mkOption {
      type = types.port;
      default = 7550;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    dataDir = mkOption {
      type = types.path;
      default = "/var/lib/refinery";
      description = "path to the data directory";
    };

    username = mkOption {
      type = types.str;
      description = "username of the first user";
    };

    initialPassword = mkOption {
      type = types.str;
      description = "initial password of the first user (should change after the first login)";
    };

    jwtSecret = mkOption {
      type = types.str;
      description = "jwt secret";
    };

    libraryDir = mkOption {
      type = types.path;
      description = "library directory";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.backend;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "refinery";
      description = "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "refinery";
      description = "group to use for this service";
    };

    openFirewall = mkOption {
      type = types.bool;
      default = false;
      description = "open the ports in the firewall";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.refinery = {
      description = "refinery";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      serviceConfig = mkMerge [
        {
          User = cfg.user;
          Group = cfg.group;

          ExecStart = "${cfg.package}/bin/refinery serve -c '${refineryConfig}'";

          Restart = "on-failure";
          RestartSec = "5s";

          PrivateTmp = true;
          ProtectHome = true;
          ProtectHostname = true;
          ProtectKernelLogs = true;
          ProtectKernelModules = true;
          ProtectKernelTunables = true;
          ProtectProc = "invisible";
          ProtectSystem = "strict";
          RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
          RestrictNamespaces = true;
          RestrictRealtime = true;
          RestrictSUIDSGID = true;
        }

        (mkIf (cfg.dataDir != "/var/lib/refinery") {
          ReadWritePaths = [ cfg.dataDir ];
        })

        (mkIf (cfg.dataDir == "/var/lib/refinery") {
          StateDirectory = "refinery";
        })

      ];
    };

    networking.firewall = lib.mkIf cfg.openFirewall {
      allowedTCPPorts = [ cfg.port ];
    };

    users.users = mkIf (cfg.user == "refinery") {
      refinery = {
        group = cfg.group;
        isSystemUser = true;
        home = "${cfg.dataDir}";
      };
    };

    users.groups = mkIf (cfg.group == "refinery") {
      refinery = {};
    };
  };
}
