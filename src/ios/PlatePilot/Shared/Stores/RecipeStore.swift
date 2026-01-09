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
        cuisineName: String?
    ) async throws -> Recipe {
        let cleanedIngredients = ingredients
            .map { ingredient in
                RecipeIngredientInput(
                    name: ingredient.name.trimmingCharacters(in: .whitespacesAndNewlines),
                    quantity: ingredient.quantity.trimmingCharacters(in: .whitespacesAndNewlines),
                    unit: ingredient.unit.trimmingCharacters(in: .whitespacesAndNewlines)
                )
            }
            .filter { !$0.name.isEmpty }

        let ingredientLines = cleanedIngredients.enumerated().map { index, ingredient in
            IngredientLineInputDTO(
                ingredientId: nil,
                ingredientName: ingredient.name,
                quantityValue: nil,
                quantityText: ingredient.quantity.isEmpty ? nil : ingredient.quantity,
                unit: ingredient.unit.isEmpty ? nil : ingredient.unit,
                isOptional: false,
                note: nil,
                sortOrder: index + 1
            )
        }

        let steps = instructions
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
            .enumerated()
            .map { index, instruction in
                RecipeStepInputDTO(
                    stepIndex: index + 1,
                    instruction: instruction,
                    durationSeconds: nil,
                    temperatureValue: nil,
                    temperatureUnit: nil,
                    mediaUrl: nil
                )
            }

        let cleanedTags = tags
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }

        let payload = CreateRecipeRequestDTO(
            name: name.trimmingCharacters(in: .whitespacesAndNewlines),
            description: description.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty
                ? nil
                : description.trimmingCharacters(in: .whitespacesAndNewlines),
            prepTimeMinutes: prepMinutes,
            cookTimeMinutes: cookMinutes,
            servings: 1,
            yieldQuantity: nil,
            yieldUnit: nil,
            mainIngredientId: nil,
            mainIngredientName: ingredientLines.first?.ingredientName,
            cuisineId: nil,
            cuisineName: cuisineName?.trimmingCharacters(in: .whitespacesAndNewlines),
            ingredientLines: ingredientLines,
            steps: steps,
            tags: cleanedTags,
            imageUrl: nil,
            nutrition: nil
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

    func loadUnits() async throws -> [String] {
        return []
    }

    func createUnit(name: String) async throws -> String {
        let cleaned = name.trimmingCharacters(in: .whitespacesAndNewlines)
        return cleaned.isEmpty ? name : cleaned
    }

    func loadCuisines() async throws -> [String] {
        let cuisines = try await apiClient.fetchCuisines()
        return cuisines
            .compactMap { $0.name?.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
    }

    func createCuisine(name: String) async throws -> String {
        let cuisine = try await apiClient.createCuisine(name: name)
        return cuisine.name?.trimmingCharacters(in: .whitespacesAndNewlines) ?? name
    }
}

struct RecipeIngredientInput: Hashable {
    let name: String
    let quantity: String
    let unit: String
}
