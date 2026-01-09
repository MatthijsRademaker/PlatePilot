import Foundation

struct Recipe: Identifiable, Hashable {
    let id: UUID
    var name: String
    var description: String
    var cuisine: String?
    var prepTime: String?
    var cookTime: String?
    var totalTime: String?
    var servings: Int?
    var ingredients: [String]
    var directions: [String]
    var tags: [String]
    var imageURL: URL?
    var imageSeed: String

    var listImageURL: URL? {
        imageURL ?? URL(string: "https://picsum.photos/seed/\(imageSeed)/360/240")
    }

    var detailImageURL: URL? {
        imageURL ?? URL(string: "https://picsum.photos/seed/\(imageSeed)/900/600")
    }
}

extension Recipe {
    static func placeholder(id: UUID = UUID()) -> Recipe {
        Recipe(
            id: id,
            name: "Loading",
            description: "",
            cuisine: nil,
            prepTime: nil,
            cookTime: nil,
            totalTime: nil,
            servings: nil,
            ingredients: [],
            directions: [],
            tags: [],
            imageURL: nil,
            imageSeed: "placeholder"
        )
    }
}

extension Recipe {
    init(dto: RecipeDTO) {
        let parsedID = UUID(uuidString: dto.id ?? "") ?? UUID()
        let trimmedName = dto.name?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        let safeName = trimmedName.isEmpty ? "Untitled Recipe" : trimmedName
        let ingredientLines = dto.ingredientLines ?? []
        let sortedLines = ingredientLines.sorted {
            ($0.sortOrder ?? 0) < ($1.sortOrder ?? 0)
        }
        let ingredients = sortedLines.compactMap { line -> String? in
            let name = line.ingredient?.name?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
            guard !name.isEmpty else { return nil }
            let quantityText = line.quantityText?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
            let unit = line.unit?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
            let quantityValue = line.quantityValue

            var quantityPart = ""
            if !quantityText.isEmpty {
                if quantityText.rangeOfCharacter(from: .letters) == nil, !unit.isEmpty {
                    quantityPart = "\(quantityText) \(unit)"
                } else {
                    quantityPart = quantityText
                }
            } else if let quantityValue {
                quantityPart = Recipe.formatNumber(quantityValue)
                if !unit.isEmpty {
                    quantityPart = "\(quantityPart) \(unit)"
                }
            } else if !unit.isEmpty {
                quantityPart = unit
            }

            let parts = [quantityPart, name].filter { !$0.isEmpty }
            return parts.joined(separator: " ")
        }

        let steps = dto.steps ?? []
        let sortedSteps = steps.sorted { ($0.stepIndex ?? 0) < ($1.stepIndex ?? 0) }
        let directions = sortedSteps.compactMap { step in
            let instruction = step.instruction?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
            return instruction.isEmpty ? nil : instruction
        }

        self.init(
            id: parsedID,
            name: safeName,
            description: dto.description ?? "",
            cuisine: dto.cuisine?.name,
            prepTime: Recipe.formatMinutes(dto.prepTimeMinutes),
            cookTime: Recipe.formatMinutes(dto.cookTimeMinutes),
            totalTime: Recipe.formatMinutes(dto.totalTimeMinutes),
            servings: dto.servings,
            ingredients: ingredients,
            directions: directions,
            tags: dto.tags ?? [],
            imageURL: URL(string: dto.imageUrl ?? ""),
            imageSeed: safeName.slugSeed()
        )
    }
}

extension Recipe {
    static func summary(id: UUID, name: String, description: String?) -> Recipe {
        let trimmedName = name.trimmingCharacters(in: .whitespacesAndNewlines)
        let safeName = trimmedName.isEmpty ? "Untitled Recipe" : trimmedName
        let trimmedDescription = description?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        return Recipe(
            id: id,
            name: safeName,
            description: trimmedDescription,
            cuisine: nil,
            prepTime: nil,
            cookTime: nil,
            totalTime: nil,
            servings: nil,
            ingredients: [],
            directions: [],
            tags: [],
            imageURL: nil,
            imageSeed: safeName.slugSeed()
        )
    }
}

private extension String {
    func slugSeed() -> String {
        let lowered = lowercased()
        let replaced = lowered.replacingOccurrences(of: "[^a-z0-9]+", with: "-", options: .regularExpression)
        return replaced.trimmingCharacters(in: CharacterSet(charactersIn: "-"))
    }
}

private extension Recipe {
    static func formatMinutes(_ minutes: Int?) -> String? {
        guard let minutes, minutes > 0 else { return nil }
        return "\(minutes) min"
    }

    static func formatNumber(_ value: Double) -> String {
        if value.truncatingRemainder(dividingBy: 1) == 0 {
            return String(format: "%.0f", value)
        }
        return String(format: "%.2f", value).replacingOccurrences(of: ".00", with: "")
    }
}
