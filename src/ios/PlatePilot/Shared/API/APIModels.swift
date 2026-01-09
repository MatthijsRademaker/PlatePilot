import Foundation

struct APIErrorResponse: Decodable {
    let error: String?
}

struct PaginatedRecipesDTO: Decodable {
    let items: [RecipeDTO]?
    let pageIndex: Int?
    let pageSize: Int?
    let totalCount: Int?
    let totalPages: Int?
}

struct RecipeDTO: Decodable {
    let id: String?
    let name: String?
    let description: String?
    let prepTimeMinutes: Int?
    let cookTimeMinutes: Int?
    let totalTimeMinutes: Int?
    let servings: Int?
    let yieldQuantity: Double?
    let yieldUnit: String?
    let mainIngredient: IngredientRefDTO?
    let cuisine: CuisineDTO?
    let ingredientLines: [IngredientLineDTO]?
    let steps: [RecipeStepDTO]?
    let tags: [String]?
    let imageUrl: String?
    let nutrition: RecipeNutritionDTO?
}

struct CuisineDTO: Decodable {
    let id: String?
    let name: String?
}

struct IngredientRefDTO: Decodable {
    let id: String?
    let name: String?
}

struct IngredientLineDTO: Decodable {
    let ingredient: IngredientRefDTO?
    let quantityValue: Double?
    let quantityText: String?
    let unit: String?
    let isOptional: Bool?
    let note: String?
    let sortOrder: Int?
}

struct RecipeStepDTO: Decodable {
    let stepIndex: Int?
    let instruction: String?
    let durationSeconds: Int?
    let temperatureValue: Double?
    let temperatureUnit: String?
    let mediaUrl: String?
}

struct RecipeNutritionDTO: Codable {
    let caloriesTotal: Int?
    let caloriesPerServing: Int?
    let proteinG: Double?
    let carbsG: Double?
    let fatG: Double?
    let fiberG: Double?
    let sugarG: Double?
    let sodiumMg: Double?
}

struct MealPlanWeekDTO: Decodable {
    let startDate: String?
    let endDate: String?
    let days: [MealPlanDayDTO]?
}

struct MealPlanDayDTO: Decodable {
    let date: String?
    let meals: [MealPlanSlotDTO]?
}

struct MealPlanSlotDTO: Decodable {
    let id: String?
    let date: String?
    let mealType: String?
    let recipe: MealPlanRecipeDTO?
}

struct MealPlanRecipeDTO: Decodable {
    let id: String?
    let name: String?
    let description: String?
}

struct CreateRecipeRequestDTO: Encodable {
    let name: String
    let description: String?
    let prepTimeMinutes: Int
    let cookTimeMinutes: Int
    let servings: Int
    let yieldQuantity: Double?
    let yieldUnit: String?
    let mainIngredientId: String?
    let mainIngredientName: String?
    let cuisineId: String?
    let cuisineName: String?
    let ingredientLines: [IngredientLineInputDTO]
    let steps: [RecipeStepInputDTO]
    let tags: [String]
    let imageUrl: String?
    let nutrition: RecipeNutritionDTO?
}

struct IngredientLineInputDTO: Encodable {
    let ingredientId: String?
    let ingredientName: String?
    let quantityValue: Double?
    let quantityText: String?
    let unit: String?
    let isOptional: Bool
    let note: String?
    let sortOrder: Int
}

struct RecipeStepInputDTO: Encodable {
    let stepIndex: Int
    let instruction: String
    let durationSeconds: Int?
    let temperatureValue: Double?
    let temperatureUnit: String?
    let mediaUrl: String?
}

struct MealPlanWeekSaveRequestDTO: Encodable {
    let startDate: String
    let endDate: String
    let days: [MealPlanDayInputDTO]
}

struct MealPlanDayInputDTO: Encodable {
    let date: String
    let meals: [MealPlanSlotInputDTO]
}

struct MealPlanSlotInputDTO: Encodable {
    let mealType: String
    let recipeId: String?
}

struct MealPlanSuggestRequestDTO: Encodable {
    let dailyConstraints: [MealPlanDailyConstraintsDTO]
    let alreadySelectedRecipeIds: [String]
    let amount: Int
}

struct MealPlanDailyConstraintsDTO: Encodable {
    let ingredientConstraints: [MealPlanEntityConstraintDTO]
    let cuisineConstraints: [MealPlanEntityConstraintDTO]
}

struct MealPlanEntityConstraintDTO: Encodable {
    let entityId: String
}

struct MealPlanSuggestResponseDTO: Decodable {
    let recipeIds: [String]?
}

struct CuisinesResponseDTO: Decodable {
    let items: [CuisineDTO]?
}

struct CreateCuisineRequestDTO: Encodable {
    let name: String
}

struct RegisterRequestDTO: Encodable {
    let email: String
    let password: String
    let displayName: String
}

struct LoginRequestDTO: Encodable {
    let email: String
    let password: String
}

struct RefreshRequestDTO: Encodable {
    let refreshToken: String
}

struct TokenResponseDTO: Decodable {
    let accessToken: String
    let refreshToken: String
    let tokenType: String
    let expiresIn: Int64
}

struct EmptyResponse: Decodable {}
