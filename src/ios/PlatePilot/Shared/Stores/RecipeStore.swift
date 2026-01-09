import Foundation
import Observation

@MainActor
@Observable
final class RecipeStore {
    private let apiClient: APIClient

    private(set) var recipes: [Recipe]
    var isLoading = false
    var errorMessage: String?

    init(apiClient: APIClient = APIClient(), recipes: [Recipe] = []) {
        self.apiClient = apiClient
        self.recipes = recipes
    }

    func refresh() async {
        await loadRecipes(pageIndex: 1)
    }

    func loadRecipes(pageIndex: Int, pageSize: Int = 20) async {
        isLoading = true
        errorMessage = nil
        do {
            let response = try await apiClient.fetchRecipes(pageIndex: pageIndex, pageSize: pageSize)
            let items = response.items ?? []
            recipes = items.map(Recipe.init)
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to fetch recipes."
        }
        isLoading = false
    }

    func recipe(id: UUID) -> Recipe? {
        recipes.first { $0.id == id }
    }

    func fetchRecipe(id: UUID) async throws -> Recipe {
        if let cached = recipe(id: id) {
            return cached
        }
        isLoading = true
        errorMessage = nil
        defer { isLoading = false }

        do {
            let dto = try await apiClient.fetchRecipe(id: id)
            let recipe = Recipe(dto: dto)
            recipes.append(recipe)
            return recipe
        } catch {
            let message = (error as? APIError)?.errorDescription ?? "Unable to fetch recipe."
            errorMessage = message
            throw error
        }
    }

    func createRecipe(
        name: String,
        description: String,
        prepMinutes: Int,
        cookMinutes: Int,
        ingredients: [RecipeIngredientInput],
        instructions: [String],
        tags: [String],
        guidedMode: Bool
    ) async throws -> Recipe {
        let ingredientNames = ingredients.map { $0.name }
        let payload = CreateRecipeRequestDTO(
            name: name,
            description: description,
            prepTime: formattedTime(prepMinutes),
            cookTime: formattedTime(cookMinutes),
            mainIngredientName: ingredientNames.first,
            cuisineName: nil,
            ingredientNames: ingredientNames,
            ingredients: ingredients.map {
                CreateRecipeIngredientDTO(
                    id: nil,
                    name: $0.name,
                    quantity: $0.quantity,
                    unit: $0.unit
                )
            },
            directions: instructions,
            tags: tags,
            guidedMode: guidedMode
        )

        do {
            let dto = try await apiClient.createRecipe(payload: payload)
            let recipe = Recipe(dto: dto)
            recipes.insert(recipe, at: 0)
            return recipe
        } catch {
            throw error
        }
    }

    func reset() {
        recipes = []
        errorMessage = nil
        isLoading = false
    }

    private func formattedTime(_ minutes: Int) -> String {
        guard minutes > 0 else { return "" }
        return "\(minutes) min"
    }

    func loadUnits() async throws -> [String] {
        let units = try await apiClient.fetchUnits()
        return units
            .compactMap { $0.name?.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
    }

    func createUnit(name: String) async throws -> String {
        let unit = try await apiClient.createUnit(name: name)
        return unit.name?.trimmingCharacters(in: .whitespacesAndNewlines) ?? name
    }
}

struct RecipeIngredientInput: Hashable {
    let name: String
    let quantity: String
    let unit: String
}
