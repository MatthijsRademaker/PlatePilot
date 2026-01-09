import PhotosUI
import SwiftUI
import UIKit

struct RecipeCreateView: View {
    @Environment(RecipeStore.self) private var recipeStore
    @Environment(\.dismiss) private var dismiss
    @Environment(\.accessibilityReduceTransparency) private var reduceTransparency

    @State private var draft = RecipeDraft()
    @State private var isSaving = false
    @State private var errorMessage: String?
    @State private var isAddingIngredient = false
    @State private var isAddingInstruction = false
    @State private var newIngredientName = ""
    @State private var newInstructionText = ""
    @State private var newIngredientQuantity = RecipeCreationMetrics.defaultQuantity
    @State private var newIngredientUnit = RecipeCreationMetrics.defaultUnit
    @State private var availableUnits: [String] = RecipeCreationMetrics.defaultUnits
    @State private var isUnitOverlayVisible = false
    @State private var unitOverlayTarget: UnitPickerTarget?
    @State private var newUnitName = ""
    @State private var timePickerTarget: TimePickerTarget?
    @State private var selectedPhoto: PhotosPickerItem?
    @State private var photoImage: Image?
    @State private var didAppear = false
    @FocusState private var focusedField: FocusField?

    var body: some View {
        ZStack {
            RecipeCreationBackground(reduceTransparency: reduceTransparency)

            ScrollView {
                VStack(spacing: RecipeCreationMetrics.sectionSpacing) {
                    header
                        .opacity(didAppear ? 1 : 0)
                        .offset(y: didAppear ? 0 : 12)
                        .animation(.easeOut(duration: 0.45), value: didAppear)

                    heroCard
                        .modifier(sectionAnimation(delay: 0.05))

                    metaCard
                        .modifier(sectionAnimation(delay: 0.1))

                    ingredientSection
                        .modifier(sectionAnimation(delay: 0.15))

                    instructionSection
                        .modifier(sectionAnimation(delay: 0.2))

                    guidedModeRow
                        .modifier(sectionAnimation(delay: 0.25))

                    dietTagsRow
                        .modifier(sectionAnimation(delay: 0.3))

                    if let errorMessage {
                        errorBanner(message: errorMessage)
                            .modifier(sectionAnimation(delay: 0.35))
                    }
                }
                .padding(.horizontal, RecipeCreationMetrics.horizontalPadding)
                .padding(.top, 28)
                .padding(.bottom, 120)
            }
            .scrollDismissesKeyboard(.immediately)

        }
        .safeAreaInset(edge: .bottom) {
            saveBar
                .padding(.bottom, RecipeCreationMetrics.hubBarHeight + 8)
        }
        .overlay {
            if isUnitOverlayVisible {
                UnitCreationOverlay(
                    unitName: $newUnitName,
                    onCancel: { hideUnitOverlay() },
                    onSave: { saveUnit() }
                )
                .transition(.opacity)
            }
        }
        .onAppear {
            withAnimation(.easeOut(duration: 0.4)) {
                didAppear = true
            }
        }
        .onChange(of: selectedPhoto) { _, newValue in
            guard let newValue else {
                photoImage = nil
                return
            }
            Task {
                if let data = try? await newValue.loadTransferable(type: Data.self),
                   let uiImage = UIImage(data: data) {
                    await MainActor.run {
                        photoImage = Image(uiImage: uiImage)
                    }
                }
            }
        }
        .sheet(item: $timePickerTarget) { target in
            TimePickerSheet(minutes: target.binding(for: $draft), title: target.title)
        }
    }

    private var header: some View {
        VStack(spacing: 8) {
            Text("New Recipe")
                .font(PlatePilotTheme.titleFont(size: 32))
                .tracking(-0.2)
                .foregroundStyle(RecipeCreationColors.textPrimary)

            Text("Let's build something delicious")
                .font(PlatePilotTheme.bodyFont(size: 16, weight: .medium))
                .foregroundStyle(RecipeCreationColors.textSecondary)
        }
        .frame(maxWidth: .infinity, alignment: .center)
    }

    private var heroCard: some View {
        GlassCard(isFocused: focusedField == .title) {
            HStack(alignment: .center, spacing: 16) {
                TextField("Recipe Title", text: $draft.title)
                    .font(PlatePilotTheme.titleFont(size: 22))
                    .foregroundStyle(RecipeCreationColors.textPrimary)
                    .focused($focusedField, equals: .title)
                    .submitLabel(.next)
                    .onSubmit { focusedField = .description }

                Spacer(minLength: 12)

                PhotosPicker(selection: $selectedPhoto, matching: .images) {
                    PhotoPickerButton(photoImage: photoImage)
                }
                .buttonStyle(.plain)
            }
        }
    }

