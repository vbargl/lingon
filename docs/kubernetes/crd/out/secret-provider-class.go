// Code generated by lingon. EDIT AS MUCH AS YOU LIKE.

package team

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
)

var TeamOnboardingAuthSecretProviderClass = &v1.SecretProviderClass{
	ObjectMeta: metav1.ObjectMeta{Name: "team-onboarding-auth"},
	Spec: v1.SecretProviderClassSpec{
		Parameters: map[string]string{
			"objects": `
- objectName: "xxx-gh-token"
  secretPath: "team-onboarding-kv/data/github-xxx-bot"
  secretKey: "token"
- objectName: "xxx-gh-username"
  secretPath: "team-onboarding-kv/data/github-xxx-bot"
  secretKey: "username"
- objectName: "harbor-username"
  secretPath: "team-onboarding-kv/data/harbor"
  secretKey: "username"
- objectName: "harbor-password"
  secretPath: "team-onboarding-kv/data/harbor"
  secretKey: "password"
- objectName: "sendgrid-api-key"
  secretPath: "team-onboarding-kv/data/sendgrid-api-key"
  secretKey: "token"
- objectName: "key-id"
  secretPath: "team-onboarding-kv/data/lakefs-xxx-admin"
  secretKey: "access-key-id"
- objectName: "secret-key"
  secretPath: "team-onboarding-kv/data/lakefs-xxx-admin"
  secretKey: "secret-access-key"
- objectName: "abk-vcc-test-xxx-xxx-gh-dev-token"
  secretPath: "team-onboarding-kv/data/abk-vcc-test-xxx-bot"
  secretKey: "token"
- objectName: "abk-vcc-test-xxx-gh-dev-username"
  secretPath: "team-onboarding-kv/data/abk-vcc-test-xxx-bot"
  secretKey: "username"
- objectName: "scim-auth-token"
  secretPath: "team-onboarding-kv/data/scim-token"
  secretKey: "token"
- objectName: "team-az-group-management-prod"
  secretPath: "team-onboarding-kv/data/team-az-group-management-prod"
  secretKey: "client-secret"
- objectName: "team-az-group-management-qa"
  secretPath: "team-onboarding-kv/data/team-az-group-management-qa"
  secretKey: "client-secret"

`,
			"roleName":     "team-onboarding-policy-read",
			"vaultAddress": "https://vault.secretstore.company.com",
		},
		Provider: v1.Provider("vault"),
		SecretObjects: []*v1.SecretObject{{
			Data: []*v1.SecretObjectData{{
				Key:        "token",
				ObjectName: "xxx-gh-token",
			}, {
				Key:        "username",
				ObjectName: "xxx-gh-username",
			}},
			SecretName: "github-auth",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "token",
				ObjectName: "abk-vcc-test-xxx-xxx-gh-dev-token",
			}, {
				Key:        "username",
				ObjectName: "abk-vcc-test-xxx-gh-dev-username",
			}},
			SecretName: "github-abk-vcc-test-auth",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "token",
				ObjectName: "scim-auth-token",
			}},
			SecretName: "scim-auth",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "username",
				ObjectName: "harbor-username",
			}, {
				Key:        "password",
				ObjectName: "harbor-password",
			}},
			SecretName: "harbor-auth",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "access-key-id",
				ObjectName: "key-id",
			}, {
				Key:        "secret-access-key",
				ObjectName: "secret-key",
			}},
			SecretName: "lakefs-auth",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "token",
				ObjectName: "sendgrid-api-key",
			}},
			SecretName: "sendgrid-api-key",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "client-secret",
				ObjectName: "team-az-group-management-prod",
			}},
			SecretName: "team-az-group-management-prod",
			Type:       "Opaque",
		}, {
			Data: []*v1.SecretObjectData{{
				Key:        "client-secret",
				ObjectName: "team-az-group-management-qa",
			}},
			SecretName: "team-az-group-management-qa",
			Type:       "Opaque",
		}},
	},
	TypeMeta: metav1.TypeMeta{
		APIVersion: "secrets-store.csi.x-k8s.io/v1",
		Kind:       "SecretProviderClass",
	},
}
