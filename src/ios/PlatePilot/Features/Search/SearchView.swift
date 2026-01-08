import SwiftUI

struct SearchView: View {
    @Environment(RecipeStore.self) private var recipeStore
    @Environment(RouterPath.self) private var router

    @State private var searchQuery = ""

    private var results: [Recipe] {
        let trimmed = searchQuery.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else { return [] }
        return recipeStore.recipes.filter {
            $0.name.localizedCaseInsensitiveContains(trimmed) ||
            $0.description.localizedCaseInsensitiveContains(trimmed)
        }
    }

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                SearchHeaderView()

                VStack(spacing: 12) {
                    if recipeStore.isLoading {
                        ProgressView()
                            .padding(.top, 40)
                    } else if results.isEmpty && !searchQuery.isEmpty {
                        EmptySearchStateView(
                            icon: "magnifyingglass",
                            title: "No recipes found",
                            message: "Try a different search term."
                        )
                    } else if results.isEmpty {
                        EmptySearchStateView(
                            icon: "fork.knife",
                            title: "Start searching",
                            message: "Find recipes by name or description."
                        )
                    } else {
                        LazyVStack(spacing: 12) {
                            ForEach(results) { recipe in
                                Button {
                                    router.push(.recipeDetail(id: recipe.id))
                                } label: {
                                    RecipeCardView(recipe: recipe)
                                }
                                .buttonStyle(.plain)
                            }
                        }
                    }
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
        .searchable(text: $searchQuery, placement: .navigationBarDrawer(displayMode: .always))
        .task {
            if recipeStore.recipes.isEmpty {
                await recipeStore.refresh()
            }
        }
    }
}

private struct SearchHeaderView: View {
    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(PlatePilotTheme.headerGradient)
                .frame(maxWidth: .infinity)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                .padding(.bottom, 4)

            HStack(spacing: 12) {
                Image(systemName: "magnifyingglass")
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .plateGlass(cornerRadius: 12, tint: .white.opacity(0.25))

                Text("Search Recipes")
                    .font(PlatePilotTheme.titleFont(size: 24))
                    .foregroundStyle(.white)

                Spacer()
            }
            .padding(.horizontal, 32)
            .padding(.vertical, 22)
        }
        .padding(.top, 12)
    }
}

private struct EmptySearchStateView: View {
    let icon: String
    let title: String
    let message: String

    var body: some View {
        VStack(spacing: 12) {
            Image(systemName: icon)
                .font(.system(size: 34, weight: .semibold))
                .foregroundStyle(PlatePilotTheme.accent)
                .frame(width: 80, height: 80)
                .background(PlatePilotTheme.tintWarm, in: RoundedRectangle(cornerRadius: 24, style: .continuous))

            Text(title)
                .font(PlatePilotTheme.titleFont(size: 20))
                .foregroundStyle(PlatePilotTheme.textPrimary)

            Text(message)
                .font(PlatePilotTheme.bodyFont(size: 14))
                .foregroundStyle(PlatePilotTheme.textSecondary)
        }
        .frame(maxWidth: .infinity)
        .padding(24)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
    }
}
