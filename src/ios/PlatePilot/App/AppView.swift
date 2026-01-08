import SwiftUI

@MainActor
struct AppView: View {
    @Bindable var appState: AppState
    let tabRouter: TabRouter
    let recipeStore: RecipeStore
    let mealPlanStore: MealPlanStore

    var body: some View {
        NavigationStack(path: tabRouter.binding(for: appState.selectedTab)) {
            appState.selectedTab.makeContentView()
                .withAppRoutes()
        }
        .id(appState.selectedTab)
        .environment(tabRouter.router(for: appState.selectedTab))
        .safeAreaInset(edge: .bottom) {
            GlassHubNavigationView()
        }
        .tint(PlatePilotTheme.accent)
        .environment(appState)
        .environment(recipeStore)
        .environment(mealPlanStore)
    }
}

extension View {
    func withAppRoutes() -> some View {
        navigationDestination(for: Route.self) { route in
            switch route {
            case .recipeDetail(let id):
                RecipeDetailView(recipeID: id)
            }
        }
    }
}