    private var metaCard: some View {
        GlassCard(isFocused: focusedField == .description) {
            VStack(alignment: .leading, spacing: 16) {
                PlaceholderTextEditor(
                    text: $draft.description,
                    placeholder: "Description",
                    isFocused: focusedField == .description
                )
                .frame(minHeight: 90)
                .focused($focusedField, equals: .description)

                TimeRow(
                    prepMinutes: draft.prepMinutes,
                    cookMinutes: draft.cookMinutes,
                    onSelectPrep: { timePickerTarget = .prep },
                    onSelectCook: { timePickerTarget = .cook }
                )
            }
        }
    }

    private var ingredientSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            SectionHeader(title: "Ingredients", actionLabel: "Add ingredient") {
                beginIngredientEntry()
            }

            VStack(spacing: 10) {
                ForEach($draft.ingredients) { $ingredient in
                    IngredientRow(
                        ingredient: $ingredient,
                        quantityOptions: RecipeCreationMetrics.quantityOptions,
                        unitOptions: availableUnits,
                        onAddUnit: { showUnitOverlay(target: .ingredient(ingredient.id)) },
                        onDelete: { removeIngredient(id: ingredient.id) }
                    )
                }

                if isAddingIngredient {
                    ingredientInputRow
                } else if draft.ingredients.isEmpty {
                    EmptyRow(text: "Add your ingredients")
                }
            }
        }
    }

    private var instructionSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            SectionHeader(title: "Instructions", actionLabel: "Add step") {
                beginInstructionEntry()
            }

            VStack(spacing: 10) {
                ForEach(Array(draft.instructions.enumerated()), id: \.element.id) { index, step in
                    InstructionRow(index: index + 1, text: step.text, onDelete: {
                        removeInstruction(step)
                    })
                }

                if isAddingInstruction {
                    instructionInputRow
                } else if draft.instructions.isEmpty {
                    EmptyRow(text: "Add your first step")
                }
            }
        }
    }

    private var ingredientInputRow: some View {
        GlassRow {
            VStack(alignment: .leading, spacing: 12) {
                HStack(spacing: 10) {
                    TextField("Add ingredient", text: $newIngredientName)
                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                        .foregroundStyle(RecipeCreationColors.textPrimary)
                        .textInputAutocapitalization(.words)
                        .autocorrectionDisabled()
                        .focused($focusedField, equals: .newIngredient)
                        .submitLabel(.done)
                        .onSubmit { addIngredient() }

                    Spacer(minLength: 8)

                    Button(action: addIngredient) {
                        Image(systemName: "plus")
                            .font(.system(size: 12, weight: .bold))
                            .foregroundStyle(.white)
                            .frame(width: 28, height: 28)
                            .background(RecipeCreationColors.saveGradient, in: Circle())
                    }
                    .buttonStyle(.plain)
                    .disabled(newIngredientName.trimmed().isEmpty)
                    .accessibilityLabel("Add ingredient")

                    Button(action: cancelIngredientEntry) {
                        Image(systemName: "xmark")
                            .font(.system(size: 11, weight: .bold))
                            .foregroundStyle(RecipeCreationColors.textSecondary)
                            .frame(width: 28, height: 28)
                            .background(Color.white.opacity(0.2), in: Circle())
                    }
                    .buttonStyle(.plain)
                    .accessibilityLabel("Cancel ingredient")
                }

                HStack(spacing: 10) {
                    SelectBox(
                        title: "Quantity",
                        selection: $newIngredientQuantity,
                        options: RecipeCreationMetrics.quantityOptions
                    )

                    HStack(spacing: 6) {
                        SelectBox(
                            title: "Unit",
                            selection: $newIngredientUnit,
                            options: availableUnits
                        )

                        UnitAddButton(action: { showUnitOverlay(target: .newIngredient) })
                    }
                }
            }
        }
    }

    private var instructionInputRow: some View {
        GlassRow {
            HStack(alignment: .top, spacing: 10) {
                TextField("Add step", text: $newInstructionText, axis: .vertical)
                    .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                    .foregroundStyle(RecipeCreationColors.textPrimary)
                    .textInputAutocapitalization(.sentences)
                    .lineLimit(1...3)
                    .focused($focusedField, equals: .newInstruction)
                    .submitLabel(.done)
                    .onSubmit { addInstruction() }

                Spacer(minLength: 8)

                Button(action: addInstruction) {
                    Image(systemName: "plus")
                        .font(.system(size: 12, weight: .bold))
                        .foregroundStyle(.white)
                        .frame(width: 28, height: 28)
                        .background(RecipeCreationColors.saveGradient, in: Circle())
                }
                .buttonStyle(.plain)
                .disabled(newInstructionText.trimmed().isEmpty)
                .accessibilityLabel("Add instruction")

                Button(action: cancelInstructionEntry) {
                    Image(systemName: "xmark")
                        .font(.system(size: 11, weight: .bold))
                        .foregroundStyle(RecipeCreationColors.textSecondary)
                        .frame(width: 28, height: 28)
                        .background(Color.white.opacity(0.2), in: Circle())
                }
                .buttonStyle(.plain)
                .accessibilityLabel("Cancel instruction")
            }
        }
    }

    private var guidedModeRow: some View {
        GlassRow {
            VStack(alignment: .leading, spacing: 8) {
                HStack {
                    Text("Guided Cooking Mode")
                        .font(PlatePilotTheme.bodyFont(size: 16, weight: .semibold))
                        .foregroundStyle(RecipeCreationColors.textPrimary)

                    Spacer()

                    Toggle("Guided Cooking Mode", isOn: $draft.guidedMode)
                        .labelsHidden()
                        .tint(RecipeCreationColors.accent)
                }

                Text("Keeps steps on-screen while you cook")
                    .font(PlatePilotTheme.bodyFont(size: 12))
                    .foregroundStyle(RecipeCreationColors.textSecondary)
            }
        }
    }

    private var dietTagsRow: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Dietary Tags")
                .font(PlatePilotTheme.bodyFont(size: 14, weight: .semibold))
                .foregroundStyle(RecipeCreationColors.textSecondary)

            HStack(spacing: 10) {
                ForEach(DietTag.allCases, id: \.self) { tag in
                    TagPill(title: tag.title, isSelected: draft.tags.contains(tag)) {
                        toggleTag(tag)
                    }
                }
            }
        }
    }

    private var saveBar: some View {
        VStack(spacing: 12) {
            Button {
                Haptics.medium()
                Task { await saveRecipe() }
            } label: {
                ZStack {
                    RoundedRectangle(cornerRadius: 28, style: .continuous)
                        .fill(RecipeCreationColors.saveGradient)
                        .frame(height: 58)
                        .shadow(color: RecipeCreationColors.accent.opacity(0.35), radius: 18, x: 0, y: 10)

                    RoundedRectangle(cornerRadius: 28, style: .continuous)
                        .stroke(Color.white.opacity(0.2), lineWidth: 1)
                        .frame(height: 58)
                        .blendMode(.softLight)

                    if isSaving {
                        ProgressView()
                            .tint(.white)
                    } else {
                        Text("Save Recipe")
                            .font(PlatePilotTheme.bodyFont(size: 18, weight: .semibold))
                            .foregroundStyle(.white)
                    }
                }
            }
            .buttonStyle(.plain)
            .disabled(!canSave)
            .opacity(canSave ? 1 : 0.55)
        }
        .padding(.horizontal, RecipeCreationMetrics.horizontalPadding)
        .padding(.top, 12)
        .padding(.bottom, 12)
        .background(
            Rectangle()
                .fill(.ultraThinMaterial)
                .opacity(reduceTransparency ? 0.9 : 0.6)
                .blur(radius: 0.5)
                .ignoresSafeArea()
        )
    }

    private var canSave: Bool {
        !draft.title.trimmed().isEmpty && !draft.ingredients.isEmpty && !draft.instructions.isEmpty && !isSaving
    }

    private func sectionAnimation(delay: Double) -> some ViewModifier {
        RecipeCreationEntranceModifier(didAppear: didAppear, delay: delay)
    }

    private func toggleTag(_ tag: DietTag) {
        if draft.tags.contains(tag) {
            draft.tags.remove(tag)
        } else {
            draft.tags.insert(tag)
        }
        Haptics.light()
    }

    private func beginIngredientEntry() {
        isAddingIngredient = true
        isAddingInstruction = false
        focusedField = .newIngredient
    }

    private func beginInstructionEntry() {
        isAddingInstruction = true
        isAddingIngredient = false
        focusedField = .newInstruction
    }

    private func showUnitOverlay(target: UnitPickerTarget) {
        unitOverlayTarget = target
        newUnitName = ""
        withAnimation(.easeOut(duration: 0.2)) {
            isUnitOverlayVisible = true
        }
    }

    private func hideUnitOverlay() {
        withAnimation(.easeOut(duration: 0.2)) {
            isUnitOverlayVisible = false
        }
        unitOverlayTarget = nil
        newUnitName = ""
    }

    private func saveUnit() {
        let cleaned = newUnitName.trimmed()
        guard !cleaned.isEmpty else { return }

        let existing = availableUnits.first { $0.caseInsensitiveCompare(cleaned) == .orderedSame }
        let resolvedUnit = existing ?? cleaned
        if existing == nil {
            availableUnits.append(cleaned)
        }

        if let target = unitOverlayTarget {
            applyUnit(resolvedUnit, to: target)
        }

        // TODO: Persist custom units and load available units from the backend.
        Haptics.light()
        hideUnitOverlay()
    }

    private func applyUnit(_ unit: String, to target: UnitPickerTarget) {
        switch target {
        case .newIngredient:
            newIngredientUnit = unit
        case .ingredient(let id):
            guard let index = draft.ingredients.firstIndex(where: { $0.id == id }) else { return }
            draft.ingredients[index].unit = unit
        }
    }

    private func addIngredient() {
        let cleaned = newIngredientName.trimmed()
        guard !cleaned.isEmpty else { return }
        withAnimation(.spring(response: 0.35, dampingFraction: 0.8)) {
            draft.ingredients.append(
                IngredientDraft(
                    name: cleaned,
                    quantity: newIngredientQuantity,
                    unit: newIngredientUnit
                )
            )
        }
        newIngredientName = ""
        newIngredientQuantity = RecipeCreationMetrics.defaultQuantity
        newIngredientUnit = RecipeCreationMetrics.defaultUnit
        Haptics.light()
    }

    private func addInstruction() {
        let cleaned = newInstructionText.trimmed()
        guard !cleaned.isEmpty else { return }
        withAnimation(.spring(response: 0.35, dampingFraction: 0.8)) {
            draft.instructions.append(InstructionDraft(text: cleaned))
        }
        newInstructionText = ""
        Haptics.light()
    }

    private func cancelIngredientEntry() {
        newIngredientName = ""
        newIngredientQuantity = RecipeCreationMetrics.defaultQuantity
        newIngredientUnit = RecipeCreationMetrics.defaultUnit
        isAddingIngredient = false
        focusedField = nil
    }

    private func cancelInstructionEntry() {
        newInstructionText = ""
        isAddingInstruction = false
        focusedField = nil
    }

    private func removeIngredient(id: UUID) {
        if let index = draft.ingredients.firstIndex(where: { $0.id == id }) {
            withAnimation(.spring(response: 0.35, dampingFraction: 0.8)) {
                draft.ingredients.remove(at: index)
            }
            Haptics.rigid()
        }
    }

    private func removeInstruction(_ instruction: InstructionDraft) {
        if let index = draft.instructions.firstIndex(of: instruction) {
            withAnimation(.spring(response: 0.35, dampingFraction: 0.8)) {
                draft.instructions.remove(at: index)
            }
            Haptics.rigid()
        }
    }

    private func saveRecipe() async {
        guard canSave else { return }
        isSaving = true
        errorMessage = nil
        defer { isSaving = false }

        do {
            let ingredientNames = draft.ingredients.map { $0.name.trimmed() }.filter { !$0.isEmpty }
            // TODO: Send ingredient quantity/unit to the backend once the API supports it.
            _ = try await recipeStore.createRecipe(
                name: draft.title.trimmed(),
                description: draft.description.trimmed(),
                prepMinutes: draft.prepMinutes,
                cookMinutes: draft.cookMinutes,
                ingredients: ingredientNames,
                instructions: draft.instructions.map { $0.text.trimmed() }.filter { !$0.isEmpty },
                tags: draft.tags.map { $0.apiValue },
                guidedMode: draft.guidedMode
            )
            dismiss()
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to save recipe."
        }
    }

    private func errorBanner(message: String) -> some View {
        Text(message)
            .font(PlatePilotTheme.bodyFont(size: 13, weight: .semibold))
            .foregroundStyle(RecipeCreationColors.textPrimary)
            .padding(.horizontal, 16)
            .padding(.vertical, 10)
            .background(
                Capsule()
                    .fill(Color.white.opacity(0.6))
            )
            .overlay(
                Capsule()
                    .stroke(Color.white.opacity(0.3), lineWidth: 1)
            )
            .frame(maxWidth: .infinity)
    }
}

