{
  description = "ansible-role - simple wrapper for Ansible to run a single role without the need to generate a playbook first";

  outputs = { self, nixpkgs }:
    let
      forAllSystems = nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
      ];

      pkgsForSystem = system: (import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      });
    in
    {
      overlays.default = _final: prev:
        let
          inherit (prev) buildGo122Module callPackage lib;
        in
        {
          ansible-role = callPackage ./package.nix { inherit buildGo122Module lib; };
        };

      packages = forAllSystems (system: rec {
        inherit (pkgsForSystem system) ansible-role;
        default = ansible-role;
      });
    };
}
