update-bff-api-dependencies:
	cd MobileBFF/Infrastructure && \
	kiota generate -l CSharp -c RecipeApi -n Infrastructure.RecipeApi.Client -d ../../RecipeApi/Application/RecipeApplication.json -o ./Recipes/ApiClient && \
	kiota generate -l CSharp -c MealPlannerApi -n Infrastructure.MealPlannerApi.Client -d ../../MealPlannerApi/Application/MealPlannerApplication.json -o ./MealPlanning/ApiClient
