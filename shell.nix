
{ pkgs ? import <nixpkgs> {
  overlays = [ (self: super: {
    # you can use this block to override with specific versions
    # nodejs = super.nodejs-10_x;
    # jre = super.jdk11;
  }) ];
} }:

pkgs.mkShell {
  name="dev-environment";
  buildInputs = [
   pkgs.go_1_19
   pkgs.golangci-lint
   pkgs.buf
   pkgs.protobuf
   pkgs.protoc-gen-go
   pkgs.gnumake
   pkgs.direnv
  ];
  shellHook = ''
    echo "Welcome to your dev env"
    '';
    }

