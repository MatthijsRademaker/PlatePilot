import SwiftUI

struct GlassHubNavigationView: View {
    @Environment(AppState.self) private var appState
    @Environment(RouterPath.self) private var router
    @State private var isStackExpanded = false

    private var currentSection: HubSection {
        HubSection(tab: appState.selectedTab) ?? .home
    }

    private var primarySections: [HubSection] {
        [.recipes, .calorieTracker, .home, .mealPlan, .insights]
    }

    private var stackAnchorSection: HubSection {
        .mealPlan
    }

    var body: some View {
        PlateGlassGroup(spacing: 12) {
            ZStack(alignment: .bottom) {
                GlassPrimaryNavBar(
                    sections: primarySections,
                    activeSection: currentSection,
                    accent: currentSection.accent,
                    onSelect: selectPrimarySection
                )

                if !currentSection.stackActions.isEmpty {
                    HStack(spacing: 0) {
                        ForEach(primarySections) { section in
                            if section == stackAnchorSection {
                                GlassContextStackMenu(
                                    actions: currentSection.stackActions,
                                    accent: currentSection.accent,
                                    isExpanded: $isStackExpanded,
                                    onToggle: toggleStack,
                                    onAction: handleAction
                                )
                                .frame(maxWidth: .infinity)
                            } else {
                                Color.clear
                                    .frame(maxWidth: .infinity)
                                    .allowsHitTesting(false)
                            }
                        }
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.horizontal, HubMetrics.barHorizontalPadding)
                    .offset(x: HubMetrics.stackGap * 8, y: -HubMetrics.stackLift)
                    .zIndex(2)
                }
            }
            .frame(maxWidth: .infinity)
        }
        .padding(.horizontal, 12)
        .padding(.top, 6)
        .padding(.bottom, 4)
        .onChange(of: appState.selectedTab) { _, _ in
            closeStack()
        }
    }

    private func selectPrimarySection(_ section: HubSection) {
        if let tab = section.tab {
            appState.selectedTab = tab
        }
        closeStack()
    }

    private func handleAction(_ action: HubAction) {
        router.push(.hubDestination(action.destination))
        closeStack()
    }

    private func toggleStack() {
        withAnimation(.spring(response: 0.5, dampingFraction: 0.75, blendDuration: 0.2)) {
            isStackExpanded.toggle()
        }
    }

