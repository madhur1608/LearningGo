package models

type Subscription struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	PlanID     string `json:"plan_id"`
	PlanName   string `json:"plan_name"`
	StartDate  string `json:"start_date"` // full timestamp in IST
	EndDate    string `json:"end_date"`   // YYYY-MM-DD
	TotalMeals int    `json:"total_meals"`
	PlanPrice      int    `json:"plan_price"`
	Status     string `json:"status"`     // active, cancelled, etc.
	MealType   string `json:"meal_type"`  // lunch or dinner
	Type       string `json:"type"`       // "subscription"
	CreatedAt  string `json:"created_at"` // timestamp
	UpdatedAt  string `json:"updated_at,omitempty"`
}