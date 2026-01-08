import SwiftUI

struct PlaceholderDetailView: View {
    let title: String
    let subtitle: String
    let icon: String
    let accent: Color

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                header

                VStack(spacing: 16) {
                    Image(systemName: icon)
                        .font(.system(size: 34, weight: .semibold))
                        .foregroundStyle(accent)
                        .frame(width: 76, height: 76)
                        .background(accent.opacity(0.12), in: RoundedRectangle(cornerRadius: 22, style: .continuous))

                    Text(title)
                        .font(PlatePilotTheme.titleFont(size: 20))
                        .foregroundStyle(PlatePilotTheme.textPrimary)

                    Text(subtitle)
                        .font(PlatePilotTheme.bodyFont(size: 14))
                        .foregroundStyle(PlatePilotTheme.textSecondary)
                        .multilineTextAlignment(.center)
                }
                .frame(maxWidth: .infinity)
                .padding(24)
                .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
            }
            .padding(.horizontal, 16)
            .padding(.bottom, 24)
        }
        .background(PlatePilotTheme.pageGradient)
    }

    private var header: some View {
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
                Image(systemName: icon)
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .plateGlass(cornerRadius: 12, tint: .white.opacity(0.25))

                VStack(alignment: .leading, spacing: 4) {
                    Text(title)
                        .font(PlatePilotTheme.titleFont(size: 24))
                        .foregroundStyle(.white)

                    Text("Scaffolded screen")
                        .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                        .foregroundStyle(.white.opacity(0.8))
                }

                Spacer()
            }
            .padding(.horizontal, 32)
            .padding(.vertical, 22)
        }
        .padding(.top, 12)
    }
}
