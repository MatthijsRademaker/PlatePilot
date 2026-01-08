import SwiftUI

struct MealPlanView: View {
    @Environment(MealPlanStore.self) private var mealPlanStore
    @Environment(RecipeStore.self) private var recipeStore

    @State private var selectedSlot: MealSlot?
    @State private var searchQuery = ""

    var body: some View {
        ScrollView {
            VStack(spacing: 16) {
                MealPlanHeaderView(
                    onToday: { mealPlanStore.goToToday() },
                    onClear: { mealPlanStore.clearWeek() }
                )

                WeekPlanView(
                    weekPlan: mealPlanStore.currentWeek,
                    onPrevious: { mealPlanStore.shiftWeek(by: -1) },
                    onNext: { mealPlanStore.shiftWeek(by: 1) },
                    onSlotTap: { selectedSlot = $0 },
                    onSlotClear: { mealPlanStore.clearRecipe(for: $0.id) }
                )
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
        .sheet(item: $selectedSlot) { slot in
            RecipePickerSheet(
                slot: slot,
                searchQuery: $searchQuery,
                recipes: recipeStore.recipes,
                onSelect: { recipe in
                    mealPlanStore.setRecipe(recipe, for: slot.id)
                    selectedSlot = nil
                }
            )
        }
        .task {
            if recipeStore.recipes.isEmpty {
                await recipeStore.refresh()
            }
        }
    }
}

private struct MealPlanHeaderView: View {
    let onToday: () -> Void
    let onClear: () -> Void

    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(PlatePilotTheme.headerGradient)
                .frame(maxWidth: .infinity)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                .padding(.bottom, 4)

            HStack {
                HStack(spacing: 12) {
                    Image(systemName: "calendar")
                        .font(.system(size: 20, weight: .semibold))
                        .foregroundStyle(.white)
                        .frame(width: 44, height: 44)
                        .plateGlass(cornerRadius: 12, tint: .white.opacity(0.25))

                    Text("Meal Plan")
                        .font(PlatePilotTheme.titleFont(size: 24))
                        .foregroundStyle(.white)
                }

                Spacer()

                HStack(spacing: 8) {
                    Button("Today", action: onToday)
                        .font(PlatePilotTheme.bodyFont(size: 13, weight: .semibold))
                        .foregroundStyle(.white)
                        .padding(.horizontal, 12)
                        .padding(.vertical, 6)
                        .plateGlass(cornerRadius: 12, tint: .white.opacity(0.2), interactive: true)

                    Button(action: onClear) {
                        Image(systemName: "trash")
                            .font(.system(size: 14, weight: .semibold))
                            .foregroundStyle(.white)
                            .frame(width: 36, height: 36)
                            .plateGlass(cornerRadius: 12, tint: .white.opacity(0.2), interactive: true)
                    }
                    .buttonStyle(.plain)
                }
            }
            .padding(.horizontal, 32)
            .padding(.vertical, 22)
        }
        .padding(.top, 12)
    }
}

private struct RecipePickerSheet: View {
    let slot: MealSlot
    @Binding var searchQuery: String
    let recipes: [Recipe]
    let onSelect: (Recipe) -> Void
    @State private var suggestedRecipes: [Recipe] = []
    @State private var isSuggesting = false