private enum FocusField: Hashable {
    case title
    case description
    case newIngredient
    case newInstruction
}

private enum TimePickerTarget: Identifiable {
    case prep
    case cook

    var id: String { title }

    var title: String {
        switch self {
        case .prep: return "Prep Time"
        case .cook: return "Cook Time"
        }
    }

    func binding(for draft: Binding<RecipeDraft>) -> Binding<Int> {
        switch self {
        case .prep:
            return draft.prepMinutes
        case .cook:
            return draft.cookMinutes
        }
    }
}

private enum UnitPickerTarget: Identifiable {
    case newIngredient
    case ingredient(UUID)

    var id: String {
        switch self {
        case .newIngredient:
            return "new"
        case .ingredient(let id):
            return id.uuidString
        }
    }
}

private struct RecipeDraft {
    var title: String = ""
    var description: String = ""
    var prepMinutes: Int = 10
    var cookMinutes: Int = 30
    var ingredients: [IngredientDraft] = []
    var instructions: [InstructionDraft] = []
    var guidedMode: Bool = false
    var tags: Set<DietTag> = []
}

private struct IngredientDraft: Identifiable, Hashable {
    let id = UUID()
    var name: String
    var quantity: String
    var unit: String
}

private struct InstructionDraft: Identifiable, Hashable {
    let id = UUID()
    var text: String
}

