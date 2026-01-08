import SwiftUI

struct MealSlotCardView: View {
    let slot: MealSlot
    let onTap: () -> Void
    let onClear: () -> Void

    var body: some View {
        Button(action: onTap) {
            VStack(alignment: .leading, spacing: 8) {
                HStack(spacing: 8) {
                    Image(systemName: slot.mealType.systemImage)
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(slot.recipe == nil ? PlatePilotTheme.textSecondary : .white)
                        .frame(width: 22, height: 22)
                        .background(iconBackground, in: RoundedRectangle(cornerRadius: 6, style: .continuous))

                    Text(slot.mealType.label)
                        .font(PlatePilotTheme.bodyFont(size: 11, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.textSecondary)
                        .textCase(.uppercase)

                    Spacer()

                    if slot.recipe != nil {
                        Button(action: onClear) {
                            Image(systemName: "xmark")
                                .font(.system(size: 10, weight: .bold))
                                .foregroundStyle(PlatePilotTheme.textSecondary)
                                .frame(width: 18, height: 18)
                                .background(Color.white.opacity(0.9), in: RoundedRectangle(cornerRadius: 6, style: .continuous))
                        }
                        .buttonStyle(.plain)
                    }
                }

                if let recipe = slot.recipe {
                    Text(recipe.name)
                        .font(PlatePilotTheme.bodyFont(size: 13, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.textPrimary)
                        .lineLimit(1)
                } else {
                    HStack(spacing: 4) {
                        Image(systemName: "plus")
                            .font(.system(size: 11, weight: .semibold))
                        Text("Add recipe")
                            .font(PlatePilotTheme.bodyFont(size: 12, weight: .medium))
                    }
                    .foregroundStyle(PlatePilotTheme.textSecondary)
                }
            }
            .padding(10)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(cardBackground)
        }
        .buttonStyle(.plain)
    }

    @ViewBuilder
    private var cardBackground: some View {
        if slot.recipe == nil {
            RoundedRectangle(cornerRadius: 14, style: .continuous)
                .strokeBorder(
                    PlatePilotTheme.accent.opacity(0.35),
                    style: StrokeStyle(lineWidth: 1.5, dash: [5, 4])
                )
                .background(
                    RoundedRectangle(cornerRadius: 14, style: .continuous)
                        .fill(PlatePilotTheme.tintWarm)
                )
        } else {
            RoundedRectangle(cornerRadius: 14, style: .continuous)
                .fill(Color.white)
                .shadow(color: Color.black.opacity(0.05), radius: 4, x: 0, y: 2)
        }
    }

    private var iconBackground: AnyShapeStyle {
        if slot.recipe == nil {
            return AnyShapeStyle(PlatePilotTheme.tintWarm)
        }
        return AnyShapeStyle(PlatePilotTheme.headerGradient)
    }
}
