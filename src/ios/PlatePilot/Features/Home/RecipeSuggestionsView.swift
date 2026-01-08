import SwiftUI

struct RecipeSuggestionsView: View {
    @Environment(RecipeStore.self) private var recipeStore
    @Environment(RouterPath.self) private var router
    @Environment(AppState.self) private var appState

    private var suggestedRecipes: [Recipe] {
        Array(recipeStore.recipes.prefix(6))
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack {
                Text("Recipe Suggestions")
                    .font(PlatePilotTheme.titleFont(size: 18))
                    .foregroundStyle(PlatePilotTheme.textPrimary)

                Spacer()

                Button("See all") {
                    appState.selectedTab = .recipes
                }
                .font(PlatePilotTheme.bodyFont(size: 13, weight: .semibold))
                .foregroundStyle(PlatePilotTheme.accent)
            }

            if recipeStore.isLoading {
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 12) {
                        ForEach(0..<4, id: \.self) { _ in
                            RoundedRectangle(cornerRadius: 16, style: .continuous)
                                .fill(PlatePilotTheme.tintWarm)
                                .frame(width: 140, height: 160)
                                .redacted(reason: .placeholder)
                        }
                    }
                    .padding(.horizontal, 4)
                }
            } else if suggestedRecipes.isEmpty {
                VStack(spacing: 12) {
                    Image(systemName: "fork.knife")
                        .font(.system(size: 28, weight: .semibold))
                        .foregroundStyle(.secondary)
                    Text("No recipes available")
                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                        .foregroundStyle(PlatePilotTheme.textSecondary)
                }
                .frame(maxWidth: .infinity)
                .padding(.vertical, 16)
            } else {
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 12) {
                        ForEach(suggestedRecipes) { recipe in
                            Button {
                                router.push(.recipeDetail(id: recipe.id))
                            } label: {
                                VStack(alignment: .leading, spacing: 8) {
                                    RemoteImageView(url: recipe.listImageURL, cornerRadius: 16)
                                        .frame(width: 140, height: 100)

                                    Text(recipe.name)
                                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .semibold))
                                        .foregroundStyle(PlatePilotTheme.textPrimary)
                                        .lineLimit(1)

                                    if let cuisine = recipe.cuisine {
                                        Text(cuisine)
                                            .font(PlatePilotTheme.bodyFont(size: 12, weight: .medium))
                                            .foregroundStyle(PlatePilotTheme.textSecondary)
                                    }
                                }
                                .padding(12)
                                .frame(width: 160, alignment: .leading)
                                .plateGlass(cornerRadius: 16, tint: .white.opacity(0.2), interactive: true)
                            }
                            .buttonStyle(.plain)
                        }
                    }
                    .padding(.horizontal, 4)
                }
            }
        }
        .padding(20)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
        .task {
            if recipeStore.recipes.isEmpty {
                await recipeStore.refresh()
            }
        }
    }
}
