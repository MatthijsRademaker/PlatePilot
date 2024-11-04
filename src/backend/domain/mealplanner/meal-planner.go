package mealplanner

type MealPlanner interface {
	Suggest(mealPlan *MealPlan, amountToSuggest int, constraints ...MealPlanConstraints) (*MealPlan, error)
	GenerateMealPlan(daysToPlanFor int, constraints ...MealPlanConstraints) (*MealPlan, error)
}

type MealPlanConstraints struct {
}
