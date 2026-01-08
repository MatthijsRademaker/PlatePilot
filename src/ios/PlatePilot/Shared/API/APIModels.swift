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
}

struct CreateRecipeRequestDTO: Encodable {
    let name: String
    let description: String
    let prepTime: String
    let cookTime: String
    let mainIngredientName: String?
    let cuisineName: String?
    let ingredientNames: [String]
    let directions: [String]
    let tags: [String]
    let guidedMode: Bool
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
