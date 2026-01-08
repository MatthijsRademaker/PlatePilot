import SwiftUI

enum AppTab: String, CaseIterable, Identifiable, Hashable {
    case home
    case recipes
    case calorieTracker
    case mealPlan
    case insights

    var id: String { rawValue }

    @ViewBuilder
    func makeContentView() -> some View {
        switch self {
        case .home:
            HomeView()
        case .recipes:
            RecipeListView()
        case .calorieTracker:
            CalorieTrackerView()
        case .mealPlan:
            MealPlanView()
        case .insights:
            InsightsView()
        }
    }

    @ViewBuilder
    var label: some View {
        switch self {
        case .home:
            Label("Home", systemImage: "house")
        case .recipes:
            Label("Recipes", systemImage: "book")
        case .calorieTracker:
            Label("Calories", systemImage: "flame")
        case .mealPlan:
            Label("Plan", systemImage: "calendar")
        case .insights:
            Label("Insights", systemImage: "chart.bar.xaxis")
        }
    }
}
