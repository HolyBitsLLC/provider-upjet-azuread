// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import "github.com/crossplane/upjet/v2/pkg/config"

// terraformPluginSDKExternalNameConfigs contains all external name configurations for this
// provider.
//
// NOTE (HolyBitsLLC): this map doubles as the upjet include-list. Only the
// resources listed here are generated into Go APIs, controllers, and package
// CRDs (see `resourceList` in provider.go / provider_namespaced.go). We
// intentionally ship ONLY the managed resources actually used by
// harvester-argo (see clusters/krubecible/components/crossplane-microsoft/),
// to keep the CRD surface minimal.
var terraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{
	// applications
	//
	// azuread_application can be imported using their object ID
	"azuread_application": config.IdentifierFromProvider,
	// No import documented
	"azuread_application_password": config.IdentifierFromProvider,

	// groups
	//
	// azuread_group can be imported using their object ID
	"azuread_group": config.IdentifierFromProvider,
	// azuread_group_member can be imported using the object ID of the group and the object ID of the member:
	// {GroupObjectID}/member/{MemberObjectID}
	// 00000000-0000-0000-0000-000000000000/member/11111111-1111-1111-1111-111111111111
	"azuread_group_member": config.IdentifierFromProvider,

	// users
	//
	// azuread_user can be imported using their object ID
	"azuread_user": config.IdentifierFromProvider,

	// serviceprincipals
	//
	// azuread_service_principal can be imported using their object ID
	"azuread_service_principal": config.IdentifierFromProvider,
}

// cliReconciledExternalNameConfigs contains all external name configurations
// belonging to Terraform resources to be reconciled under the CLI-based
// architecture for this provider.
var cliReconciledExternalNameConfigs = map[string]config.ExternalName{}

// resourceConfigurator applies all external name configs
// listed in the table terraformPluginSDKExternalNameConfigs and
// cliReconciledExternalNameConfigs and sets the version
// of those resources to v1beta1. For those resource in
// terraformPluginSDKExternalNameConfigs, it also sets
// config.Resource.UseNoForkClient to `true`.
func resourceConfigurator() config.ResourceOption {
	return func(r *config.Resource) {
		// if configured both for the no-fork and CLI based architectures,
		// no-fork configuration prevails
		e, configured := terraformPluginSDKExternalNameConfigs[r.Name]
		if !configured {
			e, configured = cliReconciledExternalNameConfigs[r.Name]
		}
		if !configured {
			return
		}
		r.Version = "v1beta1"
		r.ExternalName = e
	}
}
