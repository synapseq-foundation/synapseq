PackageIdentifier: ruanklein.synapseq
PackageVersion: __VERSION__
InstallerType: zip
NestedInstallerType: portable
Installers:
  - Architecture: x64
    InstallerUrl: __RELEASE_BASE_URL__/synapseq-v__VERSION__-windows-amd64.zip
    InstallerSha256: __WINDOWS_AMD64_SHA__
    NestedInstallerFiles:
      - RelativeFilePath: synapseq-v__VERSION__-windows-amd64.exe
        PortableCommandAlias: synapseq
  - Architecture: arm64
    InstallerUrl: __RELEASE_BASE_URL__/synapseq-v__VERSION__-windows-arm64.zip
    InstallerSha256: __WINDOWS_ARM64_SHA__
    NestedInstallerFiles:
      - RelativeFilePath: synapseq-v__VERSION__-windows-arm64.exe
        PortableCommandAlias: synapseq
ManifestType: installer
ManifestVersion: 1.6.0
