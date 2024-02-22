{ lib, buildGoModule }:

buildGoModule {
  pname = "ansible-role";
  version = "0.4.4";

  src = lib.cleanSource ../.;

  vendorHash = null;

  CGO_ENABLED = 0;

  meta = with lib; {
    description = "This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.";
    homepage = "https://github.com/czerwonk/ansible-role";
    license = licenses.mit;
  };
}
