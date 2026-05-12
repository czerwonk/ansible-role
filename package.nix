{ lib, rustPlatform }:

rustPlatform.buildRustPackage {
  pname = "ansible-role";
  version = "0.5.0";

  src = lib.cleanSource ./.;

  cargoLock.lockFile = ./Cargo.lock;

  meta = with lib; {
    description = "This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.";
    homepage = "https://github.com/czerwonk/ansible-role";
    license = licenses.mit;
  };
}
