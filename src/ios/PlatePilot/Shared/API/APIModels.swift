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
    let cuisine: CuisineDTO?
    let prepTime: String?
    let cookTime: String?
    let ingredients: [IngredientDTO]?
    let directions: [String]?
}

struct CuisineDTO: Decodable {
    let id: String?
    let name: String?
}

struct IngredientDTO: Decodable {
    let id: String?
    let name: String?
    let quantity: String?
    let unit: String?
}

struct CreateRecipeRequestDTO: Encodable {
    let name: String
    let description: String
    let prepTime: String
    let cookTime: String
    let mainIngredientName: String?
    let cuisineName: String?
    let ingredientNames: [String]
    let ingredients: [CreateRecipeIngredientDTO]
    let directions: [String]
    let tags: [String]
    let guidedMode: Bool
}

struct CreateRecipeIngredientDTO: Encodable {
    let id: String?
    let name: String
    let quantity: String
    let unit: String
}

struct UnitDTO: Decodable {
    let id: String?
    let name: String?
}

struct UnitsResponseDTO: Decodable {
    let items: [UnitDTO]?
}

struct CreateUnitRequestDTO: Encodable {
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