private enum DietTag: CaseIterable, Hashable {
    case vegetarian
    case vegan
    case glutenFree

    var title: String {
        switch self {
        case .vegetarian: return "Vegetarian"
        case .vegan: return "Vegan"
        case .glutenFree: return "Gluten-Free"
        }
    }

    var apiValue: String {
        switch self {
        case .vegetarian: return "vegetarian"
        case .vegan: return "vegan"
        case .glutenFree: return "gluten-free"
        }
    }
}

private enum RecipeCreationMetrics {
    static let horizontalPadding: CGFloat = 20
    static let sectionSpacing: CGFloat = 20
    static let cardRadius: CGFloat = 20
    static let cardPadding: CGFloat = 16
    static let rowRadius: CGFloat = 16
    static let hubBarHeight: CGFloat = 76
    static let defaultQuantity = "1"
    static let defaultUnit = "unit"
    static let quantityOptions = ["1/4", "1/2", "3/4", "1", "1 1/2", "2", "3", "4", "5"]
    static let defaultUnits = ["unit", "g", "kg", "ml", "l", "tsp", "tbsp", "cup", "pinch", "slice"]
}

private enum RecipeCreationColors {
    static let textPrimary = Color(red: 0.11, green: 0.1, blue: 0.09)
    static let textSecondary = Color(red: 0.43, green: 0.38, blue: 0.35)
    static let accent = Color(red: 0.88, green: 0.44, blue: 0.12)
    static let accentGlow = Color(red: 1.0, green: 0.7, blue: 0.42)
    static let saveGradient = LinearGradient(
        colors: [Color(red: 0.95, green: 0.5, blue: 0.24), Color(red: 0.87, green: 0.35, blue: 0.16)],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
}

private struct RecipeCreationBackground: View {
    let reduceTransparency: Bool

