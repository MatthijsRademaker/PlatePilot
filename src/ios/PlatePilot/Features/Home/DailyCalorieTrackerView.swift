import SwiftUI

struct DailyCalorieTrackerView: View {
    private let targetCalories = 2200
    private let currentCalories = 1480

    private var progress: Double {
        Double(currentCalories) / Double(targetCalories)
    }

    private let mealBreakdown: [MealBreakdown] = [
        MealBreakdown(name: "Stir Fry", calories: 450, progress: 0.75, color: .orange, imageSeed: "stirfry"),
        MealBreakdown(name: "Salad", calories: 280, progress: 0.6, color: .green, imageSeed: "salad"),
        MealBreakdown(name: "Smoothie", calories: 320, progress: 0.85, color: .purple, imageSeed: "smoothie")
    ]

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Daily Calorie Tracker")
                .font(PlatePilotTheme.titleFont(size: 18))
                .foregroundStyle(PlatePilotTheme.textPrimary)

            HStack(spacing: 20) {
                ZStack {
                    RingView(value: progress, size: 90, lineWidth: 10, color: PlatePilotTheme.accent)
                    VStack(spacing: 2) {
                        Text("\(currentCalories)")
                            .font(PlatePilotTheme.bodyFont(size: 20, weight: .bold))
                            .foregroundStyle(PlatePilotTheme.textPrimary)
                        Text("/ \(targetCalories)")
                            .font(PlatePilotTheme.bodyFont(size: 11, weight: .medium))
                            .foregroundStyle(PlatePilotTheme.textSecondary)
                    }
                }

                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 16) {
                        ForEach(mealBreakdown) { meal in
                            VStack(spacing: 8) {
                                ZStack {
                                    RingView(value: meal.progress, size: 52, lineWidth: 6, color: meal.color)

                                    RemoteImageView(
                                        url: URL(string: "https://picsum.photos/seed/\(meal.imageSeed)/100/100"),
                                        cornerRadius: 18
                                    )
                                    .frame(width: 36, height: 36)
                                }

                                Text(meal.name)
                                    .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                                    .foregroundStyle(PlatePilotTheme.textPrimary)
                                    .lineLimit(1)

                                Text("\(meal.calories) cal")
                                    .font(PlatePilotTheme.bodyFont(size: 11, weight: .medium))
                                    .foregroundStyle(PlatePilotTheme.textSecondary)
                            }
                            .frame(width: 72)
                        }
                    }
                }
            }
        }
        .padding(20)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2))
    }
}

private struct MealBreakdown: Identifiable {
    let id = UUID()
    let name: String
    let calories: Int
    let progress: Double
    let color: Color
    let imageSeed: String
}

private struct RingView: View {
    let value: Double
    let size: CGFloat
    let lineWidth: CGFloat
    let color: Color

    var body: some View {
        ZStack {
            Circle()
                .stroke(Color.gray.opacity(0.2), lineWidth: lineWidth)

            Circle()
                .trim(from: 0, to: max(0, min(value, 1)))
                .stroke(color, style: StrokeStyle(lineWidth: lineWidth, lineCap: .round))
                .rotationEffect(.degrees(-90))
        }
        .frame(width: size, height: size)
    }
}
