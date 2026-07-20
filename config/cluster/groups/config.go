// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package groups

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("azuread_group", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "azuread"
		r.ShortGroup = "groups"
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"owners"},
		}
	})
	p.AddResourceConfigurator("azuread_group_member", func(r *config.Resource) {
		r.References["group_object_id"] = config.Reference{
			TerraformName: "azuread_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("object_id",true)`,
		}
		r.References["member_object_id"] = config.Reference{
			TerraformName: "azuread_user",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("object_id",true)`,
		}

		// Duplicate the member_object_id Terraform schema field as
		// group_member_object_id so that the code generator produces a
		// separate set of Ref/Selector fields that resolve against
		// Group resources (not just Users). This enables nested-group
		// membership (e.g., adding homelab-admins to homelab-family-access).
		if schema, ok := r.TerraformResource.Schema["member_object_id"]; ok {
			c := *schema // copy the schema
			r.TerraformResource.Schema["group_member_object_id"] = &c
		}

		// Configure reference for the new field → Group
		r.References["group_member_object_id"] = config.Reference{
			TerraformName: "azuread_group",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("object_id",true)`,
		}

		// Injector to merge group_member_object_id back into member_object_id
		// before the Terraform call. Both fields feed the same underlying
		// Azure AD member_object_id attribute.
		r.TerraformConfigurationInjector = consolidateGroupMemberObjectId

		// We need to override the default group that upjet generated for
		// this resource, which would be "azuread"
		r.ShortGroup = "groups"
	})
}

// consolidateGroupMemberObjectId merges the resolved group_member_object_id
// back into member_object_id before the Terraform call. This allows a Group
// Member to reference a Group (not just a User) via the
// groupMemberObjectIdSelector / groupMemberObjectIdRef fields.
func consolidateGroupMemberObjectId(jsonMap map[string]any, tfMap map[string]any) error {
	if tfMap["group_member_object_id"] != nil {
		tfMap["member_object_id"] = tfMap["group_member_object_id"]
		delete(tfMap, "group_member_object_id")
	}
	return nil
}
