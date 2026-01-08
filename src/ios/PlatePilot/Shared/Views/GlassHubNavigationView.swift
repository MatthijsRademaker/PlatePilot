import SwiftUI

struct GlassHubNavigationView: View {
    @Environment(AppState.self) private var appState
    @Environment(RouterPath.self) private var router

    private var currentSection: HubSection {
        HubSection(tab: appState.selectedTab) ?? .home
    }

    private var primarySections: [HubSection] {
        [.recipes, .calorieTracker, .home, .mealPlan, .insights]
    }

    var body: some View {
        PlateGlassGroup(spacing: 12) {
            VStack(spacing: -HubMetrics.sideActionOverlap) {
                GlassSideActionBar(
                    leftActions: currentSection.leftActions,
                    rightActions: currentSection.rightActions,
                    accent: currentSection.accent,
                    onAction: handleAction
                )

                GlassPrimaryNavBar(
                    sections: primarySections,
                    activeSection: currentSection,
                    accent: currentSection.accent,
                    onSelect: selectPrimarySection
                )
            }
            .frame(maxWidth: .infinity)
        }
        .padding(.horizontal, 12)
        .padding(.top, 6)
        .padding(.bottom, 4)
    }

    private func selectPrimarySection(_ section: HubSection) {
        if let tab = section.tab {
            appState.selectedTab = tab
        }
    }

    private func handleAction(_ action: HubAction) {
        router.push(.hubDestination(action.destination))
    }
}

private struct GlassSideActionBar: View {
    let leftActions: [HubAction]
    let rightActions: [HubAction]
    let accent: Color
    let onAction: (HubAction) -> Void

    var body: some View {
        HStack(spacing: 0) {
            actionPanel(actions: leftActions, alignment: .leading)
            actionPanel(actions: rightActions, alignment: .trailing)
        }
        .padding(.horizontal, 22)
        .frame(height: HubMetrics.sideActionSize)
    }

    @ViewBuilder
    private func actionPanel(actions: [HubAction], alignment: Alignment) -> some View {
        let items = Array(actions.prefix(2))
        HStack(spacing: HubMetrics.sideActionSpacing) {
            ForEach(items) { action in
                SideActionButton(action: action, accent: accent) {
                    onAction(action)
                }
            }
        }
        .frame(maxWidth: .infinity, alignment: alignment)
    }
}

private struct SideActionButton: View {
    let action: HubAction
    let accent: Color
    let onTap: () -> Void

    var body: some View {
        Button {
            onTap()
        } label: {
            ZStack {
                Circle()
                    .fill(Color.white.opacity(0.18))

                Image(systemName: action.icon)
                    .font(.system(size: HubMetrics.sideActionIconSize, weight: .semibold))
                    .foregroundStyle(accent)
            }
            .frame(width: HubMetrics.sideActionSize, height: HubMetrics.sideActionSize)
            .contentShape(Circle())
            .plateGlass(
                cornerRadius: HubMetrics.sideActionSize / 2,
                tint: accent.opacity(0.22),
                interactive: true
            )
            .overlay(
                Circle()
                    .stroke(Color.white.opacity(0.2), lineWidth: 0.6)
            )
            .shadow(color: accent.opacity(0.12), radius: 6, x: 0, y: 4)
        }
        .buttonStyle(.plain)
        .hoverEffect(.lift)
        .accessibilityLabel(action.title)
    }
}

private struct GlassPrimaryNavBar: View {
    let sections: [HubSection]
    let activeSection: HubSection
    let accent: Color
    let onSelect: (HubSection) -> Void

    var body: some View {
        HStack(spacing: 0) {
            ForEach(sections) { section in
                PrimaryNavButton(
                    icon: section.icon,
                    accent: section.accent,
                    isActive: section == activeSection,
                    accessibilityLabel: section.title
                ) {
                    onSelect(section)
                }
                .frame(maxWidth: .infinity)
            }
        }
        .frame(maxWidth: .infinity)
        .padding(.horizontal, 12)
        .frame(height: HubMetrics.barHeight)
        .background(
            RoundedRectangle(cornerRadius: HubMetrics.barCornerRadius, style: .continuous)
                .fill(Color.white.opacity(0.12))
                .overlay(
                    RoundedRectangle(cornerRadius: HubMetrics.barCornerRadius, style: .continuous)
                        .fill(
                            RadialGradient(
                                colors: [accent.opacity(0.35), .clear],
                                center: .center,
                                startRadius: 0,
                                endRadius: 220
                            )
                        )
                        .blendMode(.screen)
                )
        )
        .plateGlass(cornerRadius: HubMetrics.barCornerRadius, tint: .white.opacity(0.2))
        .overlay(
            RoundedRectangle(cornerRadius: HubMetrics.barCornerRadius, style: .continuous)
                .stroke(Color.white.opacity(0.22), lineWidth: 1)
        )
        .shadow(color: .black.opacity(0.08), radius: 16, x: 0, y: 8)
    }
}

