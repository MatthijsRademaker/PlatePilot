import Observation
import SwiftUI

@MainActor
@Observable
final class RouterPath {
    var path: [Route] = []

    func push(_ route: Route) {
        path.append(route)
    }

    func reset() {
        path.removeAll()
    }
}

enum Route: Hashable {
    case recipeDetail(id: UUID)
    case hubDestination(HubDestination)
}

enum HubDestination: Hashable {
    case recipesSearch
    case recipesCreate
    case mealPlanSuggest
    case mealPlanScan
    case calorieAddExercise
    case calorieGoals
    case insightsRecipes
    case insightsMealPlan
    case insightsCalories

    var title: String {
        switch self {
        case .recipesSearch:
            return "Recipe Search"
        case .recipesCreate:
            return "Create Recipe"
        case .mealPlanSuggest:
            return "Suggest Meals"
        case .mealPlanScan:
            return "Scan Ingredients"
        case .calorieAddExercise:
            return "Add Exercise"
        case .calorieGoals:
            return "Calorie Goals"
        case .insightsRecipes:
            return "Recipe Insights"
        case .insightsMealPlan:
            return "Meal Plan Insights"
        case .insightsCalories:
            return "Calorie Tracker Insights"
        }
    }

    var subtitle: String {
        switch self {
        case .recipesSearch:
            return "Find recipes, filters, and collections live here."
        case .recipesCreate:
            return "Draft new recipes with steps, nutrition, and tags."
        case .mealPlanSuggest:
            return "Personalized meal suggestions will appear here."
        case .mealPlanScan:
            return "Scan items to build a plan from your pantry."
        case .calorieAddExercise:
            return "Log workouts and sync activity insights."
        case .calorieGoals:
            return "Set targets and tune your daily burn."
        case .insightsRecipes:
            return "Recipe trends and favorites live here."
        case .insightsMealPlan:
            return "Weekly planning insights and patterns."
        case .insightsCalories:
            return "Calorie trendlines and streaks."
        }
    }

    var icon: String {
        switch self {
        case .recipesSearch:
            return "magnifyingglass"
        case .recipesCreate:
            return "plus"
        case .mealPlanSuggest:
            return "sparkles"
        case .mealPlanScan:
            return "camera.viewfinder"
        case .calorieAddExercise:
            return "figure.run"
        case .calorieGoals:
            return "target"
        case .insightsRecipes:
            return "book.fill"
        case .insightsMealPlan:
            return "calendar"
        case .insightsCalories:
            return "flame.fill"
        }
    }

    var accent: Color {
        switch self {
        case .recipesSearch, .recipesCreate:
            return Color(red: 1.0, green: 0.63, blue: 0.34)
        case .mealPlanSuggest, .mealPlanScan:
            return Color(red: 0.33, green: 0.82, blue: 0.52)
        case .calorieAddExercise, .calorieGoals:
            return Color(red: 1.0, green: 0.52, blue: 0.28)
        case .insightsRecipes, .insightsMealPlan, .insightsCalories:
            return Color(red: 0.4, green: 0.58, blue: 0.98)
        }
    }
}

@MainActor
@Observable
final class TabRouter {
    private var routers: [AppTab: RouterPath] = [:]

    func router(for tab: AppTab) -> RouterPath {
        if let router = routers[tab] { return router }
        let router = RouterPath()
        routers[tab] = router
        return router
    }

    func binding(for tab: AppTab) -> Binding<[Route]> {
        let router = router(for: tab)
        return Binding(get: { router.path }, set: { router.path = $0 })
    }
}