    var body: some View {
        ZStack {
            LinearGradient(
                colors: [
                    Color(red: 0.95, green: 0.55, blue: 0.22),
                    Color(red: 0.97, green: 0.7, blue: 0.4),
                    Color(red: 1.0, green: 0.95, blue: 0.91)
                ],
                startPoint: .top,
                endPoint: .bottom
            )
            .ignoresSafeArea()

            if !reduceTransparency {
                NoiseOverlay()
                    .opacity(0.08)
                    .blendMode(.softLight)
                    .ignoresSafeArea()
            }
        }
    }
}

private struct GlassCard<Content: View>: View {
    @Environment(\.accessibilityReduceTransparency) private var reduceTransparency

    let isFocused: Bool
    let content: () -> Content

    var body: some View {
        content()
            .padding(RecipeCreationMetrics.cardPadding)
            .background(glassBackground)
            .offset(y: isFocused ? -2 : 0)
            .animation(.easeOut(duration: 0.2), value: isFocused)
    }

    @ViewBuilder
    private var glassBackground: some View {
        let shape = RoundedRectangle(cornerRadius: RecipeCreationMetrics.cardRadius, style: .continuous)

        let base = ZStack {
            shape
                .fill(Color.white.opacity(reduceTransparency ? 0.9 : 0.08))

            shape
                .stroke(
                    LinearGradient(
                        colors: [Color.white.opacity(isFocused ? 0.4 : 0.28), Color.white.opacity(0.08)],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    ),
                    lineWidth: 1
                )
        }

        return base
            .applyGlassIfNeeded(
                reduceTransparency: reduceTransparency,
                cornerRadius: RecipeCreationMetrics.cardRadius,
                tint: .white.opacity(isFocused ? 0.2 : 0.14)
            )
            .shadow(color: RecipeCreationColors.accent.opacity(isFocused ? 0.24 : 0.18), radius: isFocused ? 26 : 20, x: 0, y: isFocused ? 12 : 10)
            .shadow(color: Color.black.opacity(isFocused ? 0.12 : 0.08), radius: isFocused ? 20 : 16, x: 0, y: isFocused ? 8 : 6)
    }
}

private struct GlassRow<Content: View>: View {
    @Environment(\.accessibilityReduceTransparency) private var reduceTransparency

