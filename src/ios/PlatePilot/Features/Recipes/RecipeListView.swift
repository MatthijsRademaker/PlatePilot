import SwiftUI

struct RecipeListView: View {
    @Environment(RecipeStore.self) private var recipeStore
    @Environment(RouterPath.self) private var router

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                RecipeListHeaderView {
                    Task { await recipeStore.refresh() }
                }

                VStack(spacing: 12) {
                    if recipeStore.isLoading {
                        ForEach(0..<4, id: \.self) { _ in
                            RoundedRectangle(cornerRadius: PlatePilotMetrics.cardRadius, style: .continuous)
                                .fill(PlatePilotTheme.tintWarm)
                                .frame(height: 110)
                                .redacted(reason: .placeholder)
                        }
                    } else if let errorMessage = recipeStore.errorMessage {
                        ErrorStateView(message: errorMessage) {
                            Task { await recipeStore.refresh() }
                        }
                    } else {
                        ForEach(recipeStore.recipes) { recipe in
                            Button {
                                router.push(.recipeDetail(id: recipe.id))
                            } label: {
                                RecipeCardView(recipe: recipe)
                            }
                            .buttonStyle(.plain)
                        }
                    }
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
        .task {
            if recipeStore.recipes.isEmpty {
                await recipeStore.refresh()
            }
        }
    }
}

private struct RecipeListHeaderView: View {
    let onRefresh: () -> Void

    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(PlatePilotTheme.headerGradient)
                .frame(maxWidth: .infinity)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                .padding(.bottom, 4)

            HStack {
                HStack(spacing: 12) {
                    Image(systemName: "book.fill")
                        .font(.system(size: 20, weight: .semibold))
                        .foregroundStyle(.white)
                        .frame(width: 44, height: 44)
                        .plateGlass(cornerRadius: 12, tint: .white.opacity(0.25))

                    Text("Recipes")
                        .font(PlatePilotTheme.titleFont(size: 24))
                        .foregroundStyle(.white)
                }

                Spacer()

                Button(action: onRefresh) {
                    Image(systemName: "arrow.clockwise")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundStyle(.white)
                        .frame(width: 40, height: 40)
                        .plateGlass(cornerRadius: 12, tint: .white.opacity(0.2), interactive: true)
                }
                .buttonStyle(.plain)
            }
            .padding(.horizontal, 32)
            .padding(.vertical, 22)
        }
        .padding(.top, 12)
    }
}

private struct ErrorStateView: View {
    let message: String
    let onRetry: () -> Void

    var body: some View {
        VStack(spacing: 12) {
            Image(systemName: "exclamationmark.triangle")
                .font(.system(size: 28, weight: .semibold))
                .foregroundStyle(.orange)

            Text(message)
                .font(PlatePilotTheme.bodyFont(size: 15, weight: .medium))
                .foregroundStyle(PlatePilotTheme.textSecondary)

            Button("Try Again", action: onRetry)
                .plateGlassButton(prominent: true)
        }
        .frame(maxWidth: .infinity)
        .padding(20)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
    }
}
