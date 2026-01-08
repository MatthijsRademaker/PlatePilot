import SwiftUI

struct CalorieTrackerView: View {
    private let accent = Color(red: 1.0, green: 0.52, blue: 0.28)
    private let highlights: [TrackerHighlight] = [
        TrackerHighlight(
            title: "Move Goal",
            subtitle: "320 kcal left to hit today's target.",
            icon: "figure.run"
        ),
        TrackerHighlight(
            title: "Macros",
            subtitle: "Carbs 48%, Protein 28%, Fat 24%.",
            icon: "chart.pie"
        ),
        TrackerHighlight(
            title: "Hydration",
            subtitle: "1.4L logged - 0.6L to go.",
            icon: "drop.fill"
        )
    ]

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                CalorieTrackerHeaderView(accent: accent)

                VStack(spacing: 16) {
                    DailyCalorieTrackerView()

                    VStack(spacing: 12) {
                        ForEach(highlights) { highlight in
                            TrackerHighlightCard(highlight: highlight, accent: accent)
                        }
                    }
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
    }
}

private struct CalorieTrackerHeaderView: View {
    let accent: Color

    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(
                    LinearGradient(
                        colors: [accent, accent.opacity(0.75)],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                )
                .frame(maxWidth: .infinity)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                .padding(.bottom, 4)

            HStack(spacing: 12) {
                Image(systemName: "flame.fill")
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .plateGlass(cornerRadius: 12, tint: .white.opacity(0.25))

                VStack(alignment: .leading, spacing: 4) {
                    Text("Calorie Tracker")
                        .font(PlatePilotTheme.titleFont(size: 24))
                        .foregroundStyle(.white)

                    Text("Energy, balance, and daily goals")
                        .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                        .foregroundStyle(.white.opacity(0.85))
                }

                Spacer()
            }
            .padding(.horizontal, 32)
            .padding(.vertical, 22)
        }
        .padding(.top, 12)
    }
}

private struct TrackerHighlightCard: View {
    let highlight: TrackerHighlight
    let accent: Color

    var body: some View {
        HStack(spacing: 14) {
            Image(systemName: highlight.icon)
                .font(.system(size: 18, weight: .semibold))
                .foregroundStyle(accent)
                .frame(width: 44, height: 44)
                .background(PlatePilotTheme.tintWarm, in: RoundedRectangle(cornerRadius: 14, style: .continuous))

            VStack(alignment: .leading, spacing: 4) {
                Text(highlight.title)
                    .font(PlatePilotTheme.bodyFont(size: 15, weight: .semibold))
                    .foregroundStyle(PlatePilotTheme.textPrimary)

                Text(highlight.subtitle)
                    .font(PlatePilotTheme.bodyFont(size: 12))
                    .foregroundStyle(PlatePilotTheme.textSecondary)
            }

            Spacer()
        }
        .padding(16)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
    }
}

private struct TrackerHighlight: Identifiable {
    let id = UUID()
    let title: String
    let subtitle: String
    let icon: String
}