    let content: () -> Content

    var body: some View {
        content()
            .padding(12)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(glassBackground)
    }
    
    @ViewBuilder
    private var glassBackground: some View {
        let shape = RoundedRectangle(cornerRadius: RecipeCreationMetrics.rowRadius, style: .continuous)
        if reduceTransparency {
            shape
                .fill(Color.white.opacity(0.9))
                .overlay(shape.stroke(Color.white.opacity(0.2), lineWidth: 1))
        } else {
            shape
                .fill(.ultraThinMaterial)
                .overlay(shape.stroke(Color.white.opacity(0.2), lineWidth: 1))
        }
    }
}

private struct SectionHeader: View {
    let title: String
    let actionLabel: String
    let action: () -> Void

    var body: some View {
        HStack {
            Text(title)
                .font(PlatePilotTheme.bodyFont(size: 16, weight: .semibold))
                .foregroundStyle(RecipeCreationColors.textPrimary)

            Spacer()

            Button(action: action) {
                Image(systemName: "plus")
                    .font(.system(size: 14, weight: .bold))
                    .foregroundStyle(.white)
                    .frame(width: 30, height: 30)
                    .background(RecipeCreationColors.saveGradient, in: Circle())
                    .shadow(color: RecipeCreationColors.accentGlow.opacity(0.6), radius: 10, x: 0, y: 6)
            }
            .buttonStyle(.plain)
            .accessibilityLabel(actionLabel)
        }
    }
}

private struct PhotoPickerButton: View {
    let photoImage: Image?

    var body: some View {
        ZStack {
            if let photoImage {
                photoImage
                    .resizable()
                    .scaledToFill()
                    .frame(width: 44, height: 44)
                    .clipShape(Circle())
            } else {
                Image(systemName: "camera.fill")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundStyle(RecipeCreationColors.accent)
            }
        }
        .frame(width: 44, height: 44)
        .background(Circle().fill(Color.white.opacity(0.15)))
        .overlay(Circle().stroke(Color.white.opacity(0.25), lineWidth: 1))
        .plateGlass(cornerRadius: 22, tint: .white.opacity(0.2), interactive: true)
        .accessibilityLabel("Add photo")
    }
}

private struct SelectBox: View {
    let title: String
    @Binding var selection: String
    let options: [String]

    var body: some View {
        Picker(selection: $selection) {
            ForEach(options, id: \.self) { option in
                Text(option).tag(option)
            }
        } label: {
            HStack(spacing: 6) {
                Text(selection)
                    .font(PlatePilotTheme.bodyFont(size: 13, weight: .semibold))
                    .foregroundStyle(RecipeCreationColors.textSecondary)
                Image(systemName: "chevron.down")
                    .font(.system(size: 10, weight: .semibold))
                    .foregroundStyle(RecipeCreationColors.textSecondary)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(.vertical, 6)
            .padding(.horizontal, 10)
            .background(
                RoundedRectangle(cornerRadius: 10, style: .continuous)
                    .fill(Color.white.opacity(0.2))
            )
            .overlay(
                RoundedRectangle(cornerRadius: 10, style: .continuous)
                    .stroke(Color.white.opacity(0.2), lineWidth: 1)
            )
        }
        .pickerStyle(.menu)
    }
}

private struct UnitAddButton: View {
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            Image(systemName: "plus")
                .font(.system(size: 11, weight: .bold))
                .foregroundStyle(.white)
                .frame(width: 24, height: 24)
                .background(RecipeCreationColors.saveGradient, in: Circle())
        }
        .buttonStyle(.plain)
        .accessibilityLabel("Add unit")
    }
}

private struct IngredientRow: View {
    @Binding var ingredient: IngredientDraft
    let quantityOptions: [String]
    let unitOptions: [String]
    let onAddUnit: () -> Void
    let onDelete: () -> Void

    var body: some View {
        GlassRow {
            VStack(alignment: .leading, spacing: 10) {
                HStack(spacing: 12) {
                    Circle()
                        .fill(RecipeCreationColors.accent.opacity(0.7))
                        .frame(width: 6, height: 6)

                    Text(ingredient.name)
                        .font(PlatePilotTheme.bodyFont(size: 15))
                        .foregroundStyle(RecipeCreationColors.textPrimary)

                    Spacer()

                    Button(action: onDelete) {
                        Image(systemName: "trash")
                            .font(.system(size: 12, weight: .semibold))
                            .foregroundStyle(RecipeCreationColors.textSecondary)
                            .frame(width: 28, height: 28)
                            .background(Color.white.opacity(0.2), in: Circle())
                    }
                    .buttonStyle(.plain)
                    .accessibilityLabel("Delete ingredient")
                }

                HStack(spacing: 10) {
                    SelectBox(title: "Quantity", selection: $ingredient.quantity, options: quantityOptions)

                    HStack(spacing: 6) {
                        SelectBox(title: "Unit", selection: $ingredient.unit, options: unitOptions)
                        UnitAddButton(action: onAddUnit)
                    }
                }
            }
        }
    }
}

