import SwiftUI

enum AppTab: String, CaseIterable, Identifiable, Hashable {
    case home
    case recipes
    case mealPlan
    case search

    var id: String { rawValue }

    @ViewBuilder
    func makeContentView() -> some View {
        switch self {
        case .home:
            HomeView()
        case .recipes:
            RecipeListView()
        case .mealPlan:
            MealPlanView()
        case .search:
            SearchView()
        }
    }

    @ViewBuilder
    var label: some View {
        switch self {
        case .home:
            Label("Home", systemImage: "house")
        case .recipes:
            Label("Recipes", systemImage: "book")
        case .mealPlan:
            Label("Plan", systemImage: "calendar")
        case .search:
            Label("Search", systemImage: "magnifyingglass")
        }
    }
}
