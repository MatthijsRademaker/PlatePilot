import Foundation

struct Recipe: Identifiable, Hashable {
    let id: UUID
    var name: String
    var description: String
    var cuisine: String?
    var prepTime: String?
    var cookTime: String?
    var ingredients: [String]
    var directions: [String]
    var imageSeed: String

    var listImageURL: URL? {
        URL(string: "https://picsum.photos/seed/\(imageSeed)/360/240")
    }

    var detailImageURL: URL? {
        URL(string: "https://picsum.photos/seed/\(imageSeed)/900/600")
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
            ingredients: [],
            directions: [],
            imageSeed: "placeholder"
        )
    }
}

extension Recipe {
    init(dto: RecipeDTO) {
        let parsedID = UUID(uuidString: dto.id ?? "") ?? UUID()
        let name = dto.name?.trimmingCharacters(in: .whitespacesAndNewlines)
        let safeName = (name?.isEmpty == false) ? name! : "Untitled Recipe"
        let ingredients = dto.ingredients?
            .compactMap { $0.name?.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty } ?? []
        let directions = dto.directions?.filter { !$0.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty } ?? []

        self.init(
            id: parsedID,
            name: safeName,
            description: dto.description ?? "",
            cuisine: dto.cuisine?.name,
            prepTime: dto.prepTime,
            cookTime: dto.cookTime,
            ingredients: ingredients,
            directions: directions,
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