    private func closeStack() {
        guard isStackExpanded else { return }
        withAnimation(.spring(response: 0.45, dampingFraction: 0.8, blendDuration: 0.2)) {
            isStackExpanded = false
        }
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
        .padding(.horizontal, HubMetrics.barHorizontalPadding)
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

private struct GlassContextStackMenu: View {
    let actions: [HubAction]
    let accent: Color
    @Binding var isExpanded: Bool
    let onToggle: () -> Void
    let onAction: (HubAction) -> Void

    var body: some View {
        ZStack(alignment: .bottom) {
            if isExpanded {
                ForEach(Array(actions.enumerated()), id: \.element.id) { index, action in
                    let offset = fanOffset(index: index, total: actions.count)

                    Button {
                        onAction(action)
                    } label: {
                        StackActionBubble(
                            action: action,
                            accent: accent,
                            fadeInDelay: Double(index) * 0.25
                        )
                    }
                    .buttonStyle(.plain)
                    .offset(offset)
                    .transition(fanTransition(for: offset))
                    .animation(
                        .spring(response: 0.52, dampingFraction: 0.7, blendDuration: 0.2)
                            .delay(Double(index) * 0.05),
                        value: isExpanded
                    )
                    .zIndex(Double(actions.count - index))
                }
            }

            stackButton
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .bottom)
        .allowsHitTesting(isExpanded || !actions.isEmpty)
    }

    private var stackButton: some View {
        let iconName = isExpanded ? "xmark" : "square.stack.3d.up.fill"
        let labelText = isExpanded ? "Close menu" : "Open menu"

        return Button {
            onToggle()
        } label: {
            ZStack {
                Circle()
                    .fill(.white.opacity(0.12))

                Image(systemName: iconName)
                    .font(.system(size: HubMetrics.stackIconSize, weight: .semibold))
                    .foregroundStyle(.white)
            }
            .frame(width: HubMetrics.stackButtonSize, height: HubMetrics.stackButtonSize)
            .contentShape(Circle())
            .background(
                Circle()
                    .fill(
                        RadialGradient(
                            colors: [accent.opacity(0.6), accent.opacity(0.2)],
                            center: .topLeading,
                            startRadius: 4,
                            endRadius: HubMetrics.stackButtonSize
                        )
                    )
            )
            .overlay(
                Circle()
                    .stroke(Color.white.opacity(0.25), lineWidth: 0.8)
            )
            .plateGlass(
                cornerRadius: HubMetrics.stackButtonSize / 2, tint: accent.opacity(0.35),
                interactive: true
            )
            .shadow(color: accent.opacity(0.35), radius: 14, x: 0, y: 8)
        }
        .buttonStyle(.plain)
        .hoverEffect(.lift)
        .accessibilityLabel(labelText)
    }

    private func fanOffset(index: Int, total: Int) -> CGSize {
        let angles = arcAngles(for: total)
        let startAngle = angles.start
        let endAngle = angles.end
        let targetAngle: Double

        if total <= 1 {
            targetAngle = (startAngle + endAngle) / 2
        } else {
            let step = (endAngle - startAngle) / Double(total - 1)
            targetAngle = startAngle + step * Double(index)
        }

        let radians = targetAngle * Double.pi / 180
        return CGSize(
            width: cos(radians) * HubMetrics.stackRadius,
            height: sin(radians) * HubMetrics.stackRadius
        )
    }

    private func fanTransition(for offset: CGSize) -> AnyTransition {
        AnyTransition.modifier(
            active: FanTransitionModifier(
                extraOffset: CGSize(width: -offset.width, height: -offset.height),
                scale: 0.35,
                opacity: 0,
                rotation: -18,
                blur: 6
            ),
            identity: FanTransitionModifier(
                extraOffset: .zero,
                scale: 1,
                opacity: 1,
                rotation: 0,
                blur: 0
            )
        )
    }

    private struct FanTransitionModifier: ViewModifier {
        let extraOffset: CGSize
        let scale: CGFloat
        let opacity: Double
        let rotation: Double
        let blur: CGFloat

        func body(content: Content) -> some View {
            content
                .offset(extraOffset)
                .scaleEffect(scale)
                .opacity(opacity)
                .rotationEffect(.degrees(rotation))
                .blur(radius: blur)
        }
    }

    private func arcAngles(for total: Int) -> (start: Double, end: Double) {
        switch total {
        case 1:
            return (-120, -120)
        case 2:
            return (-160, -100)
        case 3:
            return (-150, -90)
        default:
            return (-150, -90)
        }
    }
}

private struct StackActionBubble: View {
    let action: HubAction
    let accent: Color
    let fadeInDelay: Double
    @State private var isVisible = false

    var body: some View {
        VStack(spacing: 4) {
            ZStack {
                Circle()
                    .fill(Color.white.opacity(0.1))

                Image(systemName: action.icon)
                    .font(.system(size: HubMetrics.stackIconSize, weight: .semibold))
                    .foregroundStyle(.white)
            }
            .frame(width: HubMetrics.stackItemSize, height: HubMetrics.stackItemSize)
            .contentShape(Circle())
            .plateGlass(
                cornerRadius: HubMetrics.stackItemSize / 2,
                tint: accent.opacity(0.25),
                interactive: true
            )
            .overlay(
                Circle()
                    .stroke(Color.white.opacity(0.25), lineWidth: 0.7)
            )
            .shadow(color: accent.opacity(0.3), radius: 10, x: 0, y: 6)
        }
        .opacity(isVisible ? 1 : 0)
        .onAppear {
            withAnimation(.easeOut(duration: 0.5).delay(fadeInDelay)) {
                isVisible = true
            }
        }
        .onDisappear {
            isVisible = false
        }
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

    var stackActions: [HubAction] {
        switch self {
        case .home:
            return []
        case .recipes:
            return [
                HubAction(destination: .recipesCreate, title: "Create", icon: "plus"),
                HubAction(destination: .recipesSearch, title: "Search", icon: "magnifyingglass"),
            ]
        case .calorieTracker:
            return [
                HubAction(
                    destination: .calorieAddExercise, title: "Add Exercise", icon: "figure.run"),
                HubAction(destination: .calorieGoals, title: "Goals", icon: "target"),
            ]
        case .mealPlan:
            return [
                HubAction(destination: .mealPlanSuggest, title: "Suggest", icon: "sparkles"),
                HubAction(destination: .mealPlanScan, title: "Scan", icon: "camera.viewfinder"),
            ]
        case .insights:
            return [
                HubAction(destination: .insightsRecipes, title: "Recipes", icon: "book.fill"),
                HubAction(destination: .insightsMealPlan, title: "Meal Plan", icon: "calendar"),
                HubAction(destination: .insightsCalories, title: "Calories", icon: "flame.fill"),
            ]
        }
    }
}

private enum HubMetrics {
    static let primaryButtonSize: CGFloat = 54
    static let primaryIconSize: CGFloat = 24
    static let barHeight: CGFloat = 76
    static let barHorizontalPadding: CGFloat = 12
    static let barCornerRadius: CGFloat = 30
    static let stackButtonSize: CGFloat = 54
    static let stackItemSize: CGFloat = 54
    static let stackIconSize: CGFloat = 24
    static let stackLabelSize: CGFloat = 10
    static let stackLabelWidth: CGFloat = 72
    static let stackRadius: CGFloat = 72
    static let stackGap: CGFloat = 8
    static let stackLift: CGFloat = barHeight + stackGap * 2
}
