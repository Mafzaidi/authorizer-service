package entity

// Claims represents the JWT token claims structure for the domain layer.
// This entity contains all the information needed to identify and authorize a user.
type Claims struct {
	// Issuer identifies the principal that issued the JWT (iss claim)
	Issuer string `json:"iss"`

	// Subject identifies the principal that is the subject of the JWT (sub claim)
	Subject string `json:"sub"`

	// Audience identifies the recipients that the JWT is intended for (aud claim)
	Audience []string `json:"aud"`

	// ExpiresAt identifies the expiration time on or after which the JWT must not be accepted (exp claim)
	ExpiresAt int64 `json:"exp"`

	// IssuedAt identifies the time at which the JWT was issued (iat claim)
	IssuedAt int64 `json:"iat"`

	// Username is the username of the authenticated user
	Username string `json:"username"`

	// Email is the email address of the authenticated user
	Email string `json:"email"`

	// Authorization contains the authorization information for the user across different applications
	Authorization []Authorization `json:"authorization"`
}

// Authorization represents the authorization information for a specific application.
// It contains the roles and permissions that a user has for a particular application.
type Authorization struct {
	// App is the application code that this authorization applies to
	App string `json:"app"`

	// Roles is the list of role codes assigned to the user for this application
	Roles []string `json:"roles"`

	// Permissions is the list of permission codes granted to the user for this application
	Permissions []string `json:"permissions"`
}
