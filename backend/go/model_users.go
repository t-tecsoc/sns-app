/*
 * sns-app
 *
 * SNSアプリ
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Users struct {
	Id string `json:"id"`

	HandleName string `json:"handle_name"`

	DisplayName string `json:"display_name"`

	EncryptedEmailAddress string `json:"encrypted_email_address"`

	EncryptedPassword string `json:"encrypted_password"`
}
