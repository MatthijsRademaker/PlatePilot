import SwiftUI

struct TodayPlanCard: View {
    @Environment(MealPlanStore.self) private var mealPlanStore
    @Environment(RouterPath.self) private var router
    @Environment(AppState.self) private var appState

    private var featuredMeal: MealSlot? {
        mealPlanStore.featuredMeal()
    }

    var body: some View {
        PlateGlassGroup(spacing: 24) {
            VStack(alignment: .leading, spacing: 12) {
                HStack(spacing: 10) {
                    Image(systemName: "fork.knife.circle.fill")
                        .font(.system(size: 18, weight: .semibold))
                        .foregroundStyle(.white)
                        .frame(width: 28, height: 28)
                        .plateGlass(cornerRadius: 10, tint: .white.opacity(0.25))

                    Text("Your Meal Plan Today")
                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .semibold))
                        .foregroundStyle(.white)
                }

                if let featuredMeal, let recipe = featuredMeal.recipe {
                    VStack(alignment: .leading, spacing: 12) {
                        RemoteImageView(url: recipe.detailImageURL, cornerRadius: 16)
                            .frame(height: 140)
                            .overlay(alignment: .topLeading) {
                                Text(featuredMeal.mealType.label)
                                    .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                                    .foregroundStyle(.white)
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 6)
                                    .plateGlass(cornerRadius: 16, tint: .black.opacity(0.35))
                                    .padding(12)
                            }

                        Text(recipe.name)
                            .font(PlatePilotTheme.titleFont(size: 20))
                            .foregroundStyle(PlatePilotTheme.textPrimary)

                        HStack(spacing: 12) {
                            Button("View Recipe") {
                                router.push(.recipeDetail(id: recipe.id))
                            }
                            .plateGlassButton(prominent: true)

                            Button {
                                appState.selectedTab = .mealPlan
                            } label: {
                                Image(systemName: "calendar")
                                    .font(.system(size: 16, weight: .semibold))
                                    .frame(width: 44, height: 44)
                            }
                            .plateGlassButton()
                        }
                    }
                    .padding(16)
                    .background(PlatePilotTheme.surface, in: RoundedRectangle(cornerRadius: 18, style: .continuous))
                } else {
                    VStack(spacing: 12) {
                        Image(systemName: "calendar.badge.exclamationmark")
                            .font(.system(size: 28, weight: .semibold))
                            .foregroundStyle(.secondary)
                            .frame(width: 56, height: 56)
                            .background(PlatePilotTheme.tintWarm, in: RoundedRectangle(cornerRadius: 16, style: .continuous))

                        Text("No meals planned for today")
                            .font(PlatePilotTheme.bodyFont(size: 15, weight: .medium))
                            .foregroundStyle(PlatePilotTheme.textSecondary)

                        Button("Plan Your Meals") {
                            appState.selectedTab = .mealPlan
                        }
                        .plateGlassButton(prominent: true)
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 16)
                    .padding(.horizontal, 12)
                    .background(PlatePilotTheme.surface, in: RoundedRectangle(cornerRadius: 18, style: .continuous))
                }
            }
            .padding(16)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(PlatePilotTheme.headerGradient, in: RoundedRectangle(cornerRadius: 24, style: .continuous))
            .shadow(color: PlatePilotTheme.accent.opacity(0.25), radius: 16, x: 0, y: 8)
        }
    }
}
