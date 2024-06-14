{ pkgs, ... }: pkgs.buildGoModule {
  name = "mach-composer";
  # TODO: tests write to config file. enable tests somehow?
  doCheck = false;
  # TODO: filter build dirs
  src = (pkgs.lib.cleanSource ./.);
  # TODO: generate vendor hash
  vendorHash = "sha256-SdeLKcQMR8Fwc64U4XRfoyWlmOVqonMhiXu13HAuEDI=";
  nativeBuildInputs = [
    pkgs.git
    pkgs.installShellFiles
    pkgs.terraform
  ];
  meta = {
    description = "Orchestration tool for modern MACH ecosystems, powered by Terraform infrastructure-as-code underneath";
    license = pkgs.lib.licenses.mit;
    mainProgram = "mach-composer";
  };
  postInstall = ''
    mv $out/bin/mach-composer-cli $out/bin/mach-composer

    installShellCompletion --cmd mach-composer \
      --bash <($out/bin/mach-composer completion bash) \
      --fish <($out/bin/mach-composer completion fish) \
      --zsh <($out/bin/mach-composer completion zsh)
  '';
}