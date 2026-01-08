import SwiftUI

struct InsightsView: View {
    private let insightsAccent = Color(red: 0.4, green: 0.58, blue: 0.98)
    private let highlights: [InsightHighlight] = [
        InsightHighlight(
            title: "Weekly Balance",
            subtitle: "Macros and calories are steady this week.",
            icon: "chart.line.uptrend.xyaxis"
        ),
        InsightHighlight(
            title: "Favorite Cuisines",
            subtitle: "Italian and Thai are leading your plans.",
            icon: "fork.knife"
        ),
        InsightHighlight(
            title: "Consistency",
            subtitle: "4 of 7 planned meals completed.",
            icon: "checkmark.seal.fill"
        )
    ]

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                InsightsHeaderView()

                VStack(spacing: 12) {
                    ForEach(highlights) { highlight in
                        InsightCardView(highlight: highlight, accent: insightsAccent)
                    }
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
    }
}

private struct InsightsHeaderView: View {
    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(
                    LinearGradient(
                        colors: [Color(red: 0.4, green: 0.58, blue: 0.98), Color(red: 0.32, green: 0.46, blue: 0.9)],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                )
                .frame(maxWidth: .infinity)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                .padding(.bottom, 4)

            HStack(spacing: 12) {
                Image(systemName: "chart.bar.xaxis")
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .plateGlass(cornerRadius: 12, tint: .white.opacity(0.25))

                VStack(alignment: .leading, spacing: 4) {
                    Text("Insights")
                        .font(PlatePilotTheme.titleFont(size: 24))
                        .foregroundStyle(.white)

                    Text("Patterns and progress")
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

private struct InsightCardView: View {
    let highlight: InsightHighlight
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

private struct InsightHighlight: Identifiable {
    let id = UUID()
    let title: String
    let subtitle: String
    let icon: String
}
