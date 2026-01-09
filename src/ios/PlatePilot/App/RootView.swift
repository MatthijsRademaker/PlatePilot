import SwiftUI

@MainActor
struct RootView: View {
    @State private var authStore: AuthStore
    @State private var appState = AppState()
    @State private var tabRouter = TabRouter()
    @State private var recipeStore: RecipeStore
    @State private var mealPlanStore: MealPlanStore

    init() {
        let authStore = AuthStore()
        let apiClient = APIClient(tokenProvider: { await MainActor.run { authStore.accessToken } })
        _authStore = State(initialValue: authStore)
        _recipeStore = State(initialValue: RecipeStore(apiClient: apiClient))
        _mealPlanStore = State(initialValue: MealPlanStore(apiClient: apiClient))
    }

    var body: some View {
        Group {
            if authStore.isAuthenticated {
                AppView(
                    appState: appState,
                    tabRouter: tabRouter,
                    recipeStore: recipeStore,
                    mealPlanStore: mealPlanStore
                )
            } else {
                AuthFlowView()
            }
        }
        .environment(authStore)
        .task {
            await authStore.refreshIfNeeded()
        }
        .onChange(of: authStore.isAuthenticated) { _, isAuthenticated in
            if !isAuthenticated {
                recipeStore.reset()
            }
        }
        .preferredColorScheme(.light)
    }
}