private struct InstructionRow: View {
    let index: Int
    let text: String
    let onDelete: () -> Void

    var body: some View {
        GlassRow {
            HStack(alignment: .top, spacing: 12) {
                Text("\(index)")
                    .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                    .foregroundStyle(RecipeCreationColors.textPrimary)
                    .frame(width: 24, height: 24)
                    .background(Color.white.opacity(0.2), in: RoundedRectangle(cornerRadius: 8, style: .continuous))

                Text(text)
                    .font(PlatePilotTheme.bodyFont(size: 14))
                    .foregroundStyle(RecipeCreationColors.textPrimary)

                Spacer()

                Button(action: onDelete) {
                    Image(systemName: "trash")
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(RecipeCreationColors.textSecondary)
                        .frame(width: 26, height: 26)
                        .background(Color.white.opacity(0.2), in: Circle())
                }
                .buttonStyle(.plain)
                .accessibilityLabel("Delete step")
            }
        }
    }
}

private struct EmptyRow: View {
    let text: String

    var body: some View {
        GlassRow {
            Text(text)
                .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                .foregroundStyle(RecipeCreationColors.textSecondary)
        }
    }
}

private struct TimeRow: View {
    let prepMinutes: Int
    let cookMinutes: Int
    let onSelectPrep: () -> Void
    let onSelectCook: () -> Void

    var body: some View {
        HStack(spacing: 0) {
            Button(action: onSelectPrep) {
                TimeSegment(label: "Prep", minutes: prepMinutes)
            }
            .buttonStyle(.plain)

            Divider()
                .frame(width: 1)
                .overlay(Color.white.opacity(0.2))

            Button(action: onSelectCook) {
                TimeSegment(label: "Cook", minutes: cookMinutes)
            }
            .buttonStyle(.plain)
        }
        .padding(.vertical, 8)
        .background(
            RoundedRectangle(cornerRadius: 14, style: .continuous)
                .fill(Color.white.opacity(0.1))
        )
        .overlay(
            RoundedRectangle(cornerRadius: 14, style: .continuous)
                .stroke(Color.white.opacity(0.2), lineWidth: 1)
        )
    }
}

private struct TimeSegment: View {
    let label: String
    let minutes: Int

    private var timeLabel: String {
        if minutes == 0 {
            return "Add"
        }
        return "\(minutes) min"
    }

    var body: some View {
        HStack(spacing: 6) {
            Image(systemName: "clock")
                .font(.system(size: 12, weight: .semibold))
                .foregroundStyle(RecipeCreationColors.accent)

            Text("\(label) \(timeLabel)")
                .font(PlatePilotTheme.bodyFont(size: 13, weight: .semibold))
                .foregroundStyle(RecipeCreationColors.textSecondary)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 8)
    }
}

private struct TagPill: View {
    let title: String
    let isSelected: Bool
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            Text(title)
                .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                .foregroundStyle(isSelected ? .white : RecipeCreationColors.textSecondary)
                .padding(.horizontal, 14)
                .padding(.vertical, 8)
                .background(
                    Group {
                        if isSelected {
                            RecipeCreationColors.saveGradient
                        } else {
                            Color.white.opacity(0.12)
                        }
                    }
                )
                .clipShape(Capsule())
                .overlay(
                    Capsule()
                        .stroke(Color.white.opacity(isSelected ? 0.25 : 0.2), lineWidth: 1)
                )
                .shadow(color: isSelected ? RecipeCreationColors.accentGlow.opacity(0.5) : .clear, radius: 10, x: 0, y: 6)
        }
        .buttonStyle(.plain)
        .accessibilityLabel(title)
        .accessibilityValue(isSelected ? "On" : "Off")
    }
}

private struct PlaceholderTextEditor: View {
    @Binding var text: String
    let placeholder: String
    let isFocused: Bool

    var body: some View {
        ZStack(alignment: .topLeading) {
            if text.isEmpty {
                Text(placeholder)
                    .font(PlatePilotTheme.bodyFont(size: 15))
                    .foregroundStyle(RecipeCreationColors.textSecondary.opacity(0.7))
                    .padding(.top, 8)
                    .padding(.leading, 4)
            }

            TextEditor(text: $text)
                .font(PlatePilotTheme.bodyFont(size: 15))
                .foregroundStyle(RecipeCreationColors.textPrimary)
                .scrollContentBackground(.hidden)
                .background(Color.clear)
        }
    }
}

