import SwiftUI

@MainActor
struct AppView: View {
    @Bindable var appState: AppState
    let tabRouter: TabRouter
    let recipeStore: RecipeStore
    let mealPlanStore: MealPlanStore

    var body: some View {
        TabView(selection: $appState.selectedTab) {
            ForEach(AppTab.allCases) { tab in
                NavigationStack(path: tabRouter.binding(for: tab)) {
                    tab.makeContentView()
                        .withAppRoutes()
                }
                .environment(tabRouter.router(for: tab))
                .tabItem { tab.label }
                .tag(tab)
            }
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
