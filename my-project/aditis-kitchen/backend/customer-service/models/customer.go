package models

type Customer struct {
	ID        string `json:"id"`                // auto: mobile_name
	Name      string `json:"name"`              // required
	Mobile    string `json:"mobile"`            // required, must be 10 digits
	Email     string `json:"email,omitempty"`   // optional
	Password  string `json:"password"`          // required (hashed)
    Address   string `json:"address,omitempty"`
	IsVerified bool	  `json:"is_verified"`
}