private struct RecipeCreationEntranceModifier: ViewModifier {
    let didAppear: Bool
    let delay: Double

    func body(content: Content) -> some View {
        content
            .opacity(didAppear ? 1 : 0)
            .offset(y: didAppear ? 0 : 14)
            .animation(.easeOut(duration: 0.5).delay(delay), value: didAppear)
    }
}

private struct UnitCreationOverlay: View {
    @Binding var unitName: String
    let onCancel: () -> Void
    let onSave: () -> Void

    var body: some View {
        ZStack {
            Color.black.opacity(0.25)
                .ignoresSafeArea()
                .onTapGesture(perform: onCancel)

            VStack(alignment: .leading, spacing: 16) {
                Text("Add Unit")
                    .font(PlatePilotTheme.titleFont(size: 20))
                    .foregroundStyle(RecipeCreationColors.textPrimary)

                TextField("e.g. tbsp, cup, g", text: $unitName)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()
                    .font(PlatePilotTheme.bodyFont(size: 15))
                    .padding(.vertical, 10)
                    .padding(.horizontal, 12)
                    .background(
                        RoundedRectangle(cornerRadius: 12, style: .continuous)
                            .fill(Color.white.opacity(0.25))
                    )
                    .overlay(
                        RoundedRectangle(cornerRadius: 12, style: .continuous)
                            .stroke(Color.white.opacity(0.2), lineWidth: 1)
                    )

                HStack(spacing: 12) {
                    Button("Cancel", action: onCancel)
                        .buttonStyle(.bordered)

                    Spacer()

                    Button("Save Unit", action: onSave)
                        .buttonStyle(.borderedProminent)
                        .disabled(unitName.trimmed().isEmpty)
                }
            }
            .padding(20)
            .background(
                RoundedRectangle(cornerRadius: 22, style: .continuous)
                    .fill(.ultraThinMaterial)
                    .overlay(
                        RoundedRectangle(cornerRadius: 22, style: .continuous)
                            .stroke(Color.white.opacity(0.25), lineWidth: 1)
                    )
            )
            .padding(.horizontal, 32)
        }
    }
}

private struct TimePickerSheet: View {
    @Environment(\.dismiss) private var dismiss
    @Binding var minutes: Int
    let title: String

    @State private var selection: Int = 0

    var body: some View {
        NavigationStack {
            VStack(spacing: 20) {
                Picker("Minutes", selection: $selection) {
                    ForEach(0...240, id: \.self) { value in
                        Text("\(value) min").tag(value)
                    }
                }
                .pickerStyle(.wheel)

                Button("Set Time") {
                    minutes = selection
                    dismiss()
                }
                .buttonStyle(.borderedProminent)
            }
            .padding(24)
            .navigationTitle(title)
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button("Done") {
                        minutes = selection
                        dismiss()
                    }
                }
            }
            .onAppear {
                selection = minutes
            }
        }
    }
}

private struct NoiseOverlay: View {
    private static let noiseImage: UIImage = NoiseOverlay.makeNoiseImage(size: 140)

    var body: some View {
        Rectangle()
            .fill(ImagePaint(image: Image(uiImage: Self.noiseImage), scale: 1))
            .ignoresSafeArea()
            .allowsHitTesting(false)
    }

    private static func makeNoiseImage(size: CGFloat) -> UIImage {
        let renderer = UIGraphicsImageRenderer(size: CGSize(width: size, height: size))
        return renderer.image { context in
            context.cgContext.setFillColor(UIColor.white.withAlphaComponent(0.12).cgColor)
            for _ in 0..<900 {
                let x = CGFloat.random(in: 0..<size)
                let y = CGFloat.random(in: 0..<size)
                let rect = CGRect(x: x, y: y, width: 1.2, height: 1.2)
                context.cgContext.fill(rect)
            }
        }
    }
}

private enum Haptics {
    static func light() {
        impact(.light)
    }

    static func medium() {
        impact(.medium)
    }

    static func rigid() {
        impact(.rigid)
    }

    private static func impact(_ style: UIImpactFeedbackGenerator.FeedbackStyle) {
        let generator = UIImpactFeedbackGenerator(style: style)
        generator.impactOccurred()
    }
}

private extension String {
    func trimmed() -> String {
        trimmingCharacters(in: .whitespacesAndNewlines)
    }
}

private extension View {
    @ViewBuilder
    func applyGlassIfNeeded(reduceTransparency: Bool, cornerRadius: CGFloat, tint: Color) -> some View {
        if reduceTransparency {
            self
        } else {
            plateGlass(cornerRadius: cornerRadius, tint: tint)
        }
    }
}
