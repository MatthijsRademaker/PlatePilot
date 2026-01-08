import SwiftUI

struct RecipeCardView: View {
    let recipe: Recipe

    var body: some View {
        HStack(spacing: 12) {
            RemoteImageView(url: recipe.listImageURL, cornerRadius: 16)
                .frame(width: 110, height: 90)

            VStack(alignment: .leading, spacing: 6) {
                Text(recipe.name)
                    .font(PlatePilotTheme.titleFont(size: 18))
                    .foregroundStyle(PlatePilotTheme.textPrimary)
                    .lineLimit(2)

                if let cuisine = recipe.cuisine {
                    Text(cuisine)
                        .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.accent)
                }

                if !recipe.description.isEmpty {
                    Text(recipe.description)
                        .font(PlatePilotTheme.bodyFont(size: 12))
                        .foregroundStyle(PlatePilotTheme.textSecondary)
                        .lineLimit(2)
                }
            }

            Spacer(minLength: 0)
        }
        .padding(12)
        .plateGlass(cornerRadius: PlatePilotMetrics.cardRadius, tint: .white.opacity(0.2), interactive: true)
    }
}
