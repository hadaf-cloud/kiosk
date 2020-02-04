/*
Copyright 2020 DevSpace Technologies Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
)

var (

	// UserPrefix is used to prefix users for the index
	UserPrefix = "user:"
	// GroupPrefix is used to prefix groups for the index
	GroupPrefix = "group:"
)

// ConvertSubject converts the given subject into an unqiue id string
func ConvertSubject(namespace string, subject *rbacv1.Subject) string {
	// if subject.APIGroup != rbacv1.GroupName {
	//	return ""
	// }

	switch subject.Kind {
	case rbacv1.UserKind:
		return UserPrefix + subject.Name

	case rbacv1.GroupKind:
		return GroupPrefix + subject.Name

	case rbacv1.ServiceAccountKind:
		saNamespace := namespace
		if len(subject.Namespace) > 0 {
			saNamespace = subject.Namespace
		}
		if len(saNamespace) == 0 {
			return ""
		}

		return UserPrefix + serviceaccount.MakeUsername(saNamespace, subject.Name)
	default:
	}

	return ""
}
