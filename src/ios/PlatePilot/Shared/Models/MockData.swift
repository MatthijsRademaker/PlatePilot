import Foundation

enum MockData {
    static let recipes: [Recipe] = [
        Recipe(
            id: UUID(),
            name: "Citrus Herb Salmon",
            description: "Bright citrus glaze with fresh herbs and roasted vegetables.",
            cuisine: "Mediterranean",
            prepTime: "15 min",
            cookTime: "25 min",
            ingredients: [
                "2 salmon fillets",
                "1 orange, zested",
                "2 tbsp olive oil",
                "1 tbsp honey",
                "Fresh dill",
                "Sea salt"
            ],
            directions: [
                "Whisk orange zest, olive oil, and honey.",
                "Brush salmon with glaze and season.",
                "Roast at 400F until flaky.",
                "Finish with fresh dill."
            ],
            imageSeed: "citrus-salmon"
        ),
        Recipe(
            id: UUID(),
            name: "Miso Ginger Noodles",
            description: "Silky noodles tossed with miso, ginger, and crunchy veggies.",
            cuisine: "Japanese",
            prepTime: "20 min",
            cookTime: "10 min",
            ingredients: [
                "Soba noodles",
                "1 tbsp white miso",
                "1 tsp grated ginger",
                "Snap peas",
                "Sesame seeds"
            ],
            directions: [
                "Cook noodles until tender.",
                "Whisk miso, ginger, and warm water.",
                "Toss noodles with sauce and veggies.",
                "Top with sesame seeds."
            ],
            imageSeed: "miso-noodles"
        ),
        Recipe(
            id: UUID(),
            name: "Tuscan Chickpea Bowl",
            description: "Hearty chickpeas with roasted tomatoes and garlic yogurt.",
            cuisine: "Italian",
            prepTime: "15 min",
            cookTime: "20 min",
            ingredients: [
                "Chickpeas",
                "Cherry tomatoes",
                "Garlic",
                "Greek yogurt",
                "Lemon"
            ],
            directions: [
                "Roast tomatoes with garlic and olive oil.",
                "Warm chickpeas with herbs.",
                "Serve with lemon yogurt drizzle."
            ],
            imageSeed: "tuscan-bowl"
        ),
        Recipe(
            id: UUID(),
            name: "Berry Oat Parfait",
            description: "Layered oats, greek yogurt, and berries for a quick start.",
            cuisine: "American",
            prepTime: "10 min",
            cookTime: nil,
            ingredients: [
                "Rolled oats",
                "Greek yogurt",
                "Mixed berries",
                "Honey"
            ],
            directions: [
                "Toast oats lightly.",
                "Layer oats with yogurt and berries.",
                "Finish with honey drizzle."
            ],
            imageSeed: "berry-parfait"
        ),
        Recipe(
            id: UUID(),
            name: "Spiced Lentil Soup",
            description: "Comforting lentils with warming spices and fresh herbs.",
            cuisine: "Middle Eastern",
            prepTime: "15 min",
            cookTime: "35 min",
            ingredients: [
                "Red lentils",
                "Carrot",
                "Cumin",
                "Vegetable broth",
                "Parsley"
            ],
            directions: [
                "Sweat aromatics with spices.",
                "Add lentils and broth; simmer.",
                "Blend slightly and finish with parsley."
            ],
            imageSeed: "lentil-soup"
        ),
        Recipe(
            id: UUID(),
            name: "Avocado Citrus Toast",
            description: "Creamy avocado with citrus pop on toasted sourdough.",
            cuisine: "California",
            prepTime: "8 min",
            cookTime: nil,
            ingredients: [
                "Sourdough",
                "Ripe avocado",
                "Lime juice",
                "Chili flakes"
            ],
            directions: [
                "Toast sourdough slices.",
                "Mash avocado with lime and salt.",
                "Spread and sprinkle chili flakes."
            ],
            imageSeed: "avocado-toast"
        )
    ]

    static func weekPlan(starting startDate: Date, recipes: [Recipe]) -> WeekPlan {
        let calendar = Calendar.current
        let startOfWeek = calendar.date(from: calendar.dateComponents([.yearForWeekOfYear, .weekOfYear], from: startDate)) ?? startDate
        let days = (0..<7).compactMap { offset -> DayPlan? in
            guard let date = calendar.date(byAdding: .day, value: offset, to: startOfWeek) else { return nil }
            let meals = MealType.allCases.map { mealType in
                MealSlot(id: UUID(), mealType: mealType, recipe: nil)
            }
            return DayPlan(id: UUID(), date: date, meals: meals)
        }
        let endDate = calendar.date(byAdding: .day, value: 6, to: startOfWeek) ?? startOfWeek
        var week = WeekPlan(id: UUID(), startDate: startOfWeek, endDate: endDate, days: days)

        if let first = recipes.first, let firstIndex = week.days.indices.first {
            var firstDay = week.days[firstIndex]
            firstDay.meals = firstDay.meals.map { slot in
                if slot.mealType == .dinner {
                    var updated = slot
                    updated.recipe = first
                    return updated
                }
                return slot
            }
            week.days[firstIndex] = firstDay
        }

        return week
    }
}
