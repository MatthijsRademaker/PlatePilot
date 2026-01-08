import SwiftUI

struct RecipeDetailView: View {
    let recipeID: UUID

    @Environment(RecipeStore.self) private var recipeStore
    @Environment(\.dismiss) private var dismiss
    @State private var isLoading = false
    @State private var errorMessage: String?

    private var recipe: Recipe? {
        recipeStore.recipe(id: recipeID)
    }

    var body: some View {
        ScrollView {
            if isLoading {
                VStack(spacing: 12) {
                    ProgressView()
                    Text("Loading recipe...")
                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                        .foregroundStyle(PlatePilotTheme.textSecondary)
                }
                .padding(24)
            } else if let errorMessage {
                VStack(spacing: 16) {
                    Image(systemName: "exclamationmark.triangle")
                        .font(.system(size: 32, weight: .semibold))
                        .foregroundStyle(.orange)

                    Text(errorMessage)
                        .font(PlatePilotTheme.bodyFont(size: 16, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.textSecondary)

                    Button("Go Back") {
                        dismiss()
                    }
                    .plateGlassButton(prominent: true)
                }
                .padding(24)
            } else if let recipe {
                VStack(spacing: 0) {
                    heroView(for: recipe)

                    VStack(alignment: .leading, spacing: 16) {
                        if !recipe.description.isEmpty {
                            Text(recipe.description)
                                .font(PlatePilotTheme.bodyFont(size: 15))
                                .foregroundStyle(PlatePilotTheme.textSecondary)
                        }

                        SectionCard(title: "Ingredients", icon: "checklist") {
                            VStack(alignment: .leading, spacing: 10) {
                                ForEach(recipe.ingredients, id: \.self) { ingredient in
                                    HStack(spacing: 10) {
                                        Image(systemName: "checkmark.circle.fill")
                                            .foregroundStyle(PlatePilotTheme.accent)
                                        Text(ingredient)
                                            .font(PlatePilotTheme.bodyFont(size: 15))
                                            .foregroundStyle(PlatePilotTheme.textPrimary)
                                    }
                                }
                            }
                        }

                        SectionCard(title: "Directions", icon: "list.number") {
                            VStack(alignment: .leading, spacing: 12) {
                                ForEach(Array(recipe.directions.enumerated()), id: \.offset) { index, step in
                                    HStack(alignment: .top, spacing: 12) {
                                        Text("\(index + 1)")
                                            .font(PlatePilotTheme.bodyFont(size: 13, weight: .bold))
                                            .foregroundStyle(.white)
                                            .frame(width: 28, height: 28)
                                            .background(PlatePilotTheme.headerGradient, in: RoundedRectangle(cornerRadius: 8, style: .continuous))

                                        Text(step)
                                            .font(PlatePilotTheme.bodyFont(size: 15))
                                            .foregroundStyle(PlatePilotTheme.textPrimary)
                                    }
                                }
                            }
                        }
                    }
                    .padding(20)
                }
            } else {
                VStack(spacing: 16) {
                    Image(systemName: "exclamationmark.triangle")
                        .font(.system(size: 32, weight: .semibold))
                        .foregroundStyle(.orange)

                    Text("Recipe not found")
                        .font(PlatePilotTheme.bodyFont(size: 16, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.textSecondary)

                    Button("Go Back") {
                        dismiss()
                    }
                    .plateGlassButton(prominent: true)
                }
                .padding(24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
        .ignoresSafeArea(edges: .top)
        .navigationBarBackButtonHidden(true)
        .toolbar {
            ToolbarItem(placement: .topBarLeading) {
                Button {
                    dismiss()
                } label: {
                    Image(systemName: "chevron.left")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundStyle(.white)
                        .frame(width: 34, height: 34)
                        .plateGlass(cornerRadius: 12, tint: .black.opacity(0.35), interactive: true)
                }
            }
        }
        .task(id: recipeID) {
            guard recipeStore.recipe(id: recipeID) == nil else { return }
            isLoading = true
            errorMessage = nil
            do {
                _ = try await recipeStore.fetchRecipe(id: recipeID)
            } catch {
                errorMessage = (error as? APIError)?.errorDescription ?? "Unable to fetch recipe."
            }
            isLoading = false
        }
    }

    @ViewBuilder
    private func heroView(for recipe: Recipe) -> some View {
        ZStack(alignment: .bottomLeading) {
            RemoteImageView(url: recipe.detailImageURL, cornerRadius: 0)
                .frame(height: 280)

            LinearGradient(
                colors: [Color.black.opacity(0.05), Color.black.opacity(0.6)],
                startPoint: .top,
                endPoint: .bottom
            )

            VStack(alignment: .leading, spacing: 12) {
                if let cuisine = recipe.cuisine {
                    Text(cuisine)
                        .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                        .foregroundStyle(.white)
                        .padding(.horizontal, 12)
                        .padding(.vertical, 6)
                        .plateGlass(cornerRadius: 16, tint: PlatePilotTheme.accent.opacity(0.4))
                }

                Text(recipe.name)
                    .font(PlatePilotTheme.titleFont(size: 28))
                    .foregroundStyle(.white)

                PlateGlassGroup(spacing: 12) {
                    HStack(spacing: 10) {
                        if let prepTime = recipe.prepTime {
                            TimeBadge(label: "Prep", value: prepTime)
                        }
                        if let cookTime = recipe.cookTime {
                            TimeBadge(label: "Cook", value: cookTime)
                        }
                    }
                }
            }
            .padding(20)
        }
    }
}

private struct TimeBadge: View {
    let label: String
    let value: String

    var body: some View {
        HStack(spacing: 6) {
            Image(systemName: "clock")
                .font(.system(size: 12, weight: .semibold))
            Text("\(label): \(value)")
                .font(PlatePilotTheme.bodyFont(size: 12, weight: .medium))
        }
        .foregroundStyle(.white)
        .padding(.horizontal, 10)
        .padding(.vertical, 6)
        .plateGlass(cornerRadius: 10, tint: .white.opacity(0.2))
    }
}

private struct SectionCard<Content: View>: View {
    let title: String
    let icon: String
    let content: () -> Content

    init(title: String, icon: String, @ViewBuilder content: @escaping () -> Content) {
        self.title = title
        self.icon = icon
        self.content = content
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack(spacing: 12) {
                Image(systemName: icon)
                    .font(.system(size: 16, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 36, height: 36)
                    .background(PlatePilotTheme.headerGradient, in: RoundedRectangle(cornerRadius: 10, style: .continuous))

                Text(title)
                    .font(PlatePilotTheme.titleFont(size: 18))
                    .foregroundStyle(PlatePilotTheme.textPrimary)
            }

            content()
        }
        .padding(20)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
    }
}
