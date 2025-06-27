# SPDX-License-Identifier: Apache-2.0

%global goipath github.com/complytime/complyctl
%global base_url https://%{goipath}
%global gopath %{_builddir}/go

Name:           complyctl
Version:        0.0.6
Release:        %autorelease
Summary:        Tool to perform compliance assessment activities, scaled by plugins
License:        Apache-2.0
URL:            %{base_url}
Source0:        %{base_url}/archive/refs/tags/v%{version}.tar.gz

# git is temporarily used
BuildRequires:  git
BuildRequires:  golang
BuildRequires:  go-rpm-macros

%gometa -f

%description
%{name} leverages OSCAL to perform compliance assessment activities, using
plugins for each stage of the lifecycle.

%package        openscap-plugin
Summary:        A plugin which extends complyctl capabilities to use OpenSCAP
Requires:       %{name}%{?_isa} = %{version}-%{release}
Requires:       scap-security-guide
%description    openscap-plugin
openscap-plugin is a plugin which extends the complyctl capabilities to use
OpenSCAP. The plugin communicates with complyctl via gRPC, providing a
standard and consistent communication mechanism that gives independence for
plugin developers to choose their preferred languages.

%prep
%autosetup -n %{name}-%{version}

%build
BUILD_DATE_GO=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

# Set up environment variables and flags to build properly and securely
%set_build_flags

# Align GIT_COMMIT and GIT_TAG with version for simplicity
GO_LD_EXTRAFLAGS="-X %{goipath}/internal/version.version=%{version} \
                  -X %{goipath}/internal/version.gitTreeState=clean \
                  -X %{goipath}/internal/version.commit=%{version} \
                  -X %{goipath}/internal/version.buildDate=${BUILD_DATE_GO}"

# Adapt go env to RPM build environment
export GO111MODULE=on

# Define and create the output directory for binaries
GO_BUILD_BINDIR=./bin
mkdir -p ${GO_BUILD_BINDIR}

# Not calling the macro for more control on go env variables
go build -mod=vendor -o ${GO_BUILD_BINDIR}/ -ldflags="${GO_LD_EXTRAFLAGS}" ./cmd/...

%install
# Install complyctl directories
install -d %{buildroot}%{_bindir}
install -d -m 0755 %{buildroot}%{_datadir}/%{name}/{plugins,bundles,controls}
install -d -m 0755 %{buildroot}%{_libexecdir}/%{name}/plugins
install -d -m 0755 %{buildroot}%{_sysconfdir}/%{name}/config.d
install -d -m 0755 %{buildroot}%{_mandir}/{man1,man5}

# Copy sample data to be consumed by complyctl CLI
cp -rp docs/samples %{buildroot}%{_datadir}/%{name}

# Install files for complyctl CLI
install -p -m 0755 bin/complyctl %{buildroot}%{_bindir}/complyctl
install -p -m 0644 docs/man/complyctl.1 %{buildroot}%{_mandir}/man1/complyctl.1

# Install files for openscap-plugin package
install -p -m 0755 bin/openscap-plugin %{buildroot}%{_libexecdir}/%{name}/plugins/openscap-plugin
install -p -m 0644 docs/man/c2p-openscap-manifest.5 %{buildroot}%{_mandir}/man5/c2p-openscap-manifest.5

%check
# Run unit tests
go test -mod=vendor -race -v ./...

%files
%defattr(0644, root, root, 0755)
%attr(0755, root, root) %{_bindir}/complyctl
%license LICENSE
%{_mandir}/man1/complyctl.1*
%{_datadir}/%{name}/samples/{sample-catalog.json,sample-component-definition.json,sample-profile.json}
%{_datadir}/%{name}/{plugins,bundles,controls}
%{_sysconfdir}/%{name}/config.d

%files          openscap-plugin
%attr(0755, root, root) %{_libexecdir}/%{name}/plugins/openscap-plugin
%{_mandir}/man5/c2p-openscap-manifest.5*

%changelog
* Mon Jun 16 2025 George Vauter <gvauter@redhat.com>
- Update package name to complyctl

* Wed Jun 11 2025 Marcus Burghardt <maburgha@redhat.com>
- Bump to upstream version v0.0.6
- Align with Fedora Package Guidelines

* Tue May 6 2025 Qingmin Duanmu <qduanmu@redhat.com>
- Add complytime and openscap plugin man pages

* Wed Apr 30 2025 Qingmin Duanmu <qduanmu@redhat.com>
- Separate plugin binary from manifest

* Fri Apr 11 2025 Qingmin Duanmu <qduanmu@redhat.com>
- Separate package for openscap-plugin

* Tue Apr 08 2025 Marcus Burghardt <maburgha@redhat.com>
- Initial RPM
