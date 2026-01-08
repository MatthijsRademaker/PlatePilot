import SwiftUI

struct GlassHubNavigationView: View {
    @Environment(AppState.self) private var appState

    @State private var isExpanded = false
    @State private var activeSection: HubSection = .calorieTracker
    @State private var hubPulse = false

    var body: some View {
        PlateGlassGroup(spacing: 24) {
            ZStack(alignment: .bottom) {
                GlassHubBottomBar(section: activeSection)

                GlassHubOrbitFan(
                    sections: HubSection.allCases,
                    isExpanded: $isExpanded,
                    activeSection: $activeSection,
                    hubPulse: $hubPulse,
                    onSelect: selectSection,
                    onToggle: toggleExpanded,
                    onReset: resetHub
                )
                .frame(height: HubMetrics.fanHeight)
            }
            .frame(maxWidth: .infinity)
            .frame(height: HubMetrics.fanHeight + HubMetrics.barHeight)
        }
        .padding(.horizontal, 16)
        .padding(.top, 6)
        .padding(.bottom, 4)
        .onAppear {
            hubPulse = true
            withAnimation(.spring(response: 0.6, dampingFraction: 0.7, blendDuration: 0.2)) {
                isExpanded = true
            }
        }
    }

    private func toggleExpanded() {
        withAnimation(.spring(response: 0.6, dampingFraction: 0.7, blendDuration: 0.2)) {
            isExpanded.toggle()
        }
    }

    private func resetHub() {
        withAnimation(.spring(response: 0.6, dampingFraction: 0.8, blendDuration: 0.2)) {
            isExpanded = false
            activeSection = .calorieTracker
        }
        appState.selectedTab = .home
    }

    private func selectSection(_ section: HubSection) {
        activeSection = section
        if let tab = section.tab {
            appState.selectedTab = tab
        }
    }
}

private struct GlassHubOrbitFan: View {
    let sections: [HubSection]
    @Binding var isExpanded: Bool
    @Binding var activeSection: HubSection
    @Binding var hubPulse: Bool
    let onSelect: (HubSection) -> Void
    let onToggle: () -> Void
    let onReset: () -> Void

