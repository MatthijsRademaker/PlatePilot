import SwiftUI

struct WeekPlanView: View {
    let weekPlan: WeekPlan
    let onPrevious: () -> Void
    let onNext: () -> Void
    let onSlotTap: (MealSlot) -> Void
    let onSlotClear: (MealSlot) -> Void

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Button(action: onPrevious) {
                    Image(systemName: "chevron.left")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.accent)
                        .frame(width: 34, height: 34)
                        .background(PlatePilotTheme.tintWarm, in: RoundedRectangle(cornerRadius: 12, style: .continuous))
                }
                .buttonStyle(.plain)

                Spacer()

                Text(weekPlan.startDate.rangeLabel(to: weekPlan.endDate))
                    .font(PlatePilotTheme.bodyFont(size: 16, weight: .semibold))
                    .foregroundStyle(PlatePilotTheme.textPrimary)

                Spacer()

                Button(action: onNext) {
                    Image(systemName: "chevron.right")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.accent)
                        .frame(width: 34, height: 34)
                        .background(PlatePilotTheme.tintWarm, in: RoundedRectangle(cornerRadius: 12, style: .continuous))
                }
                .buttonStyle(.plain)
            }

            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 12) {
                    ForEach(weekPlan.days) { day in
                        DayPlanCardView(day: day, onSlotTap: onSlotTap, onSlotClear: onSlotClear)
                    }
                }
            }
        }
    }
}

private struct DayPlanCardView: View {
    let day: DayPlan
    let onSlotTap: (MealSlot) -> Void
    let onSlotClear: (MealSlot) -> Void

    var body: some View {
        VStack(spacing: 12) {
            VStack(spacing: 2) {
                Text(day.date.dayLabel())
                    .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                    .foregroundStyle(.white)
                Text(day.date.shortDate())
                    .font(PlatePilotTheme.bodyFont(size: 12, weight: .medium))
                    .foregroundStyle(.white.opacity(0.8))
            }
            .frame(maxWidth: .infinity)
            .padding(.vertical, 10)
            .background(PlatePilotTheme.headerGradient)

            VStack(spacing: 8) {
                ForEach(day.meals) { meal in
                    MealSlotCardView(slot: meal, onTap: { onSlotTap(meal) }, onClear: { onSlotClear(meal) })
                }
            }
            .padding(.horizontal, 8)
            .padding(.bottom, 10)
        }
        .frame(width: 200)
        .background(PlatePilotTheme.surface, in: RoundedRectangle(cornerRadius: 18, style: .continuous))
        .overlay(
            RoundedRectangle(cornerRadius: 18, style: .continuous)
                .stroke(Color.black.opacity(0.05), lineWidth: 1)
        )
    }
}