    private var filteredRecipes: [Recipe] {
        if searchQuery.isEmpty { return recipes }
        return recipes.filter {
            $0.name.lowercased().contains(searchQuery.lowercased()) ||
            $0.description.lowercased().contains(searchQuery.lowercased())
        }
    }

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(alignment: .leading, spacing: 16) {
                    Text("Select a Recipe")
                        .font(PlatePilotTheme.titleFont(size: 22))
                        .foregroundStyle(PlatePilotTheme.textPrimary)

                    HStack(spacing: 12) {
                        TextField("Search recipes", text: $searchQuery)
                            .textFieldStyle(.roundedBorder)

                        SparkleSuggestButton(isLoading: isSuggesting, action: handleSuggest)
                    }

                    if !suggestedRecipes.isEmpty {
                        HStack {
                            Text("Suggested")
                                .font(PlatePilotTheme.bodyFont(size: 14, weight: .semibold))
                                .foregroundStyle(PlatePilotTheme.textSecondary)

                            Spacer()

                            Button("Clear") {
                                suggestedRecipes = []
                            }
                            .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                            .foregroundStyle(PlatePilotTheme.textSecondary)
                        }

                        VStack(spacing: 12) {
                            ForEach(suggestedRecipes) { recipe in
                                Button {
                                    onSelect(recipe)
                                } label: {
                                    RecipePickerRow(recipe: recipe, accent: true)
                                }
                                .buttonStyle(.plain)
                            }
                        }
                    }

                    Text("All Recipes")
                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .semibold))
                        .foregroundStyle(PlatePilotTheme.textSecondary)

                    VStack(spacing: 10) {
                        ForEach(filteredRecipes) { recipe in
                            Button {
                                onSelect(recipe)
                            } label: {
                                RecipePickerRow(recipe: recipe, accent: false)
                            }
                            .buttonStyle(.plain)
                        }
                    }
                }
                .padding(20)
            }
            .navigationTitle(slot.mealType.label)
            .navigationBarTitleDisplayMode(.inline)
        }
    }

    private func handleSuggest() {
        guard !isSuggesting else { return }
        isSuggesting = true
        suggestedRecipes = []

        Task { @MainActor in
            try? await Task.sleep(nanoseconds: 750_000_000)
            let shuffled = recipes.shuffled()
            suggestedRecipes = Array(shuffled.prefix(3))
            isSuggesting = false
        }
    }
}

private struct SparkleSuggestButton: View {
    let isLoading: Bool
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            ZStack {
                if isLoading {
                    SparkleGlowRing(size: 54, lineWidth: 3, duration: 1.2, reverse: false)
                    SparkleGlowRing(size: 68, lineWidth: 2, duration: 1.8, reverse: true)
                }

                Image(systemName: "sparkles")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundStyle(.white)
                    .frame(width: 44, height: 44)
                    .background(PlatePilotTheme.headerGradient, in: RoundedRectangle(cornerRadius: 14, style: .continuous))
                    .shadow(color: PlatePilotTheme.accent.opacity(0.35), radius: 8, x: 0, y: 4)
            }
            .frame(width: 68, height: 68)
        }
        .buttonStyle(.plain)
        .disabled(isLoading)
        .accessibilityLabel(isLoading ? "Suggesting recipes" : "Suggest recipes")
    }
}

private struct SparkleGlowRing: View {
    let size: CGFloat
    let lineWidth: CGFloat
    let duration: Double
    let reverse: Bool
    @State private var spin = false

    var body: some View {
        Circle()
            .stroke(
                AngularGradient(
                    colors: [
                        PlatePilotTheme.accent.opacity(0.1),
                        PlatePilotTheme.accent.opacity(0.6),
                        PlatePilotTheme.accentDeep.opacity(0.9),
                        PlatePilotTheme.accent.opacity(0.1)
                    ],
                    center: .center
                ),
                style: StrokeStyle(lineWidth: lineWidth, lineCap: .round)
            )
            .frame(width: size, height: size)
            .rotationEffect(.degrees(spin ? (reverse ? -360 : 360) : 0))
            .animation(.linear(duration: duration).repeatForever(autoreverses: false), value: spin)
            .onAppear { spin = true }
    }
}

private struct RecipePickerRow: View {
    let recipe: Recipe
    let accent: Bool

    var body: some View {
        HStack(spacing: 12) {
            Image(systemName: accent ? "sparkles" : "fork.knife")
                .font(.system(size: 16, weight: .semibold))
                .foregroundStyle(.white)
                .frame(width: 32, height: 32)
                .background(
                    accent ? PlatePilotTheme.headerGradient : LinearGradient(colors: [PlatePilotTheme.accent, PlatePilotTheme.accentDeep], startPoint: .topLeading, endPoint: .bottomTrailing),
                    in: RoundedRectangle(cornerRadius: 10, style: .continuous)
                )

            VStack(alignment: .leading, spacing: 4) {
                Text(recipe.name)
                    .font(PlatePilotTheme.bodyFont(size: 15, weight: .semibold))
                    .foregroundStyle(PlatePilotTheme.textPrimary)
                Text(recipe.description)
                    .font(PlatePilotTheme.bodyFont(size: 12))
                    .foregroundStyle(PlatePilotTheme.textSecondary)
                    .lineLimit(1)
            }

            Spacer()
        }
        .padding(12)
        .plateGlass(cornerRadius: 16, tint: .white.opacity(0.15), interactive: true)
    }
}
