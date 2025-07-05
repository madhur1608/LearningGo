package models

type MenuItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Days        []string `json:"days"`
	Type        string   `json:"type"`
}

type MealPlan struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Name       string `json:"name"`
	PlanPrice      int    `json:"plan_price"`
	TotalMeals int    `json:"total_meals"` // total meal for the plan selected
	Type       string `json:"type"`     // "plan"
}
