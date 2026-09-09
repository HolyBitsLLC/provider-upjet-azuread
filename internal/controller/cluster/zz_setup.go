// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	application "github.com/upbound/provider-azuread/v2/internal/controller/cluster/applications/application"
	password "github.com/upbound/provider-azuread/v2/internal/controller/cluster/applications/password"
	group "github.com/upbound/provider-azuread/v2/internal/controller/cluster/groups/group"
	member "github.com/upbound/provider-azuread/v2/internal/controller/cluster/groups/member"
	providerconfig "github.com/upbound/provider-azuread/v2/internal/controller/cluster/providerconfig"
	principal "github.com/upbound/provider-azuread/v2/internal/controller/cluster/serviceprincipals/principal"
	user "github.com/upbound/provider-azuread/v2/internal/controller/cluster/users/user"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		application.Setup,
		password.Setup,
		group.Setup,
		member.Setup,
		providerconfig.Setup,
		principal.Setup,
		user.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		application.SetupGated,
		password.SetupGated,
		group.SetupGated,
		member.SetupGated,
		providerconfig.SetupGated,
		principal.SetupGated,
		user.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
