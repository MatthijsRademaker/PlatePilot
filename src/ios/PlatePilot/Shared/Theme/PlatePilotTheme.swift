import SwiftUI

enum PlatePilotTheme {
    static let accent = Color(red: 1.0, green: 0.5, blue: 0.31)
    static let accentDeep = Color(red: 0.98, green: 0.39, blue: 0.28)
    static let background = Color(red: 1.0, green: 0.97, blue: 0.96)
    static let surface = Color.white
    static let textPrimary = Color(red: 0.18, green: 0.12, blue: 0.1)
    static let textSecondary = Color(red: 0.66, green: 0.62, blue: 0.62)
    static let tintWarm = Color(red: 1.0, green: 0.93, blue: 0.91)

    static let headerGradient = LinearGradient(
        colors: [accent, accentDeep],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )

    static let pageGradient = LinearGradient(
        colors: [background, .white],
        startPoint: .top,
        endPoint: .bottom
    )

    static func titleFont(size: CGFloat) -> Font {
        .system(size: size, weight: .semibold, design: .serif)
    }

    static func bodyFont(size: CGFloat, weight: Font.Weight = .regular) -> Font {
        .system(size: size, weight: weight, design: .rounded)
    }
}

enum PlatePilotMetrics {
    static let cardRadius: CGFloat = 20
    static let smallRadius: CGFloat = 14
}