    var body: some View {
        ZStack(alignment: .bottom) {
            GlassHubGlow(accent: activeSection.accent)
                .opacity(isExpanded ? 1 : 0)
                .animation(.easeInOut(duration: 0.4), value: isExpanded)

            ZStack {
                if isExpanded {
                    connector(for: activeSection)
                }

                ForEach(Array(sections.enumerated()), id: \.element.id) { index, section in
                    let offset = orbitOffset(for: index)

                    Button {
                        withAnimation(.spring(response: 0.5, dampingFraction: 0.75, blendDuration: 0.2)) {
                            onSelect(section)
                        }
                    } label: {
                        HubOrbitBubble(section: section, isSelected: section == activeSection)
                    }
                    .buttonStyle(.plain)
                    .offset(isExpanded ? offset : .zero)
                    .opacity(isExpanded ? 1 : 0)
                    .scaleEffect(isExpanded ? (section == activeSection ? 1.08 : 1.0) : 0.7)
                    .animation(
                        .spring(response: 0.5, dampingFraction: 0.72, blendDuration: 0.2)
                            .delay(Double(index) * 0.05),
                        value: isExpanded
                    )
                    .animation(.spring(response: 0.4, dampingFraction: 0.8, blendDuration: 0.2), value: activeSection)
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .bottom)
            .offset(y: -HubMetrics.hubSize / 2)

            hubButton
                .scaleEffect(hubPulse ? 1.02 : 1.0)
                .animation(.easeInOut(duration: 2).repeatForever(autoreverses: true), value: hubPulse)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .bottom)
    }

    private var hubButton: some View {
        let accent = activeSection.accent
        let doubleTap = TapGesture(count: 2).onEnded(onReset)
        let singleTap = TapGesture(count: 1).onEnded(onToggle)

        return ZStack {
            Circle()
                .fill(.white.opacity(0.1))

            Image(systemName: "house.fill")
                .font(.system(size: 24, weight: .semibold))
                .foregroundStyle(.white)
        }
        .frame(width: HubMetrics.hubSize, height: HubMetrics.hubSize)
        .background(
            Circle()
                .fill(
                    RadialGradient(
                        colors: [accent.opacity(0.6), accent.opacity(0.15)],
                        center: .topLeading,
                        startRadius: 4,
                        endRadius: HubMetrics.hubSize
                    )
                )
        )
        .overlay(
            LinearGradient(
                colors: [Color.white.opacity(0.7), .clear],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .clipShape(Circle())
            .opacity(0.7)
        )
        .overlay(
            Circle()
                .stroke(Color.white.opacity(0.22), lineWidth: 1)
        )
        .plateGlass(cornerRadius: HubMetrics.hubSize / 2, tint: accent.opacity(0.35), interactive: true)
        .shadow(color: accent.opacity(0.45), radius: 18, x: 0, y: 10)
        .accessibilityLabel("Open navigation hub")
        .accessibilityAddTraits(.isButton)
        .gesture(doubleTap.exclusively(before: singleTap))
    }

    private func connector(for section: HubSection) -> some View {
        let offset = orbitOffset(for: sections.firstIndex(of: section) ?? 0)
        let distance = max(1, hypot(offset.width, offset.height))
        let hubRadius = HubMetrics.hubSize / 2
        let bubbleRadius = HubMetrics.bubbleSize / 2
        let length = max(12, distance - hubRadius - bubbleRadius + HubMetrics.connectorOverlap)
        let unitX = offset.width / distance
        let unitY = offset.height / distance
        let centerDistance = (distance + hubRadius - bubbleRadius) / 2
        let centerX = unitX * centerDistance
        let centerY = unitY * centerDistance
        let angle = Angle(radians: atan2(offset.height, offset.width))

        return LiquidConnector(
            length: length,
            thickness: HubMetrics.connectorThickness,
            accent: section.accent
        )
        .rotationEffect(angle)
        .offset(x: centerX, y: centerY)
        .animation(.spring(response: 0.4, dampingFraction: 0.75, blendDuration: 0.2), value: activeSection)
    }

    private func orbitOffset(for index: Int) -> CGSize {
        guard sections.count > 1 else {
            return CGSize(width: 0, height: -HubMetrics.orbitRadius)
        }

        let startAngle = -150.0
        let endAngle = -30.0
        let step = (endAngle - startAngle) / Double(sections.count - 1)
        let angle = (startAngle + step * Double(index)) * Double.pi / 180

        return CGSize(
            width: cos(angle) * HubMetrics.orbitRadius,
            height: sin(angle) * HubMetrics.orbitRadius
        )
    }
}

private struct HubOrbitBubble: View {
    let section: HubSection
    let isSelected: Bool

    var body: some View {
        VStack(spacing: 6) {
            ZStack {
                Circle()
                    .fill(.white.opacity(0.08))

                Image(systemName: section.icon)
                    .font(.system(size: 16, weight: .semibold))
                    .foregroundStyle(.white)
            }
            .frame(width: HubMetrics.bubbleSize, height: HubMetrics.bubbleSize)
            .background(
                Circle()
                    .fill(
                        RadialGradient(
                            colors: [section.accent.opacity(0.5), section.accent.opacity(0.1)],
                            center: .topLeading,
                            startRadius: 4,
                            endRadius: HubMetrics.bubbleSize
                        )
                    )
            )
            .overlay(
                LinearGradient(
                    colors: [Color.white.opacity(0.6), .clear],
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
                .clipShape(Circle())
                .opacity(0.65)
            )
            .overlay(
                Circle()
                    .stroke(isSelected ? section.accent.opacity(0.7) : Color.white.opacity(0.22), lineWidth: 1)
            )
            .plateGlass(
                cornerRadius: HubMetrics.bubbleSize / 2,
                tint: section.accent.opacity(isSelected ? 0.35 : 0.22),
                interactive: true
            )
            .shadow(
                color: section.accent.opacity(isSelected ? 0.35 : 0.2),
                radius: isSelected ? 12 : 8,
                x: 0,
                y: 6
            )

            Text(section.title)
                .font(PlatePilotTheme.bodyFont(size: 11, weight: .semibold))
                .foregroundStyle(isSelected ? PlatePilotTheme.textPrimary : PlatePilotTheme.textSecondary)
                .multilineTextAlignment(.center)
                .lineLimit(2)
        }
    }
}

private struct GlassHubBottomBar: View {
    let section: HubSection

    var body: some View {
        let accent = section.accent

        HStack {
            HubActionButton(action: section.leftAction, accent: accent)
            Spacer()
            Color.clear.frame(width: HubMetrics.hubSize)
            Spacer()
            HubActionButton(action: section.rightAction, accent: accent)
        }
        .padding(.horizontal, 18)
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

private struct HubActionButton: View {
    let action: HubAction
    let accent: Color

    var body: some View {
        Button {
        } label: {
            HStack(spacing: 6) {
                Image(systemName: action.icon)
                    .font(.system(size: 12, weight: .semibold))
                    .foregroundStyle(accent)
                Text(action.title)
                    .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                    .foregroundStyle(PlatePilotTheme.textPrimary)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 8)
        }
        .buttonStyle(.plain)
        .background(
            RoundedRectangle(cornerRadius: 14, style: .continuous)
                .fill(Color.white.opacity(0.18))
        )
        .plateGlass(cornerRadius: 14, tint: accent.opacity(0.18), interactive: true)
        .overlay(
            RoundedRectangle(cornerRadius: 14, style: .continuous)
                .stroke(Color.white.opacity(0.2), lineWidth: 0.6)
        )
        .shadow(color: accent.opacity(0.12), radius: 6, x: 0, y: 4)
    }
}

private struct GlassHubGlow: View {
    let accent: Color

    var body: some View {
        Ellipse()
            .fill(
                RadialGradient(
                    colors: [accent.opacity(0.35), accent.opacity(0.08), .clear],
                    center: .center,
                    startRadius: 0,
                    endRadius: HubMetrics.glowWidth / 1.6
                )
            )
            .frame(width: HubMetrics.glowWidth, height: HubMetrics.glowHeight)
            .blur(radius: 8)
            .offset(y: HubMetrics.glowOffset)
    }
}

private struct LiquidConnector: View {
    let length: CGFloat
    let thickness: CGFloat
    let accent: Color

    var body: some View {
        ZStack {
            Capsule()
                .fill(Color.white.opacity(0.14))
                .frame(width: length, height: thickness)
                .plateGlass(cornerRadius: thickness / 2, tint: accent.opacity(0.25))
                .overlay(
                    LinearGradient(
                        colors: [Color.white.opacity(0.65), .clear],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                    .clipShape(Capsule())
                    .opacity(0.7)
                )
                .overlay(
                    Capsule()
                        .stroke(Color.white.opacity(0.22), lineWidth: 0.7)
                )
                .shadow(color: accent.opacity(0.18), radius: 10, x: 0, y: 6)

            Circle()
                .fill(accent.opacity(0.22))
                .frame(width: thickness * 0.9, height: thickness * 0.9)
                .offset(x: length / 2 - thickness * 0.6)
                .blur(radius: 0.6)
        }
        .blur(radius: 0.4)
    }
}

private struct HubAction: Hashable {
    let title: String
    let icon: String
}

private enum HubSection: String, CaseIterable, Identifiable {
    case recipes
    case calorieTracker
    case mealPlan
    case insights

    var id: String { rawValue }

    var title: String {
        switch self {
        case .recipes: return "Recipes"
        case .calorieTracker: return "Calorie Tracker"
        case .mealPlan: return "Meal Plan"
        case .insights: return "Insights"
        }
    }

    var icon: String {
        switch self {
        case .recipes: return "book.fill"
        case .calorieTracker: return "flame.fill"
        case .mealPlan: return "calendar"
        case .insights: return "chart.bar.xaxis"
        }
    }

    var accent: Color {
        switch self {
        case .recipes: return Color(red: 1.0, green: 0.63, blue: 0.34)
        case .calorieTracker: return Color(red: 1.0, green: 0.52, blue: 0.28)
        case .mealPlan: return Color(red: 0.33, green: 0.82, blue: 0.52)
        case .insights: return Color(red: 0.4, green: 0.58, blue: 0.98)
        }
    }

    var tab: AppTab? {
        switch self {
        case .recipes: return .recipes
        case .calorieTracker: return .home
        case .mealPlan: return .mealPlan
        case .insights: return .search
        }
    }

    var leftAction: HubAction {
        switch self {
        case .recipes: return HubAction(title: "Discover", icon: "sparkles")
        case .calorieTracker: return HubAction(title: "Log", icon: "flame.fill")
        case .mealPlan: return HubAction(title: "Plan", icon: "calendar.badge.plus")
        case .insights: return HubAction(title: "Trends", icon: "chart.line.uptrend.xyaxis")
        }
    }

    var rightAction: HubAction {
        switch self {
        case .recipes: return HubAction(title: "Saved", icon: "bookmark.fill")
        case .calorieTracker: return HubAction(title: "History", icon: "clock.arrow.circlepath")
        case .mealPlan: return HubAction(title: "Groceries", icon: "cart.fill")
        case .insights: return HubAction(title: "Weekly", icon: "calendar")
        }
    }
}

private enum HubMetrics {
    static let hubSize: CGFloat = 64
    static let bubbleSize: CGFloat = 46
    static let orbitRadius: CGFloat = 112
    static let barHeight: CGFloat = 64
    static let barCornerRadius: CGFloat = 28
    static let connectorThickness: CGFloat = 14
    static let connectorOverlap: CGFloat = 12
    static let glowWidth: CGFloat = 240
    static let glowHeight: CGFloat = 120
    static let glowOffset: CGFloat = -6
    static let fanHeight: CGFloat = orbitRadius + bubbleSize / 2 + hubSize / 2 + 28
}