private struct PrimaryNavButton: View {
    let icon: String
    let accent: Color
    let isActive: Bool
    let accessibilityLabel: String
    let onTap: () -> Void

    var body: some View {
        Button {
            onTap()
        } label: {
            ZStack {
                Circle()
                    .fill(.white.opacity(isActive ? 0.18 : 0.12))

                Image(systemName: icon)
                    .font(.system(size: HubMetrics.primaryIconSize, weight: .semibold))
                    .foregroundStyle(isActive ? accent : .white.opacity(0.85))
            }
            .frame(width: HubMetrics.primaryButtonSize, height: HubMetrics.primaryButtonSize)
            .contentShape(Circle())
            .plateGlass(
                cornerRadius: HubMetrics.primaryButtonSize / 2,
                tint: accent.opacity(isActive ? 0.35 : 0.18),
                interactive: true
            )
            .overlay(
                Circle()
                    .stroke(
                        isActive ? accent.opacity(0.75) : Color.white.opacity(0.2),
                        lineWidth: isActive ? 1.2 : 0.6
                    )
            )
            .shadow(
                color: accent.opacity(isActive ? 0.3 : 0.12),
                radius: isActive ? 8 : 6,
                x: 0,
                y: 4
            )
        }
        .buttonStyle(.plain)
        .hoverEffect(.lift)
        .accessibilityLabel(accessibilityLabel)
    }
}

private struct HubAction: Identifiable, Hashable {
    let destination: HubDestination
    let title: String
    let icon: String

    var id: HubDestination { destination }
}

private enum HubSection: String, CaseIterable, Identifiable {
    case home
    case recipes
    case calorieTracker
    case mealPlan
    case insights

    var id: String { rawValue }

    init?(tab: AppTab) {
        switch tab {
        case .home:
            self = .home
        case .recipes:
            self = .recipes
        case .calorieTracker:
            self = .calorieTracker
        case .mealPlan:
            self = .mealPlan
        case .insights:
            self = .insights
        }
    }

    var title: String {
        switch self {
        case .home: return "Home"
        case .recipes: return "Recipes"
        case .calorieTracker: return "Calorie Tracker"
        case .mealPlan: return "Meal Plan"
        case .insights: return "Insights"
        }
    }

    var icon: String {
        switch self {
        case .home: return "house"
        case .recipes: return "book.fill"
        case .calorieTracker: return "flame.fill"
        case .mealPlan: return "calendar"
        case .insights: return "chart.bar.xaxis"
        }
    }

    var accent: Color {
        switch self {
        case .home: return PlatePilotTheme.accent
        case .recipes: return Color(red: 1.0, green: 0.63, blue: 0.34)
        case .calorieTracker: return Color(red: 1.0, green: 0.52, blue: 0.28)
        case .mealPlan: return Color(red: 0.33, green: 0.82, blue: 0.52)
        case .insights: return Color(red: 0.4, green: 0.58, blue: 0.98)
        }
    }

    var tab: AppTab? {
        switch self {
        case .home: return .home
        case .recipes: return .recipes
        case .calorieTracker: return .calorieTracker
        case .mealPlan: return .mealPlan
        case .insights: return .insights
        }
    }

    var leftActions: [HubAction] {
        switch self {
        case .home:
            return []
        case .recipes:
            return [
                HubAction(destination: .recipesSearch, title: "Search", icon: "magnifyingglass")
            ]
        case .calorieTracker:
            return [
                HubAction(destination: .calorieAddExercise, title: "Add Exercise", icon: "figure.run")
            ]
        case .mealPlan:
            return [
                HubAction(destination: .mealPlanSuggest, title: "Suggest", icon: "sparkles")
            ]
        case .insights:
            return [
                HubAction(destination: .insightsRecipes, title: "Recipes", icon: "book.fill"),
                HubAction(destination: .insightsMealPlan, title: "Meal Plan", icon: "calendar")
            ]
        }
    }

    var rightActions: [HubAction] {
        switch self {
        case .home:
            return []
        case .recipes:
            return [
                HubAction(destination: .recipesCreate, title: "Create", icon: "plus")
            ]
        case .calorieTracker:
            return [
                HubAction(destination: .calorieGoals, title: "Goals", icon: "target")
            ]
        case .mealPlan:
            return [
                HubAction(destination: .mealPlanScan, title: "Scan", icon: "camera.viewfinder")
            ]
        case .insights:
            return [
                HubAction(destination: .insightsCalories, title: "Calories", icon: "flame.fill")
            ]
        }
    }
}

private enum HubMetrics {
    static let primaryButtonSize: CGFloat = 44
    static let primaryIconSize: CGFloat = 20
    static let primarySpacing: CGFloat = 14
    static let sideActionSize: CGFloat = 40
    static let sideActionIconSize: CGFloat = 18
    static let sideActionSpacing: CGFloat = 12
    static let sideActionOverlap: CGFloat = 14
    static let barHeight: CGFloat = 76
    static let barCornerRadius: CGFloat = 30
}
