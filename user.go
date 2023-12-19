package main

type User struct {
	Name    string `json:"name"`
	SavedAt string `json:"saved_at"`
}

func (u User) ToMetrics() map[string]any {
	return map[string]any{
		"user_name":     u.Name,
		"user_saved_at": u.SavedAt,
	}
}
